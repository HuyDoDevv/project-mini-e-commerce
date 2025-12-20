package repository

import "project-mini-e-commerce/internal/models"

type UserRepository interface {
	GetAll() ([]models.User, error)
	Create(user models.User) error
	GetByUUID(uuid string) (models.User, bool)
	Update(uuid string, updateUser models.User) error
	Delete(uuid string) error
	FindUserByEmail(email string) (models.User, bool)
}
