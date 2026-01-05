package v1service

import (
	"project-mini-e-commerce/internal/db/sqlc"

	"github.com/gin-gonic/gin"
)

type UserService interface {
	GetAllUser(*gin.Context)
	CreateUser(ctx *gin.Context, userParams sqlc.CreateUserParams) (sqlc.User, error)
	GetByUserUUID()
	UpdateUser()
	DeleteUser()
}
