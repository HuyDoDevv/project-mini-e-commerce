package repository

import (
	"project-mini-e-commerce/internal/models"
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
func (ur *InmemoryUserRepository) Update() {

}
func (ur *InmemoryUserRepository) Delete() {

}
func (ur *InmemoryUserRepository) FindUserByEmail(email string) bool {
	for _, user := range ur.users {
		if user.Email == email {
			return true
		}
	}
	return false
}
