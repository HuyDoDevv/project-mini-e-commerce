package handler

import (
	"net/http"
	"project-mini-e-commerce/internal/dto"
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

type GetAllUserQuery struct {
	Search string `form:"search" binding:"omitempty,min=3,max=100,search"`
	Limit  int    `form:"limit" binding:"omitempty,gte=1,lte=100"`
	Page   int    `form:"page" binding:"omitempty,gte=1,lte=100"`
}

func (uh *UserHandler) GetAllUser(ctx *gin.Context) {
	var query GetAllUserQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	if query.Limit == 0 {
		query.Limit = 10
	}
	if query.Page == 0 {
		query.Page = 1
	}

	userList, err := uh.service.GetAllUser(query.Search, query.Limit, query.Page)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	userDtos := dto.MapUsersDTO(userList)
	utils.ResponseSuccess(ctx, http.StatusOK, userDtos)
}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {
	var input dto.CreateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	userParam := input.MapCreateUserToModel()

	createUser, err := uh.service.CreateUser(userParam)
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
	var param GetByUserUUIDParam
	if err := ctx.ShouldBindUri(&param); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	var input dto.UpdateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	userParam := input.MapUpdateUserToModel()

	updateUser, err := uh.service.UpdateUser(param.Uuid, userParam)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	userDTO := dto.MapUserDTO(updateUser)
	utils.ResponseSuccess(ctx, http.StatusCreated, userDTO)
}
func (uh *UserHandler) DeleteUser(ctx *gin.Context) {
	var param GetByUserUUIDParam
	if err := ctx.ShouldBindUri(&param); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	if err := uh.service.DeleteUser(param.Uuid); err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	utils.ResponseStatusCode(ctx, http.StatusNoContent)
}
