package main

import (
	"path/filepath"
	"project-mini-e-commerce/internal/app"
	"project-mini-e-commerce/internal/common"
	"project-mini-e-commerce/internal/config"
	"project-mini-e-commerce/internal/utils"
	"project-mini-e-commerce/pkg/logger"

	"github.com/joho/godotenv"
)

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
		logger.Logger.Info().Msg(".env file loaded successfully in api")
	}
	configFile := config.NewConfig()
	application, err := app.NewApplication(configFile)
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("Failed to create application")
	}

	if err := application.Run(); err != nil {
		logger.Logger.Fatal().Err(err).Msg("Failed to run application")
	}
}
