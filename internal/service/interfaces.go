package service

import (
	"project-mini-e-commerce/internal/models"
)

type UserService interface {
	GetAllUser(search string, limit, page int) ([]models.User, error)
	CreateUser(user models.User) (models.User, error)
	GetByUserUUID(uuid string) (models.User, error)
	UpdateUser(uuid string, updateUser models.User) (models.User, error)
	DeleteUser(uuid string) error
}
