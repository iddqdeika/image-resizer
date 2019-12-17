package main

import (
	"image-resizer/internal/service"
	"image-resizer/pkg/cfg"
	"image-resizer/pkg/logger"
)

func main() {
	s, err := service.NewResizeService(logger.ConsoleLogger(logger.LogAll), cfg.EnvCfg)
	if err != nil {
		panic(err)
	}
	s.Run()
}
