package repository

import (
	"context"
	"project-mini-e-commerce/internal/db/sqlc"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetAll(ctx context.Context, search, orderBy, sort string, limit, offset int32) ([]sqlc.User, error)
	GetAll2(ctx context.Context, search, orderBy, sort string, limit, offset int32, deleted bool) ([]sqlc.User, error)
	Create(ctx context.Context, userParams sqlc.CreateUserParams) (sqlc.User, error)
	FindUUID(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error)
	Update(ctx context.Context, input sqlc.UpdateUserParams) (sqlc.User, error)
	Delete(ctx context.Context, userUuid uuid.UUID) error
	Restore(ctx context.Context, userUuid uuid.UUID) error
	Trash(ctx context.Context, userUuid uuid.UUID) error
	FindUserByEmail()
	CountAllUsers(ctx context.Context) (int64, error)
}
