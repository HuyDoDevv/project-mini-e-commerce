package v1service

import (
	"fmt"
	"project-mini-e-commerce/internal/db/sqlc"
	"project-mini-e-commerce/internal/repository"
	"project-mini-e-commerce/internal/utils"
	"project-mini-e-commerce/pkg/auth"
	"project-mini-e-commerce/pkg/cache"
	"project-mini-e-commerce/pkg/logger"
	"project-mini-e-commerce/pkg/mail"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
)

type authService struct {
	userRepo     repository.UserRepository
	tokenService auth.TokenService
	cacheService cache.RedisCacheService
	mailService  mail.EmailProviderService
}

func NewAuthService(repo repository.UserRepository, tokenService auth.TokenService, cacheService cache.RedisCacheService, mailService mail.EmailProviderService) AuthService {
	return &authService{
		userRepo:     repo,
		tokenService: tokenService,
		cacheService: cacheService,
		mailService:  mailService,
	}
}

type LimiterConfig struct {
	Rps   rate.Limit    // Requests per second
	Burst int           // Max burst size
	TTL   time.Duration // Time to live for client data
}
type LoginAttempt struct {
	Limiter  *rate.Limiter
	LastSeen time.Time
}

var (
	mu              sync.Mutex
	clients         = make(map[string]*LoginAttempt)
	LoginAttemptTTL = 5 * time.Minute
	MaxLoginAttempt = 5
)

func (as *authService) getIpClient(ctx *gin.Context) string {
	ip := ctx.ClientIP()
	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}
	return ip
}

func (as *authService) GetLoginAttempt(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()
	client, exits := clients[ip]
	if !exits {
		limiter := rate.NewLimiter(rate.Limit(float32(MaxLoginAttempt)/float32(LoginAttemptTTL.Seconds())), MaxLoginAttempt)
		newClient := &LoginAttempt{limiter, time.Now()}
		clients[ip] = newClient
		return limiter
	}
	client.LastSeen = time.Now()
	return client.Limiter
}

func (as *authService) CheckLoginAttempt(ip string) error {
	limiter := as.GetLoginAttempt(ip)
	if !limiter.Allow() {
		return utils.NewError("Too many login attempt. Please retry again login", utils.ErrCodeTooManyRequests)
	}
	return nil
}

func (as *authService) CleanupClient(ip string) {
	mu.Lock()
	defer mu.Unlock()
	delete(clients, ip)
}

func (as *authService) Login(ctx *gin.Context, email, password string) (string, string, int, error) {
	context := ctx.Request.Context()

	ip := as.getIpClient(ctx)
	err := as.CheckLoginAttempt(ip)
	if err != nil {
		return "", "", 0, err
	}

	email = utils.NormalizeString(email)
	user, err := as.userRepo.FindUserByEmail(context, email)
	if err != nil {
		as.GetLoginAttempt(ip)
		return "", "", 0, utils.NewError("Invalid email of password", utils.ErrCodeUnauthorized)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(password))
	if err != nil {
		as.GetLoginAttempt(ip)
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

	as.CleanupClient(ip)

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

func (as *authService) ForgotPassword(ctx *gin.Context, email string) error {
	context := ctx.Request.Context()
	rateLimitKey := fmt.Sprintf("reset:ratelimiter:%s:", email)

	exitst, err := as.cacheService.Exists(rateLimitKey)
	if err == nil && exitst {
		return utils.NewError("Too many requests", utils.ErrCodeTooManyRequests)
	}

	user, err := as.userRepo.FindUserByEmail(context, utils.NormalizeString(email))
	if err != nil {
		return utils.NewError("Email not found", utils.ErrCodeNotFound)
	}

	token, err := utils.GenerateTokenString(16)
	if err != nil {
		return utils.WrapError(err, "Failed to generate forgot password token", utils.ErrCodeInternal)
	}

	err = as.cacheService.Set("reset:"+token, user.UserUuid, 1*time.Hour)
	if err != nil {
		return utils.WrapError(err, "Failed to save reset token", utils.ErrCodeInternal)
	}

	err = as.cacheService.Set(rateLimitKey, "1", 5*time.Minute)
	if err != nil {
		return utils.WrapError(err, "Failed to set rate limit", utils.ErrCodeInternal)
	}

	resetLink := fmt.Sprintf("https://your-frontend-url/reset-password?token=%s", token)
	mailContent := &mail.Email{
		To: []mail.Address{
			{
				Email: email,
			},
		},
		Subject: "Password Reset Request",
		Text:    fmt.Sprintf("Hi %s,\n\nYou requested a password reset. Click the link below to reset your password:\n\n%s\n\nIf you didn't request this, please ignore this email.", user.UserName, resetLink),
	}
	err = as.mailService.SendEmail(context, mailContent)
	if err != nil {
		return utils.WrapError(err, "Failed to send reset password email", utils.ErrCodeInternal)
	}

	return nil
}

func (as *authService) ResetPassword(ctx *gin.Context, token, newPassword string) error {
	context := ctx.Request.Context()
	var userUUIDStr string
	err := as.cacheService.Get("reset:"+token, &userUUIDStr)
	if err == redis.Nil || userUUIDStr == "" {
		return utils.NewError("Invalid or expired token", utils.ErrCodeBadRequest)
	}

	if err != nil {
		return utils.WrapError(err, "Failed to retrieve token from cache", utils.ErrCodeInternal)
	}

	userUuid, err := uuid.Parse(userUUIDStr)
	if err != nil {
		return utils.NewError("Invalid user UUID in token", utils.ErrCodeInternal)
	}

	user, err := as.userRepo.FindUUID(context, userUuid)
	if err != nil {
		return utils.NewError("User not found", utils.ErrCodeNotFound)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return utils.WrapError(err, "Failed to hash password", utils.ErrCodeInternal)
	}

	user.UserPassword = string(hashedPassword)
	input := sqlc.UpdatePasswordParams{
		UserUuid:     user.UserUuid,
		UserPassword: user.UserPassword,
	}
	_, err = as.userRepo.UpdatePasswordParams(context, input)
	if err != nil {
		return utils.WrapError(err, "Failed to update password", utils.ErrCodeInternal)
	}

	err = as.cacheService.Clear("reset:" + token)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("Failed to delete reset token from cache")
	}

	return nil
}
