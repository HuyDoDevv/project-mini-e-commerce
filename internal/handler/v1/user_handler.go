package v1handler

import (
	"fmt"
	v1service "project-mini-e-commerce/internal/service/v1"
	"time"

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
	fmt.Println("getAllUser start")
	time.Sleep(10 * time.Second)
	fmt.Println("getAllUser end")
}
func (uh *UserHandler) CreateUser(*gin.Context) {

}
func (uh *UserHandler) GetByUserUUID(*gin.Context) {}
func (uh *UserHandler) UpdateUser(*gin.Context) {

}
func (uh *UserHandler) DeleteUser(*gin.Context) {

}
