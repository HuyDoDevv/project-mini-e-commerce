package v1handler

import (
	v1service "project-mini-e-commerce/internal/service/v1"

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

func (uh *UserHandler) GetAllUser(ctx *gin.Context) {

}

func (uh *UserHandler) CreateUser(ctx *gin.Context) {

}

func (uh *UserHandler) GetByUserUUID(ctx *gin.Context) {

}
func (uh *UserHandler) UpdateUser(ctx *gin.Context) {

}
func (uh *UserHandler) DeleteUser(ctx *gin.Context) {

}
