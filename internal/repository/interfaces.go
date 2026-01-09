package repository

import (
	"context"
	"project-mini-e-commerce/internal/db/sqlc"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]sqlc.User, error)
	Create(ctx context.Context, userParams sqlc.CreateUserParams) (sqlc.User, error)
	GetByUUID()
	Update(ctx context.Context, input sqlc.UpdateUserParams) (sqlc.User, error)
	Delete(ctx context.Context, userUuid uuid.UUID) error
	Restore(ctx context.Context, userUuid uuid.UUID) error
	Trash(ctx context.Context, userUuid uuid.UUID) error
	FindUserByEmail()
}
