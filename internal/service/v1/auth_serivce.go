package v1service

import (
	"project-mini-e-commerce/internal/repository"
	"project-mini-e-commerce/internal/utils"
	"project-mini-e-commerce/pkg/auth"
	"project-mini-e-commerce/pkg/cache"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo     repository.UserRepository
	tokenService auth.TokenService
	cacheService cache.RedisCacheService
}

func NewAuthService(repo repository.UserRepository, tokenService auth.TokenService, cacheService cache.RedisCacheService) AuthService {
	return &authService{
		userRepo:     repo,
		tokenService: tokenService,
		cacheService: cacheService,
	}
}

func (as *authService) Login(ctx *gin.Context, email, password string) (string, string, int, error) {
	context := ctx.Request.Context()

	email = utils.NormalizeString(email)
	user, err := as.userRepo.FindUserByEmail(context, email)
	if err != nil {
		return "", "", 0, utils.NewError("Invalid email of password", utils.ErrCodeUnauthorized)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(password)); err != nil {
		return "", "", 0, utils.NewError("Invalid email of password", utils.ErrCodeUnauthorized)
	}

	accessToken, err := as.tokenService.GenerateAccessToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "failed to generate token", utils.ErrCodeInternal)
	}

	refreshToken, err := as.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "failed to generate token", utils.ErrCodeInternal)
	}
	if err := as.tokenService.StoreRefreshToken(refreshToken); err != nil {
		return "", "", 0, utils.WrapError(err, "cannot save refresh token", utils.ErrCodeInternal)
	}

	return accessToken, refreshToken.Token, int(auth.AccessTokenTTL.Seconds()), nil
}
func (as *authService) Logout(ctx *gin.Context, refreshToken string) error {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return utils.NewError("Missing Authorization header", utils.ErrCodeUnauthorized)
	}

	accessToken := strings.TrimPrefix(authHeader, "Bearer ")
	_, claims, err := as.tokenService.ParseToken(accessToken)
	if err != nil {
		return utils.WrapError(err, "Invalid access token", utils.ErrCodeInternal)
	}

	if jti, ok := claims["jti"].(string); ok {
		expUnix, _ := claims["exp"].(float64)
		exp := time.Unix(int64(expUnix), 0)
		key := "blacklist:" + jti
		ttl := time.Until(exp)
		as.cacheService.Set(key, "revoked", ttl)
	}

	_, err = as.tokenService.ValidationRefreshToken(refreshToken)
	if err != nil {
		return utils.WrapError(err, "Refresh token is invalid or revoked", utils.ErrCodeUnauthorized)
	}

	if err := as.tokenService.RevokeRefreshToken(refreshToken); err != nil {
		return utils.WrapError(err, "Cannot revoke refresh token", utils.ErrCodeInternal)
	}

	return nil
}

func (as *authService) RefreshToken(ctx *gin.Context, refreshTokenString string) (string, string, int, error) {
	context := ctx.Request.Context()
	var err error

	token, err := as.tokenService.ValidationRefreshToken(refreshTokenString)
	if err != nil {
		return "", "", 0, utils.NewError("Invalid token", utils.ErrCodeUnauthorized)
	}

	userUuid, err := uuid.Parse(token.UserUUID)
	if err != nil {
		return "", "", 0, utils.NewError("Cannot find user with uuid", utils.ErrCodeUnauthorized)
	}

	user, err := as.userRepo.FindUUID(context, userUuid)
	if err != nil {
		return "", "", 0, utils.NewError("User not found", utils.ErrCodeUnauthorized)
	}

	accessToken, err := as.tokenService.GenerateAccessToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "failed to generate token", utils.ErrCodeInternal)
	}

	refreshToken, err := as.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "failed to generate token", utils.ErrCodeInternal)
	}

	err = as.tokenService.RevokeRefreshToken(refreshTokenString)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "cannot revoke refresh token", utils.ErrCodeInternal)
	}

	err = as.tokenService.StoreRefreshToken(refreshToken)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "cannot save refresh token", utils.ErrCodeInternal)
	}

	return accessToken, refreshToken.Token, int(auth.AccessTokenTTL.Seconds()), nil
}
