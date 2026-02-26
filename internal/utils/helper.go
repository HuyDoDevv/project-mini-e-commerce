package utils

import (
	"log"
	"os"
	"path/filepath"
	"project-mini-e-commerce/internal/common"
	"project-mini-e-commerce/pkg/logger"
	"strconv"

	"github.com/rs/zerolog"
)

func GetEnv(key, defaultValue string) string {
	valueKey := os.Getenv(key)
	if valueKey == "" {
		return defaultValue
	}
	return valueKey
}

func GetIntEnv(key string, defaultValue int) int {
	valueKey := os.Getenv(key)
	if valueKey == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(valueKey)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func GetWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln("unable to get working dir: ", err)
	}
	return dir
}

func NewLoggerWithPath(fileName, level string) *zerolog.Logger {
	path := filepath.Join(GetWorkingDir(), "internal/logs", fileName)

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
