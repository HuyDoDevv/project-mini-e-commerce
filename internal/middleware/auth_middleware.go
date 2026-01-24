package middleware

import (
	"net/http"
	"project-mini-e-commerce/pkg/auth"
	"project-mini-e-commerce/pkg/cache"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	jwtService   auth.TokenService
	cacheService cache.RedisCacheService
)

func InitAuthService(service auth.TokenService, cache cache.RedisCacheService) {
	jwtService = service
	cacheService = cache
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
		_, claims, err := jwtService.ParseToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			return
		}

		if jti, ok := claims["jti"].(string); ok {
			key := "blacklist:" + jti
			exists, err := cacheService.Exists(key)
			if err == nil && exists {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Token revoked",
				})
				return
			}
		}

		payload, err := jwtService.DecryptAccessToken(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
				"error":   err.Error(),
			})
			return
		}
		ctx.Set("user_uuid", payload.UserUUID)
		ctx.Set("user_role", payload.Role)
		ctx.Set("user_email", payload.Email)
		ctx.Next()
	}
}
