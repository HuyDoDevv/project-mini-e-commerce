package v1service

import (
	"project-mini-e-commerce/internal/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserService interface {
	GetAllUser(ctx *gin.Context, search, orderBy, sort string, limit, offset int32) ([]sqlc.User, int64, error)
	CreateUser(ctx *gin.Context, userParams sqlc.CreateUserParams) (sqlc.User, error)
	GetByUserUUID()
	UpdateUser(ctx *gin.Context, userParams sqlc.UpdateUserParams) (sqlc.User, error)
	DeleteUser(ctx *gin.Context, userUuid uuid.UUID) error
	RestoreUser(ctx *gin.Context, userUuid uuid.UUID) error
	TrashUser(ctx *gin.Context, userUuid uuid.UUID) error
}
