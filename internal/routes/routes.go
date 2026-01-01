package routes

import (
	"project-mini-e-commerce/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routers ...Route) {
	httpLogger := zerolog.New(&lumberjack.Logger{
		Filename:   "../../internal/logs/http.log",
		MaxSize:    1, // megabytes MB
		MaxAge:     5, // 5 days
		MaxBackups: 5,
		Compress:   true, // cos nen la khong
		LocalTime:  true, // gio vi tri hien tai
	}).With().Timestamp().Logger()

	recoveryLogger := zerolog.New(&lumberjack.Logger{
		Filename:   "../../internal/logs/recovery.log",
		MaxSize:    1, // megabytes MB
		MaxAge:     5, // 5 days
		MaxBackups: 5,
		Compress:   true, // cos nen la khong
		LocalTime:  true, // gio vi tri hien tai
	}).With().Timestamp().Logger()

	r.Use(
		middleware.LoggerMiddleware(httpLogger),
		middleware.RecoveryMiddleware(recoveryLogger),
		middleware.ApiKeyMiddleware(),
		middleware.LimiterMiddleware(),
	)

	v1api := r.Group("api/v1")

	for _, route := range routers {
		route.Register(v1api)
	}
}
