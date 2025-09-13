package main

import (
	"fmt"
	
	"repo-starter/internal/config"
	"repo-starter/internal/logger"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Printf("Error while getting config: %s", err)
	}
	fmt.Printf("%+v\n", cfg)
	log := logger.NewLogger(cfg)
	appLogger := log.With("component", "app")
	appLogger.With("key", "value").Info("Hello, World!")
}
