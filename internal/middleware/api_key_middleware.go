package middleware

import (
	"net/http"
	"project-mini-e-commerce/internal/utils"

	"github.com/gin-gonic/gin"
)

func ApiKeyMiddleware() gin.HandlerFunc {
	expect := utils.GetEnv("API_KEY", "API_KEY")

	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-API-KEY")
		if apiKey == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, "Missing X API KEY")
			return
		}
		if apiKey != expect {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Incorrect X API KEY")
			return
		}
		ctx.Next()
	}
}
