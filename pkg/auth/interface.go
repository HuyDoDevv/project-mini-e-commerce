package auth

import (
	"project-mini-e-commerce/internal/db/sqlc"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateAccessToken(user sqlc.User) (string, error)
	GenerateRefreshToken(user sqlc.User) (RefreshToken, error)
	ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error)
	DecryptAccessToken(tokenString string) (*EncryptedPayload, error)
	StoreRefreshToken(token RefreshToken) error
	ValidationRefreshToken(token string) (RefreshToken, error)
	RevokeRefreshToken(token string) error
}
