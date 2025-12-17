package repository

import (
	"log"
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

func (ur *InmemoryUserRepository) GetAll() {
	log.Println("in repo")
}
func (ur *InmemoryUserRepository) Create() {

}
func (ur *InmemoryUserRepository) GetByUUID() {

}
func (ur *InmemoryUserRepository) Update() {

}
func (ur *InmemoryUserRepository) Delete() {

}
