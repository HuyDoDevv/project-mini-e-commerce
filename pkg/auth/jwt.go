package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"project-mini-e-commerce/internal/db/sqlc"
	"project-mini-e-commerce/internal/utils"
	"project-mini-e-commerce/pkg/cache"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	cacheService *cache.RedisCacheService
}

type EncryptedPayload struct {
	UserUUID string `json:"user_uuid"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type RefreshToken struct {
	Token     string    `json:"token"`
	UserUUID  string    `json:"user_uuid"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
}

var (
	jwtSecret     = []byte(utils.GetEnv("JWT_SECRET", "JWT-Secret-Cho-Khoa-Lap-Trinh-Golang-Tu-Hoc-Cua-Huy"))
	jwtEncryptKey = []byte(utils.GetEnv("JWT_ENCRYPT_KEY", "JWT-Encrypt-Key-Cho-Khoa-Lap-Trinh-Golang-Tu-Hoc-Cua-Huy"))
)

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
)

func NewJWTService(caseService *cache.RedisCacheService) TokenService {
	return &JWTService{
		cacheService: caseService,
	}
}

func (js *JWTService) GenerateAccessToken(user sqlc.User) (string, error) {
	payload := EncryptedPayload{
		UserUUID: user.UserUuid.String(),
		Email:    user.UserEmail,
		Role:     string(user.UserRole),
	}

	rawData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	encryptedData, err := utils.EncryptAES(rawData, jwtEncryptKey)
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"data": encryptedData,
		"jti":  uuid.NewString(),
		"exp":  jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
		"iat":  jwt.NewNumericDate(time.Now()),
		"iss":  "HuyDo",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
func (js *JWTService) GenerateRefreshToken(user sqlc.User) (RefreshToken, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return RefreshToken{}, err
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	return RefreshToken{
		Token:     token,
		UserUUID:  user.UserUuid.String(),
		ExpiresAt: time.Now().Add(RefreshTokenTTL),
		Revoked:   false,
	}, nil
}

func (js *JWTService) ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, nil, utils.NewError("Invalid token", utils.ErrCodeUnauthorized)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, utils.NewError("Invalid token", utils.ErrCodeUnauthorized)
	}
	return token, claims, nil
}

func (js *JWTService) DecryptAccessToken(tokenString string) (*EncryptedPayload, error) {
	_, claims, err := js.ParseToken(tokenString)
	if err != nil {
		return nil, utils.NewError("Cannot parse token", utils.ErrCodeUnauthorized)
	}
	encryptedData, ok := claims["data"].(string)
	if !ok {
		return nil, utils.NewError("Invalid token", utils.ErrCodeUnauthorized)
	}
	decryptedData, err := utils.DecryptAES(encryptedData, jwtEncryptKey)
	if err != nil {
		return nil, utils.WrapError(err, "Cannot decrypt token", utils.ErrCodeUnauthorized)
	}
	var payload EncryptedPayload
	err = json.Unmarshal(decryptedData, &payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}

func (js *JWTService) StoreRefreshToken(token RefreshToken) error {
	cacheKey := "refresh_token:" + token.Token
	return js.cacheService.Set(cacheKey, token, RefreshTokenTTL)
}

func (js *JWTService) ValidationRefreshToken(token string) (RefreshToken, error) {
	cacheKey := "refresh_token:" + token
	var refreshToken RefreshToken
	err := js.cacheService.Get(cacheKey, &refreshToken)
	if err != nil || refreshToken.Revoked || refreshToken.ExpiresAt.Before(time.Now()) {
		return RefreshToken{}, utils.WrapError(err, "Cannot get refresh token from cache", utils.ErrCodeInternal)
	}

	return refreshToken, nil
}

func (js *JWTService) RevokeRefreshToken(token string) error {
	cacheKey := "refresh_token:" + token
	var refreshToken RefreshToken
	err := js.cacheService.Get(cacheKey, &refreshToken)
	if err != nil {
		return utils.WrapError(err, "Cannot get refresh token from cache", utils.ErrCodeInternal)
	}
	refreshToken.Revoked = true
	return js.cacheService.Set(cacheKey, refreshToken, time.Until(refreshToken.ExpiresAt))
}
