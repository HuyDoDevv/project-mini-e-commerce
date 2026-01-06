package routes

import (
	"project-mini-e-commerce/internal/middleware"
	"project-mini-e-commerce/internal/utils"

	"github.com/gin-gonic/gin"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routers ...Route) {
	rateLimiterLogger := utils.NewLoggerWithPath("../../internal/logs/ratelimiter.log", "warning")
	httpLogger := utils.NewLoggerWithPath("../../internal/logs/http.log", "info")
	recoveryLogger := utils.NewLoggerWithPath("../../internal/logs/recovery.log", "warning")

	r.Use(
		middleware.LimiterMiddleware(rateLimiterLogger),
		middleware.TraceMiddleware(),
		middleware.LoggerMiddleware(httpLogger),
		middleware.RecoveryMiddleware(recoveryLogger),
		middleware.ApiKeyMiddleware(),
	)

	v1api := r.Group("api/v1")

	for _, route := range routers {
		route.Register(v1api)
	}
}
