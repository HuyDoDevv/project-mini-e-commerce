package auth

import (
	"encoding/json"
	"project-mini-e-commerce/internal/db/sqlc"
	"project-mini-e-commerce/internal/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
}

//type Claims struct {
//	jwt.RegisteredClaims
//}

type EncryptedPayload struct {
	UserUUID string `json:"user_uuid"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

var (
	jwtSecret    = []byte(utils.GetEnv("JWT_SECRET", "JWT-Secret-Cho-Khoa-Lap-Trinh-Golang-Tu-Hoc-Cua-Huy"))
	jwtEncrytKey = []byte(utils.GetEnv("JWT_ENCRYPT_KEY", "JWT-Encrypt-Key-Cho-Khoa-Lap-Trinh-Golang-Tu-Hoc-Cua-Huy"))
)

const (
	AccessTokenTTL = 15 * time.Minute
)

func NewJWTService() TokenService {
	return &JWTService{}
}

func (js *JWTService) GenerateAccessToken(user sqlc.User) (string, error) {
	payload := EncryptedPayload{
		UserUUID: user.UserUuid.String(),
		Email:    user.UserEmail,
		Role:     string(rune(user.UserRole)),
	}

	rawData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	encryptedData, err := utils.EncryptAES(rawData, jwtEncrytKey)
	if err != nil {
		return "", err
	}
	//claims := &Claims{
	//	UserUUID: user.UserUuid.String(),
	//	Email:    user.UserEmail,
	//	Role:     string(rune(user.UserRole)),
	//	RegisteredClaims: jwt.RegisteredClaims{
	//		ID:        uuid.NewString(),
	//		ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
	//		IssuedAt:  jwt.NewNumericDate(time.Now()),
	//		Issuer:    "HuyDo",
	//	},
	//}

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
func (js *JWTService) GenerateRefreshToken() {}

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
	decryptedData, err := utils.DecryptAES(encryptedData, jwtEncrytKey)
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
