package handler

import (
	"net/http"
	"project-mini-e-commerce/internal/dto"
	"project-mini-e-commerce/internal/models"
	"project-mini-e-commerce/internal/service"
	"project-mini-e-commerce/internal/utils"
	"project-mini-e-commerce/internal/validation"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

type GetByUserUUIDParam struct {
	Uuid string `uri:"uuid" binding:"uuid"`
}

func (uh *UserHandler) GetAllUser(ctx *gin.Context) {
	userList, err := uh.service.GetAllUser()
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	userDtos := dto.MapUsersDTO(userList)
	utils.ResponseSuccess(ctx, http.StatusOK, userDtos)
}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	createUser, err := uh.service.CreateUser(user)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	userDTO := dto.MapUserDTO(createUser)
	utils.ResponseSuccess(ctx, http.StatusCreated, userDTO)
}

func (uh *UserHandler) GetByUserUUID(ctx *gin.Context) {
	var param GetByUserUUIDParam
	if err := ctx.ShouldBindUri(&param); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	user, err := uh.service.GetByUserUUID(param.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
	}

	userDTO := dto.MapUserDTO(user)
	utils.ResponseSuccess(ctx, http.StatusCreated, userDTO)
}
func (uh *UserHandler) UpdateUser(ctx *gin.Context) {

}
func (uh *UserHandler) DeleteUser(ctx *gin.Context) {

}
