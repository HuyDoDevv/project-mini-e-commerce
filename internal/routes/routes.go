package routes

import (
	"project-mini-e-commerce/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routers ...Route) {
	r.Use(
		middleware.LoggerMiddleware(),
		middleware.ApiKeyMiddleware(),
		middleware.LimiterMiddleware(),
	)

	api := r.Group("api/v1")

	for _, route := range routers {
		route.Register(api)
	}
}
