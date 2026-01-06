package utils

import (
	"os"
	"project-mini-e-commerce/internal/common"
	"project-mini-e-commerce/pkg/logger"

	"github.com/rs/zerolog"
)

func GetEnv(key, defaultValue string) string {
	valueKey := os.Getenv(key)
	if valueKey == "" {
		return defaultValue
	}
	return valueKey
}

func NewLoggerWithPath(path, level string) *zerolog.Logger {
	config := logger.Config{
		Level:       level,
		Filename:    path,
		MaxSize:     1,
		MaxAge:      5,
		MaxBackups:  5,
		Compress:    true,
		Environment: common.Environment(GetEnv("APP_ENV", "development")),
	}

	return logger.NewLogger(config)
}
