package repository

import (
	"context"
	"fmt"
	"project-mini-e-commerce/internal/db/sqlc"

	"github.com/google/uuid"
)

type QueryUserRepository struct {
	db sqlc.Querier
}

func NewQueryUserRepository(db sqlc.Querier) UserRepository {
	return &QueryUserRepository{
		db: db,
	}
}

func (ur *QueryUserRepository) GetAll(ctx context.Context) ([]sqlc.User, error) {
	users, err := ur.db.GetAllUsers(ctx)
	if err != nil {
		return []sqlc.User{}, err
	}
	return users, nil
}
func (ur *QueryUserRepository) Create(ctx context.Context, input sqlc.CreateUserParams) (sqlc.User, error) {
	user, err := ur.db.CreateUser(ctx, input)
	if err != nil {
		return sqlc.User{}, err
	}

	return user, nil
}
func (ur *QueryUserRepository) GetByUUID() {}
func (ur *QueryUserRepository) Update(ctx context.Context, input sqlc.UpdateUserParams) (sqlc.User, error) {
	user, err := ur.db.UpdateUser(ctx, input)
	if err != nil {
		return sqlc.User{}, err
	}
	return user, nil
}
func (ur *QueryUserRepository) FindUserByEmail() {}
func (ur *QueryUserRepository) Delete(ctx context.Context, userUuid uuid.UUID) error {
	rowsAffected, err := ur.db.DeleteUser(ctx, userUuid)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found or already deleted")
	}
	return nil
}
func (ur *QueryUserRepository) Restore(ctx context.Context, userUuid uuid.UUID) error {
	rowsAffected, err := ur.db.RestoreUser(ctx, userUuid)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found or already deleted")
	}
	return nil
}
func (ur *QueryUserRepository) Trash(ctx context.Context, userUuid uuid.UUID) error {
	rowsAffected, err := ur.db.TrashUser(ctx, userUuid)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("user not found or already deleted")
	}
	return nil
}
