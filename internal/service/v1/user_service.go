package v1service

import (
	"errors"
	"project-mini-e-commerce/internal/db/sqlc"
	"project-mini-e-commerce/internal/repository"
	"project-mini-e-commerce/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (us *userService) GetAllUser(ctx *gin.Context, search, orderBy, sort string, limit, page int32) ([]sqlc.User, int64, error) {
	context := ctx.Request.Context()
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	if orderBy == "" {
		orderBy = "user_id"
	}
	if sort == "" {
		sort = "asc"
	}
	offset := (page - 1) * limit
	search = utils.NormalizeString(search)
	users, err := us.repo.GetAll(context, search, orderBy, sort, limit, offset)
	if err != nil {
		return []sqlc.User{}, 0, utils.NewError("cannot fetch users", utils.ErrCodeInternal)
	}

	countUser, err := us.repo.CountAllUsers(ctx)
	if err != nil {
		return []sqlc.User{}, 0, utils.WrapError(err, "failed to count all users", utils.ErrCodeInternal)
	}
	return users, countUser, nil
}
func (us *userService) CreateUser(ctx *gin.Context, input sqlc.CreateUserParams) (sqlc.User, error) {
	context := ctx.Request.Context()

	input.UserEmail = utils.NormalizeString(input.UserEmail)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return sqlc.User{}, utils.WrapError(err, "failed to hash password", utils.ErrCodeInternal)
	}

	input.UserPassword = string(hashedPassword)

	user, err := us.repo.Create(context, input)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return sqlc.User{}, utils.NewError("email already exist", utils.ErrCodeConflict)
		}

		return sqlc.User{}, utils.WrapError(err, "failed to create a new user", utils.ErrCodeInternal)
	}

	return user, nil
}
func (us *userService) GetByUserUUID() {}
func (us *userService) UpdateUser(ctx *gin.Context, input sqlc.UpdateUserParams) (sqlc.User, error) {
	context := ctx.Request.Context()
	if input.UserPassword != nil && *input.UserPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.UserPassword), bcrypt.DefaultCost)
		if err != nil {
			return sqlc.User{}, utils.WrapError(err, "failed to hash password", utils.ErrCodeInternal)
		}
		hasher := string(hashedPassword)
		input.UserPassword = &hasher
	}

	user, err := us.repo.Update(context, input)
	if err != nil {
		return sqlc.User{}, utils.WrapError(err, "failed to update user", utils.ErrCodeInternal)
	}
	return user, nil
}
func (us *userService) DeleteUser(ctx *gin.Context, userUuid uuid.UUID) error {
	context := ctx.Request.Context()

	err := us.repo.Delete(context, userUuid)
	if err != nil {
		return utils.WrapError(err, "failed to delete user", utils.ErrCodeInternal)
	}
	return nil
}
func (us *userService) RestoreUser(ctx *gin.Context, userUuid uuid.UUID) error {
	context := ctx.Request.Context()

	if err := us.repo.Restore(context, userUuid); err != nil {
		return utils.WrapError(err, "failed to restore user", utils.ErrCodeInternal)
	}
	return nil
}
func (us *userService) TrashUser(ctx *gin.Context, userUuid uuid.UUID) error {
	context := ctx.Request.Context()
	if err := us.repo.Trash(context, userUuid); err != nil {
		return utils.WrapError(err, "failed to trash user", utils.ErrCodeInternal)
	}
	return nil
}
