package v1service

import (
	"project-mini-e-commerce/internal/repository"

	"github.com/gin-gonic/gin"
)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (us *userService) GetAllUser(*gin.Context) {}
func (us *userService) CreateUser(*gin.Context) {}
func (us *userService) GetByUserUUID()          {}
func (us *userService) UpdateUser()             {}
func (us *userService) DeleteUser()             {}
