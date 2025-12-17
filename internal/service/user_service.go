package service

import (
	"log"
	"project-mini-e-commerce/internal/repository"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (us *userService) GetAllUser() {
	log.Println("service call to repo")
	us.repo.GetAll()
}
func (us *userService) CreateUser() {

}
func (us *userService) GetByUserUUID() {

}
func (us *userService) UpdateUser() {

}
func (us *userService) DeleteUser() {

}
