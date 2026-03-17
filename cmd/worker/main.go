package main

import (
	"path/filepath"
	"project-mini-e-commerce/internal/common"
	"project-mini-e-commerce/internal/config"
	"project-mini-e-commerce/internal/utils"
	"project-mini-e-commerce/pkg/logger"

	"github.com/joho/godotenv"
)

func NewWorker(cfg *config.Config) {
}

func main() {
	rootDir := utils.GetWorkingDir()
	logFile := filepath.Join(rootDir, "internal/logs/app.log")
	logger.InitLogger(logger.Config{
		Level:       "info",
		Filename:    logFile,
		MaxSize:     1,
		MaxAge:      5,
		MaxBackups:  5,
		Compress:    true,
		Environment: common.Environment(utils.GetEnv("APP_ENV", "development")),
	})

	if err := godotenv.Load(filepath.Join(rootDir, ".env")); err != nil {
		logger.Logger.Warn().Msg("No .env file found")
	} else {
		logger.Logger.Info().Msg(".env file loaded successfully in worker")
	}
	configFile := config.NewConfig()

	NewWorker(configFile)
}
