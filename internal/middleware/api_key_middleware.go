package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func ApiKeyMiddleware() gin.HandlerFunc {
	expeckey := os.Getenv("API_KEY")
	if expeckey == "" {
		expeckey = "Ex Key"
	}

	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-API-KEY")
		if apiKey == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, "Missing X API KEY")
			return
		}
		if apiKey != expeckey {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, "Incorrect X API KEY")
			return
		}
		ctx.Next()
	}
}
