package service

import (
	"project-mini-e-commerce/internal/models"
	"project-mini-e-commerce/internal/repository"
	"project-mini-e-commerce/internal/utils"
	"strings"

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

func (us *userService) GetAllUser(search string, limit, page int) ([]models.User, error) {
	users, err := us.repo.GetAll()
	if err != nil {
		return []models.User{}, utils.WrapError(err, "can not fetch users", utils.ErrCodeBadRequest)
	}

	var filteredUsers []models.User
	if search != "" {
		for _, user := range users {
			if strings.Contains(search, strings.ToLower(user.Email)) || strings.Contains(search, strings.ToLower(user.Name)) {
				filteredUsers = append(filteredUsers, user)
			}
		}
	} else {
		filteredUsers = users
	}

	lengUser := len(filteredUsers)
	start := (page - 1) * limit
	if start >= lengUser {
		return []models.User{}, nil
	}
	end := start + limit

	if end > lengUser {
		end = lengUser
	}
	return filteredUsers[start:end], nil
}
func (us *userService) CreateUser(user models.User) (models.User, error) {
	user.Email = utils.NormalizeString(user.Email)

	if _, exitst := us.repo.FindUserByEmail(user.Email); exitst {
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

func (us *userService) UpdateUser(uuid string, updateUser models.User) (models.User, error) {
	updateUser.Email = utils.NormalizeString(updateUser.Email)

	if u, exitst := us.repo.FindUserByEmail(updateUser.Email); exitst && u.Uuid != uuid {
		return models.User{}, utils.NewError("email already exist", utils.ErrCodeConflict)
	}

	currentUser, exitst := us.repo.GetByUUID(uuid)
	if !exitst {
		return models.User{}, utils.NewError("user not exitst", utils.ErrCodeConflict)
	}

	if updateUser.Name != "" {
		currentUser.Name = updateUser.Name
	}

	if updateUser.Email != "" {
		currentUser.Email = updateUser.Email
	}
	if updateUser.Age != 0 {
		currentUser.Age = updateUser.Age
	}

	if updateUser.Status != 0 {
		currentUser.Status = updateUser.Status
	}

	if updateUser.Level != 0 {
		currentUser.Level = updateUser.Level
	}

	if updateUser.Password != "" {
		hashePassword, err := bcrypt.GenerateFromPassword([]byte(currentUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return models.User{}, utils.WrapError(err, "user not found", utils.ErrCodeBadRequest)
		}
		currentUser.Password = string(hashePassword)
	}
	if err := us.repo.Update(uuid, currentUser); err != nil {
		return models.User{}, utils.WrapError(err, "fail to update user", utils.ErrCodeBadRequest)

	}
	return currentUser, nil
}
func (us *userService) DeleteUser(uuid string) error {
	if err := us.repo.Delete(uuid); err != nil {
		return err
	}
	return nil
}
