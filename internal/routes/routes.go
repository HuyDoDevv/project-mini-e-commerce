package routes

import (
	"project-mini-e-commerce/internal/middleware"
	v1routes "project-mini-e-commerce/internal/routes/v1"
	"project-mini-e-commerce/internal/utils"
	"project-mini-e-commerce/pkg/auth"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, authService auth.TokenService, routers ...Route) {
	rateLimiterLogger := utils.NewLoggerWithPath("../../internal/logs/ratelimiter.log", "warning")
	httpLogger := utils.NewLoggerWithPath("../../internal/logs/http.log", "info")
	recoveryLogger := utils.NewLoggerWithPath("../../internal/logs/recovery.log", "warning")

	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.Use(
		middleware.LimiterMiddleware(rateLimiterLogger),
		middleware.TraceMiddleware(),
		middleware.LoggerMiddleware(httpLogger),
		middleware.RecoveryMiddleware(recoveryLogger),
		middleware.ApiKeyMiddleware(),
	)

	v1api := r.Group("api/v1")

	middleware.InitAuthService(authService)
	protected := v1api.Group("")
	protected.Use(middleware.AuthMiddleware())

	for _, route := range routers {
		switch route.(type) {
		case *v1routes.AuthRoutes:
			route.Register(v1api)
		default:
			route.Register(protected)
		}
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "not found"})
	})
}
