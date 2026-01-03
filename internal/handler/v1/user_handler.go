package v1handler

import (
	"fmt"
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

func (uh *UserHandler) GetAllUser(ctx *gin.Context) {}
func (uh *UserHandler) CreateUser(*gin.Context) {

}
func (uh *UserHandler) GetByUserUUID(*gin.Context) {}
func (uh *UserHandler) UpdateUser(*gin.Context) {

}
func (uh *UserHandler) DeleteUser(*gin.Context) {

}
func (uh *UserHandler) PanicUser(*gin.Context) {
	var a []int
	fmt.Println(a[1])
}
