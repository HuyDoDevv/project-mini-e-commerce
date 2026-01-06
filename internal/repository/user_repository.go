package repository

import (
	"context"
	"project-mini-e-commerce/internal/db/sqlc"
)

type QueryUserRepository struct {
	db sqlc.Querier
}

func NewQueryUserRepository(db sqlc.Querier) UserRepository {
	return &QueryUserRepository{
		db: db,
	}
}

func (ur *QueryUserRepository) GetAll() {}
func (ur *QueryUserRepository) Create(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error) {
	user, err := ur.db.CreateUser(ctx, input)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}
func (ur *QueryUserRepository) GetByUUID()       {}
func (ur *QueryUserRepository) Update()          {}
func (ur *QueryUserRepository) Delete()          {}
func (ur *QueryUserRepository) FindUserByEmail() {}
