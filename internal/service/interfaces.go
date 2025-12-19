package service

import (
	"project-mini-e-commerce/internal/models"
)

type UserService interface {
	GetAllUser() ([]models.User, error)
	CreateUser(user models.User) (models.User, error)
	GetByUserUUID(uuid string) (models.User, error)
	UpdateUser()
	DeleteUser()
}
