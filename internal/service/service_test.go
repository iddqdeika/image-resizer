package service

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestParseParams(t *testing.T) {

	imgProvider := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "test_img.png")
	}
	http.HandleFunc("/testimg", imgProvider)
	go http.ListenAndServe("localhost:8089", nil)

	cfg := &cfgStub{m: map[string]string{
		"relative-path":       "/",
		"concurrency":         "10",
		"download-queue-size": "100",
		"incoming-timeout":    "25",
		"download-timeout":    "25",
		"http-port":           "",
	}}

	s, err := newResizeService(nil, cfg)
	if err != nil {
		t.Fatal("unexpected err during service init: ", err)
	}

	req500, _ := http.NewRequest("GET", "http://localhost/?width=1920&height=1080&url=https://some-incorrect-image", nil)
	req400, _ := http.NewRequest("GET", "http://localhost/?width=1920&height=1080", nil)
	req200, _ := http.NewRequest("GET", "http://localhost/?width=1920&height=1080&url=http://localhost:8089/testimg", nil)

	r := httptest.NewRecorder()
	s.e.ServeHTTP(r, req500)
	if r.Code != 500 {
		t.Errorf("service must return 500 code for unreachable url")
	}

	r = httptest.NewRecorder()
	s.e.ServeHTTP(r, req400)
	if r.Code != 400 {
		t.Errorf("service must return 400 with uncomplete query params list")
	}

	r = httptest.NewRecorder()
	s.e.ServeHTTP(r, req200)
	if r.Code != 200 {
		b, _ := ioutil.ReadAll(r.Body)
		println(string(b))
		t.Errorf("service must return 200 with correct query and reachable correct image")
	}

}

type cfgStub struct {
	m map[string]string
}

func (c *cfgStub) StringWithDefaults(key string, defaultValue string) string {
	if val, ok := c.m[key]; ok {
		return val
	}
	return defaultValue
}

func (c *cfgStub) IntWithDefaults(key string, defaultValue int) int {
	if val, ok := c.m[key]; ok {
		result, err := strconv.Atoi(val)
		if err != nil {
			return defaultValue
		}
		return result
	}
	return defaultValue
}
