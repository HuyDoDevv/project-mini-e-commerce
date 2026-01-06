package middleware

import (
	"context"
	"project-mini-e-commerce/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := ctx.GetHeader("X-TRACE-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		contextValue := context.WithValue(ctx.Request.Context(), logger.TraceIDKey, traceID)
		ctx.Request = ctx.Request.WithContext(contextValue)

		ctx.Writer.Header().Set("X-TRACE-ID", traceID)
		ctx.Set(string(logger.TraceIDKey), traceID)

		ctx.Next()
	}
}
