package main

import (
	"os"
	"os/signal"
	"syscall"

	"convertpdfgo/api"
	"convertpdfgo/config"
	"convertpdfgo/pkg/gotenberg"
	"convertpdfgo/pkg/logger"
	"convertpdfgo/service"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.ServiceName)

	gotClient := gotenberg.New(cfg.GotenbergURL)
	services := service.New(nil, log, gotClient)

	// HTTP Server
	server := api.New(cfg, log, services)
	go func() {
		if err := server.Run(); err != nil {
			log.Error("failed to start http server", logger.Error(err))
		}
	}()
	log.Info("HTTP server started", logger.String("port", cfg.AppPort))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down...")
}
