package main

import (
	"flag"
	"fmt"
	"log/slog"

	"word_of_wisdom/internal/config"
	"word_of_wisdom/internal/pkg/cache"
	"word_of_wisdom/internal/server"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./../../config/config.yaml", "config path")

	config := config.InitConfig(configPath)
	if config == nil {
		slog.Error("main.main: config is nil")
		return
	}

	err := runServer(config)
	if err != nil {
		slog.Error(fmt.Sprintf("main.main: %v", err))
	}
}

func runServer(cfg *config.Config) error {
	serverAddress := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	cache := cache.NewCache()
	srv, err := server.NewServer(cache, serverAddress, config.GetConfig().Hash.FirstZerosCount)
	if err != nil {
		return fmt.Errorf("main.runServer: %w", err)
	}

	err = srv.Start()
	if err != nil {
		return fmt.Errorf("server.NewServer: %w", err)
	}

	return nil
}
