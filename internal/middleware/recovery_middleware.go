package middleware

import (
	"bytes"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func RecoveryMiddleware(recoveryLogger zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()

				recoveryLogger.Error().
					Str("path", ctx.Request.URL.Path).
					Str("method", ctx.Request.Method).
					Str("client_ip", ctx.ClientIP()).
					Str("panic", fmt.Sprintf("%v", err)).
					Str("statck_at", ExtractFirstStackLine(stack)).
					Str("stack", string(stack)).
					Msg("panic recovered")
				ctx.AbortWithStatusJSON(500, gin.H{"code": "INTERNAL_SERVER_ERROR", "error": "try again later"})
			}
		}()

		ctx.Next()
	}
}

func ExtractFirstStackLine(stack []byte) string {
	lines := bytes.Split(stack, []byte("\n"))
	for _, line := range lines {
		if bytes.Contains(line, []byte(".go")) &&
			!bytes.Contains(line, []byte("/runtime/")) &&
			!bytes.Contains(line, []byte("/debug/")) &&
			!bytes.Contains(line, []byte("recovery_middleware.go")) {
			cleanLine := strings.TrimSpace(string(line))
			return cleanLine
		}
	}
	return ""
}
