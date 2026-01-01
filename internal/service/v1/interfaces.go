package v1service

import "github.com/gin-gonic/gin"

type UserService interface {
	GetAllUser(*gin.Context)
	CreateUser(*gin.Context)
	GetByUserUUID()
	UpdateUser()
	DeleteUser()
}
