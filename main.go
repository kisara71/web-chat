package main

import (
	"os"
	"web-chat/api"
	"web-chat/configs"
	"web-chat/internal/svc"
	"web-chat/pkg/logger"
)

func main() {
	lgr := logger.L()
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/dev.yaml"
	}
	cfg, err := configs.Load(configPath)
	if err != nil {
		lgr.Fatalf("load config error: %v", err)
	}
	svcCtx := svc.NewContext(cfg)
	router := api.NewRouter(svcCtx)
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	lgr.Printf("server start on %s", addr)
	if err := router.Run(addr); err != nil {
		lgr.Fatalf("server stopped: %v", err)
	}
}
