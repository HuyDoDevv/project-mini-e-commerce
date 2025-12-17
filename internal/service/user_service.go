package service

import "project-mini-e-commerce/internal/repository"

type UserService struct {
	repo *repository.MemoryUserRepository
}

func NewUserService(repo *repository.MemoryUserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}
