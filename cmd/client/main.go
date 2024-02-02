package main

import (
	"flag"
	"fmt"
	"log/slog"

	"word_of_wisdom/internal/client"
	"word_of_wisdom/internal/config"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./../../config/config.yaml", "config path")

	config := config.InitConfig(configPath)

	err := runClient(config)
	if err != nil {
		slog.Error(fmt.Sprintf("main.main: %v", err))
	}
}

func runClient(cfg *config.Config) error {
	serverAddress := fmt.Sprintf("%s:%d", cfg.Client.Host, cfg.Client.Port)

	cli, err := client.NewClient(serverAddress)
	if err != nil {
		return fmt.Errorf("main.runClient - create: %w", err)
	}
	err = cli.Start()
	if err != nil {
		return fmt.Errorf("main.runClient - start: %w", err)
	}
	return nil
}
