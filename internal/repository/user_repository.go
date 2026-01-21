package repository

import (
	"context"
	"fmt"
	"project-mini-e-commerce/internal/db"
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

func (ur *QueryUserRepository) GetAll(ctx context.Context, search, orderBy, sort string, limit, offset int32) ([]sqlc.User, error) {
	var (
		users []sqlc.User
		err   error
	)
	switch {
	case orderBy == "user_id" && sort == "asc":
		users, err = ur.db.GetAllUserIdASC(ctx, sqlc.GetAllUserIdASCParams{
			Limit:  limit,
			Offset: offset,
			Search: search,
		})
	case orderBy == "user_id" && sort == "desc":
		users, err = ur.db.GetAllUserIdDESC(ctx, sqlc.GetAllUserIdDESCParams{
			Limit:  limit,
			Offset: offset,
			Search: search,
		})
	case orderBy == "user_create" && sort == "asc":
		users, err = ur.db.GetAllUserCreateASC(ctx, sqlc.GetAllUserCreateASCParams{
			Search: search,
			Limit:  limit,
			Offset: offset,
		})
	case orderBy == "user_create" && sort == "desc":
		users, err = ur.db.GetAllUserCreateDESC(ctx, sqlc.GetAllUserCreateDESCParams{
			Search: search,
			Limit:  limit,
			Offset: offset,
		})
	}
	if err != nil {
		return []sqlc.User{}, err
	}
	return users, nil
}
func (ur *QueryUserRepository) GetAll2(ctx context.Context, search, orderBy, sort string, limit, offset int32) ([]sqlc.User, error) {
	query := `SELECT user_id, user_uuid, user_email, user_password, user_name, user_age, user_status, user_role, user_deleted_at, user_created_at, user_updated_at
			FROM users
			WHERE user_deleted_at IS NULL
			AND (
				$3::TEXT IS NULL
				OR $3 = '::TEXT'
				OR user_email ILIKE '%' || $3 || '%'
				OR user_name  ILIKE '%' || $3 || '%'
			)`
	order := "ASC"
	if sort == "desc" {
		order = "DESC"
	}

	if orderBy == "user_id" || orderBy == "user_updated_at" {
		query += fmt.Sprintf(" ORDER BY %s %s", orderBy, order)
	} else {
		query += " ORDER BY user_id ASC"
	}
	query += fmt.Sprintf(" LIMIT $1 OFFSET $2")

	rows, err := db.DBPool.Query(ctx, query, limit, offset, search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []sqlc.User
	for rows.Next() {
		var i sqlc.User
		if err := rows.Scan(
			&i.UserID,
			&i.UserUuid,
			&i.UserEmail,
			&i.UserPassword,
			&i.UserName,
			&i.UserAge,
			&i.UserStatus,
			&i.UserRole,
			&i.UserDeletedAt,
			&i.UserCreatedAt,
			&i.UserUpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
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

func (ur *QueryUserRepository) CountAllUsers(ctx context.Context) (int64, error) {
	counterUser, err := ur.db.CountAllUsers(ctx)
	if err != nil {
		return 0, err
	}
	return counterUser, nil
}
