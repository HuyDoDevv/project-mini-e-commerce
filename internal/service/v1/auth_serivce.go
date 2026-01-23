package v1service

import (
	"project-mini-e-commerce/internal/repository"
	"project-mini-e-commerce/internal/utils"
	"project-mini-e-commerce/pkg/auth"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo     repository.UserRepository
	tokenService auth.TokenService
}

func NewAuthService(repo repository.UserRepository, tokenService auth.TokenService) AuthService {
	return &authService{
		userRepo:     repo,
		tokenService: tokenService,
	}
}

func (as *authService) Login(ctx *gin.Context, email, password string) (string, int, error) {
	context := ctx.Request.Context()

	email = utils.NormalizeString(email)
	user, err := as.userRepo.FindUserByEmail(context, email)
	if err != nil {
		return "", 0, utils.NewError("Invalid email of password", utils.ErrCodeUnauthorized)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(password)); err != nil {
		return "", 0, utils.NewError("Invalid email of password", utils.ErrCodeUnauthorized)
	}

	token, err := as.tokenService.GenerateAccessToken(user)
	if err != nil {
		return "", 0, utils.WrapError(err, "failed to generate token", utils.ErrCodeInternal)
	}

	return token, 100, nil
}
func (as *authService) Logout(ctx *gin.Context) error {
	return nil
}
