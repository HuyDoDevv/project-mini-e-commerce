package repository

import "project-mini-e-commerce/internal/models"

type UserRepository interface {
	GetAll() ([]models.User, error)
	Create(user models.User) error
	GetByUUID(uuid string) (models.User, bool)
	Update()
	Delete()
	FindUserByEmail(email string) bool
}
