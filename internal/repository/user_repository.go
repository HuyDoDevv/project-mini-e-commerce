package repository

import (
	"fmt"
	"project-mini-e-commerce/internal/models"
	"slices"
)

type InmemoryUserRepository struct {
	users []models.User
}

func NewMemoryUserRepository() UserRepository {
	return &InmemoryUserRepository{
		users: make([]models.User, 0),
	}
}

func (ur *InmemoryUserRepository) GetAll() ([]models.User, error) {
	return ur.users, nil
}
func (ur *InmemoryUserRepository) Create(user models.User) error {
	ur.users = append(ur.users, user)
	return nil
}
func (ur *InmemoryUserRepository) GetByUUID(uuid string) (models.User, bool) {
	for _, user := range ur.users {
		if user.Uuid == uuid {
			return user, true
		}
	}
	return models.User{}, false
}
func (ur *InmemoryUserRepository) Update(uuid string, updateUser models.User) error {
	for i, user := range ur.users {
		if user.Uuid == uuid {
			ur.users[i] = updateUser
		}
	}
	return nil
}
func (ur *InmemoryUserRepository) Delete(uuid string) error {
	for i, user := range ur.users {
		if user.Uuid == uuid {
			ur.users = slices.Delete(ur.users, i, i+1)
			return nil
		}
	}
	return fmt.Errorf("user not found")
}
func (ur *InmemoryUserRepository) FindUserByEmail(email string) (models.User, bool) {
	for _, user := range ur.users {
		if user.Email == email {
			return user, true
		}
	}
	return models.User{}, false
}
