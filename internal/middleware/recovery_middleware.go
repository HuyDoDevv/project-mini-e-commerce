package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func RecoveryMiddleware(recoveryLogger zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				recoveryLogger.Error().
					Str("path", ctx.Request.URL.Path).
					Str("method", ctx.Request.Method).
					Str("client_ip", ctx.ClientIP()).
					Str("panic", fmt.Sprintf("%v", err)).
					Str("statck", string(debug.Stack())).
					Msg("panic recovered")
				ctx.AbortWithStatusJSON(500, gin.H{"code": "INTERNAL_SERVER_ERROR", "error": "try again later"})
			}
		}()

		ctx.Next()
	}
}
