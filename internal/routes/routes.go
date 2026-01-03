package routes

import (
	"project-mini-e-commerce/internal/common"
	"project-mini-e-commerce/internal/middleware"
	"project-mini-e-commerce/internal/utils"
	"project-mini-e-commerce/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routers ...Route) {
	httpLogger := newLoggerWithPath("../../internal/logs/http.log", "info")
	recoveryLogger := newLoggerWithPath("../../internal/logs/recovery.log", "warning")

	r.Use(
		middleware.LoggerMiddleware(*httpLogger),
		middleware.RecoveryMiddleware(*recoveryLogger),
		middleware.ApiKeyMiddleware(),
		middleware.LimiterMiddleware(),
	)

	v1api := r.Group("api/v1")

	for _, route := range routers {
		route.Register(v1api)
	}
}

func newLoggerWithPath(level, path string) *zerolog.Logger {
	config := logger.Config{
		Level:       level,
		Filename:    path,
		MaxSize:     1,
		MaxAge:      5,
		MaxBackups:  5,
		Compress:    true,
		Environment: common.Environment(utils.GetEnv("APP_ENV", "development")),
	}

	return logger.NewLogger(config)
}
