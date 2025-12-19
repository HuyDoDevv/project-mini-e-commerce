package service

import (
	"project-mini-e-commerce/internal/models"
	"project-mini-e-commerce/internal/repository"
	"project-mini-e-commerce/internal/utils"

	"github.com/google/uuid"
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

func (us *userService) GetAllUser() ([]models.User, error) {
	users, err := us.repo.GetAll()
	if err != nil {
		return []models.User{}, utils.WrapError(err, "can not fetch users", utils.ErrCodeBadRequest)
	}
	return users, nil
}
func (us *userService) CreateUser(user models.User) (models.User, error) {
	user.Email = utils.NormalizeString(user.Email)

	if exitst := us.repo.FindUserByEmail(user.Email); exitst {
		return models.User{}, utils.NewError("email exitst", utils.ErrCodeConflict)
	}

	user.Uuid = uuid.New().String()

	hashePassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, utils.WrapError(err, "can not hash user password", utils.ErrCodeBadRequest)
	}
	user.Password = string(hashePassword)

	if err := us.repo.Create(user); err != nil {
		return models.User{}, utils.WrapError(err, "can not craete user", utils.ErrCodeBadRequest)
	}
	return user, nil
}
func (us *userService) GetByUserUUID(uuid string) (models.User, error) {
	user, exitst := us.repo.GetByUUID(uuid)
	if !exitst {
		return models.User{}, utils.NewError("user not exitst", utils.ErrCodeConflict)
	}
	return user, nil
}

func (us *userService) UpdateUser() {

}
func (us *userService) DeleteUser() {

}
