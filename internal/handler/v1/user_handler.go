package v1handler

import (
	"net/http"
	v1dto "project-mini-e-commerce/internal/dto/v1"
	v1service "project-mini-e-commerce/internal/service/v1"
	"project-mini-e-commerce/internal/utils"
	"project-mini-e-commerce/internal/validation"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service v1service.UserService
}

func NewUserHandler(service v1service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) GetAllUser(ctx *gin.Context) {}
func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var input v1dto.CreateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	user := input.MapCreateInputToModel()

	createdUser, err := uh.service.CreateUser(ctx, user)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	userDTO := v1dto.MapUserToDTO(createdUser)

	utils.ResponseSuccess(ctx, http.StatusCreated, userDTO)
}
func (uh *UserHandler) GetByUserUUID(ctx *gin.Context) {}
func (uh *UserHandler) UpdateUser(ctx *gin.Context) {

}
func (uh *UserHandler) DeleteUser(ctx *gin.Context) {

}
