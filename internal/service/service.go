package service

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"image-resizer/internal/defaultlogger"
	"image-resizer/internal/definitions"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	"golang.org/x/image/webp"
)

const (
	DefaultResizeRelativePath    = "/"
	DefaultProcessingConcurrency = 16
	DefaultDownloadQueueSize     = 1024
	DefaultDownloadTimeoutInSec  = 25
	DefaultIncomingTimeoutInSec  = 25
)

// Creates new instance of ResizerService with given (or default) logger and config
// Returns error in case of incorrect config given
func NewResizeService(logger definitions.Logger, config definitions.Config) (definitions.Service, error) {
	if logger == nil {
		logger = defaultlogger.Logger
	}
	cfg, err := parseCfg(config)
	if err != nil {
		return nil, fmt.Errorf("cant parse service config: %v", err)
	}

	res := &resizer{
		e:   gin.Default(),
		l:   logger,
		cfg: cfg,
		db:  make(chan struct{}, cfg.downloadQueueSize),
		pb:  make(chan struct{}, cfg.concurrency),
	}

	res.fillRoutes()
	return res, nil
}

// Parses values from given config and returns filled Cfg or error
func parseCfg(config definitions.Config) (resizerCfg, error) {
	cfg := resizerCfg{}

	cfg.resizeRelativePath = config.StringWithDefaults("path", DefaultResizeRelativePath)
	cfg.concurrency = config.IntWithDefaults("concurrency", DefaultProcessingConcurrency)
	cfg.downloadQueueSize = config.IntWithDefaults("download-queue-size", DefaultDownloadQueueSize)
	cfg.incomingTimeoutInSec = config.IntWithDefaults("incoming-timeout", DefaultIncomingTimeoutInSec)
	cfg.downloadTimeoutInSec = config.IntWithDefaults("download-timeout", DefaultDownloadTimeoutInSec)

	return cfg, nil
}

type resizerCfg struct {
	// relative path
	resizeRelativePath string
	// maximum concurrent processing units
	concurrency int
	// maximum active requests awaiting download.
	// if active requests exeeds this value, we must return sorry
	downloadQueueSize int
	// must return either result or sorry message within this timeout
	incomingTimeoutInSec int
	// http timeout while image downloading
	downloadTimeoutInSec int
}

type resizer struct {
	e   *gin.Engine
	l   definitions.Logger
	cfg resizerCfg
	db  chan struct{} //download balanser
	pb  chan struct{} //concurrency balanser
}

func (r *resizer) Run(addr ...string) {
	r.e.Run(addr...)
}

func (r *resizer) fillRoutes() {
	r.e.GET(r.cfg.resizeRelativePath, r.resizeHandler)
}

func (r *resizer) resizeHandler(ctx *gin.Context) {

	r.l.Infof("incoming request, url: %v\r\n", ctx.Request.URL.String())
	t := time.NewTicker(time.Second * time.Duration(r.cfg.incomingTimeoutInSec))
	defer t.Stop()

	//check and  incoming params
	params, err := r.parseParams(ctx)
	if err != nil {
		r.l.Errorf("cant parse config: %v\r\n", err)
		ctx.JSON(http.StatusBadRequest, &message{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
		return
	}
	r.l.Infof("params got: %v\r\n", params)

	//wait for download slot
	select {
	case r.db <- struct{}{}:
		defer func() { <-r.db }()
	default:
		r.l.Errorf("there is no slots for download await\r\n", err)
		ctx.JSON(http.StatusInternalServerError, &message{
			StatusCode: http.StatusInternalServerError,
			Message:    "too many requests in a row",
		})
		return
	}

	sourceData, err := r.download(params)
	if err != nil {
		r.l.Errorf("error during image downloading: %v\r\n", err)
		ctx.JSON(http.StatusInternalServerError, &message{
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		})
		return
	}

	select {
	case r.pb <- struct{}{}:
		defer func() { <-r.pb }()
	case <-t.C:
		r.l.Errorf("timeout exceeded while awaiting processing quota\r\n")
		ctx.JSON(http.StatusInternalServerError, &message{
			StatusCode: http.StatusInternalServerError,
			Message:    fmt.Sprintf("cannot process image within timeout %v", r.cfg.incomingTimeoutInSec),
		})
		return
	}

	resultData, mimeType, err := r.resize(sourceData, params)
	if err != nil {
		r.l.Errorf("error during image resizing: %v\r\n", err)
		ctx.JSON(http.StatusInternalServerError, &message{
			StatusCode: http.StatusInternalServerError,
			Message:    "internal error while image processing",
		})
		return
	}

	ctx.Data(http.StatusOK, mimeType, resultData)
}

type message struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type params struct {
	width  int
	height int
	url    string
}

// Check and parse query params
func (r *resizer) parseParams(ctx *gin.Context) (params, error) {
	w := ctx.Request.URL.Query().Get("width")
	if len(w) == 0 {
		return params{}, fmt.Errorf("width param needed")
	}
	width, err := strconv.Atoi(w)
	if err != nil {
		return params{}, fmt.Errorf("width must be integer")
	}
	h := ctx.Request.URL.Query().Get("height")
	if len(h) == 0 {
		return params{}, fmt.Errorf("height param needed")
	}
	height, err := strconv.Atoi(h)
	if err != nil {
		return params{}, fmt.Errorf("height must be integer")
	}
	url := ctx.Request.URL.Query().Get("url")
	if len(url) == 0 {
		return params{}, fmt.Errorf("url param needed")
	}
	return params{height: height, width: width, url: url}, nil
}

func (r *resizer) download(params params) ([]byte, error) {
	c := http.DefaultClient
	//if duration is 0 - may be connection leaks
	c.Timeout = time.Duration(r.cfg.downloadTimeoutInSec) * time.Second

	//prepare request
	req, err := http.NewRequest("GET", params.url, nil)
	if err != nil {
		return nil, fmt.Errorf("error during image download request creation: %v", err)
	}

	//send request
	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error during image downloading: %v", err)
	}

	//be aware: http.timeot can interrupt res.body reading
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error during response data retrieve: %v", err)
	}
	return data, nil
}

// Resize photo and return data and MIME-type or error
func (r *resizer) resize(data []byte, p params) ([]byte, string, error) {
	t := http.DetectContentType(data)
	var i image.Image
	var err error
	switch t {
	case "image/gif":
		i, err = gif.Decode(bytes.NewReader(data))
	case "image/jpeg", "image/pjpeg":
		i, err = jpeg.Decode(bytes.NewReader(data))
	case "image/png":
		i, err = png.Decode(bytes.NewReader(data))
	case "image/webp":
		i, err = webp.Decode(bytes.NewReader(data))
	default:
		return nil, "", fmt.Errorf("unsupportet filetype: %v", t)
	}
	if err != nil {
		return nil, "", fmt.Errorf("error while image decoding: %v", err)
	}
	res := resize.Resize(uint(p.width), uint(p.height), i, resize.Bilinear)
	buf := &bytes.Buffer{}
	switch t {
	case "image/gif":
		err = gif.Encode(buf, res, nil)
	case "image/jpeg", "image/pjpeg":
		err = jpeg.Encode(buf, res, nil)
	case "image/png":
		err = png.Encode(buf, res)
	case "image/webp":
		err = jpeg.Encode(buf, res, nil)
		t = "image/jpeg"
	}
	if err != nil {
		return nil, "", fmt.Errorf("error while image encoding: %v", err)
	}
	out, err := ioutil.ReadAll(buf)
	if err != nil {
		return nil, "", fmt.Errorf("error while image encoding: %v", err)
	}
	return out, t, nil
}
