package repository

import "project-mini-e-commerce/internal/models"

type MemoryUserRepository struct {
	users []models.User
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: make([]models.User, 0),
	}
}
