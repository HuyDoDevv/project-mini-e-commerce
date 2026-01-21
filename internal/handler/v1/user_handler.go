package v1handler

import (
	"net/http"
	v1dto "project-mini-e-commerce/internal/dto/v1"
	v1service "project-mini-e-commerce/internal/service/v1"
	"project-mini-e-commerce/internal/utils"
	"project-mini-e-commerce/internal/validation"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service v1service.UserService
}

func NewUserHandler(service v1service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (uh *UserHandler) GetAllUser(ctx *gin.Context) {
	var params v1dto.GetUsersParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	users, countUser, err := uh.service.GetAllUser(ctx, params.Search, params.Order, params.Sort, params.Limit, params.Page)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	userDTO := v1dto.MapUsersToDTO(users)
	userPagination := utils.NewPaginationResponse(userDTO, params.Limit, params.Page, countUser)
	utils.ResponseSuccess(ctx, http.StatusOK, "ok", userPagination)
}
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

	utils.ResponseSuccess(ctx, http.StatusCreated, "ok", userDTO)
}
func (uh *UserHandler) GetByUserUUID(ctx *gin.Context) {}

func (uh *UserHandler) UpdateUser(ctx *gin.Context) {
	var params v1dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	uuidParse, err := uuid.Parse(params.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	var input v1dto.UpdateUserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}

	user := input.MapUpdateInputToModel(uuidParse)

	updateUser, err := uh.service.UpdateUser(ctx, user)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, "ok", updateUser)
}
func (uh *UserHandler) DeleteUser(ctx *gin.Context) {
	var params v1dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	uuidParse, err := uuid.Parse(params.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	if err := uh.service.DeleteUser(ctx, uuidParse); err != nil {
		utils.ResponseError(ctx, err)
	}
	utils.ResponseStatusCode(ctx, http.StatusOK)
}
func (uh *UserHandler) RestoreUser(ctx *gin.Context) {
	var params v1dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	uuidParse, err := uuid.Parse(params.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	if err := uh.service.RestoreUser(ctx, uuidParse); err != nil {
		utils.ResponseError(ctx, err)
	}
	utils.ResponseStatusCode(ctx, http.StatusOK)
}
func (uh *UserHandler) TrashUser(ctx *gin.Context) {
	var params v1dto.GetUserByUuidParam
	if err := ctx.ShouldBindUri(&params); err != nil {
		utils.ResponseValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	uuidParse, err := uuid.Parse(params.Uuid)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	if err := uh.service.TrashUser(ctx, uuidParse); err != nil {
		utils.ResponseError(ctx, err)
	}
	utils.ResponseStatusCode(ctx, http.StatusOK)
}
