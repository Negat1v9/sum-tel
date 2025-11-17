package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Negat1v9/sum-tel/services/core/internal/app"
	"github.com/Negat1v9/sum-tel/shared/config"
	"github.com/Negat1v9/sum-tel/shared/logger"
)

func main() {
	// load config
	cfg, err := config.LoadCoreConfig("./config/config")

	if err != nil {
		log.Printf("[error] load config: %v", err)
		os.Exit(1)
	}
	fmt.Println(cfg.AppConfig.Env)

	logger := logger.NewLogger(cfg.AppConfig.Env)

	app := app.NewApp(cfg, logger)
	if err := app.Run(); err != nil {
		logger.Errorf("app run error: %v", err)
		os.Exit(1)
	}
}
