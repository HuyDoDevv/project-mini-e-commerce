package middleware

import (
	"net/http"
	"project-mini-e-commerce/pkg/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	jwtService auth.TokenService
)

func InitAuthService(service auth.TokenService) {
	jwtService = service
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
				"error":   "Missing or invalid Authorization header",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		_, _, err := jwtService.ParseToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}
		payload, err := jwtService.DecryptAccessToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
		}
		ctx.Set("user_uuid", payload.UserUUID)
		ctx.Set("user_role", payload.Role)
		ctx.Set("user_email", payload.Email)
		ctx.Next()
	}
}
