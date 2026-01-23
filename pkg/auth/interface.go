package auth

import (
	"project-mini-e-commerce/internal/db/sqlc"
)

type TokenService interface {
	GenerateAccessToken(user sqlc.User) (string, error)
	GenerateRefreshToken()
}
