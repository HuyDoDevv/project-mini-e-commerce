package v1service

import (
	"errors"
	"fmt"
	"project-mini-e-commerce/internal/db/sqlc"
	"project-mini-e-commerce/internal/repository"
	"project-mini-e-commerce/internal/utils"
	"project-mini-e-commerce/pkg/cache"
	"project-mini-e-commerce/pkg/logger"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo  repository.UserRepository
	cache cache.RedisCacheService
}

func NewUserService(repo repository.UserRepository, redisClient *redis.Client) UserService {
	return &userService{
		repo:  repo,
		cache: cache.NewRedisCacheService(redisClient),
	}
}

func (us *userService) GetUserByUUID(ctx *gin.Context, userUuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()
	user, err := us.repo.FindUUID(context, userUuid)
	if err != nil {
		utils.ResponseError(ctx, err)
	}
	return user, nil
}

func (us *userService) GetAllUser(ctx *gin.Context, search, orderBy, sort string, limit, page int32, deleted bool) ([]sqlc.User, int64, error) {
	context := ctx.Request.Context()
	var cacheData struct {
		Users []sqlc.User
		Total int64
	}
	cacheKey := us.generateCacheKey(search, orderBy, sort, limit, page, deleted)
	if err := us.cache.Get(cacheKey, &cacheData); err == nil && cacheData.Users != nil {
		return cacheData.Users, cacheData.Total, nil
	}
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

	users, err := us.repo.GetAll2(context, search, orderBy, sort, limit, offset, deleted)
	if err != nil {
		return []sqlc.User{}, 0, utils.NewError("cannot fetch users", utils.ErrCodeInternal)
	}

	countUser, err := us.repo.CountAllUsers(ctx)
	if err != nil {
		return []sqlc.User{}, 0, utils.WrapError(err, "failed to count all users", utils.ErrCodeInternal)
	}

	cacheData = struct {
		Users []sqlc.User
		Total int64
	}{
		Users: users,
		Total: countUser,
	}

	if err := us.cache.Set(cacheKey, cacheData, 5*time.Minute); err != nil {
		return []sqlc.User{}, 0, utils.WrapError(err, "failed to set cache", utils.ErrCodeInternal)
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

	if err := us.cache.Clear("users:*"); err != nil {
		logger.Logger.Warn().Err(err).Msg("Failed to clear cache after creating user")
	}

	return user, nil
}

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

func (us *userService) generateCacheKey(search, orderBy, sort string, limit, page int32, deleted bool) string {
	search = utils.NormalizeString(search)
	if search == "" {
		search = "none"
	}
	orderBy = strings.TrimSpace(orderBy)
	if orderBy == "" {
		orderBy = "user_created_at"
	}
	sort = strings.ToLower(strings.TrimSpace(sort))
	if sort == "" {
		sort = "desc"
	}
	return fmt.Sprintf("users:%s-%s-%s-%d-%d-%v", search, orderBy, sort, limit, page, deleted)
}
