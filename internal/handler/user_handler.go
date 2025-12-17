package handler

import (
	"project-mini-e-commerce/internal/service"

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

func (uh *UserHandler) GetAllUser(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"status": "ok",
	})
	uh.service.GetAllUser()
}
func (uh *UserHandler) CreateUser(ctx *gin.Context) {

}
func (uh *UserHandler) GetByUserUUID(ctx *gin.Context) {

}
func (uh *UserHandler) UpdateUser(ctx *gin.Context) {

}
func (uh *UserHandler) DeleteUser(ctx *gin.Context) {

}
