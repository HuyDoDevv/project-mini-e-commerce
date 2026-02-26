package v1handler

import (
	"net/http"
	v1dto "project-mini-e-commerce/internal/dto/v1"
	v1service "project-mini-e-commerce/internal/service/v1"
	"project-mini-e-commerce/internal/utils"
	"project-mini-e-commerce/internal/validation"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service v1service.AuthService
}

func NewAuthHandler(service v1service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (ah *AuthHandler) Login(ctx *gin.Context) {
	var params v1dto.LoginInput
	if err := ctx.ShouldBindJSON(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	accessToken, refreshToken, expiresIn, err := ah.service.Login(ctx, params.Email, params.Password)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	responseToken := v1dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}

	utils.ResponseSuccess(ctx, http.StatusOK, "Login success", responseToken)
}

func (ah *AuthHandler) Logout(ctx *gin.Context) {
	var params v1dto.RefreshTokenInput
	if err := ctx.ShouldBindJSON(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	if err := ah.service.Logout(ctx, params.RefreshToken); err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, "Logout success")
}

func (ah *AuthHandler) RefreshToken(ctx *gin.Context) {
	var params v1dto.RefreshTokenInput
	if err := ctx.ShouldBindJSON(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	accessToken, refreshToken, expiresIn, err := ah.service.RefreshToken(ctx, params.RefreshToken)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	responseToken := v1dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}
	utils.ResponseSuccess(ctx, http.StatusOK, "Refresh token success", responseToken)
}

func (ah *AuthHandler) ForgotPassword(ctx *gin.Context) {
	var params v1dto.ForgotPasswordInput
	if err := ctx.ShouldBindJSON(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	if err := ah.service.ForgotPassword(ctx, params.Email); err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, "Forgot password success")
}

func (ah *AuthHandler) ResetPassword(ctx *gin.Context) {
	var params v1dto.ResetPasswordInput
	if err := ctx.ShouldBindJSON(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	if err := ah.service.ResetPassword(ctx, params.Token, params.NewPassword); err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, "Reset password success")
}
