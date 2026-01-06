package repository

import (
	"context"
	"project-mini-e-commerce/internal/db/sqlc"
)

type UserRepository interface {
	GetAll()
	Create(ctx context.Context, userParams sqlc.CreateUserParams) (sqlc.User, error)
	GetByUUID()
	Update()
	Delete()
	FindUserByEmail()
}
