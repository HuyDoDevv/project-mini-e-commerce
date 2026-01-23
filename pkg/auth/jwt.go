package auth

import (
	"project-mini-e-commerce/internal/db/sqlc"
	"project-mini-e-commerce/internal/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
}

type Claims struct {
	UserUUID string `json:"user_uuid"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(utils.GetEnv("JWT_SECRET", "JWT-Secret-Cho-Khoa-Lap-Trinh-Golang-Tu-Hoc-Cua-Huy"))

const (
	AccessTokenTTL = 15 * time.Minute
)

func NewJWTService() TokenService {
	return &JWTService{}
}

func (js *JWTService) GenerateAccessToken(user sqlc.User) (string, error) {
	claims := &Claims{
		UserUUID: user.UserUuid.String(),
		Email:    user.UserEmail,
		Role:     string(rune(user.UserRole)),
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "HuyDo",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)

}
func (js *JWTService) GenerateRefreshToken() {}
