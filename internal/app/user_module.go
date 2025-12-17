package app

import (
	"project-mini-e-commerce/internal/handler"
	"project-mini-e-commerce/internal/repository"
	"project-mini-e-commerce/internal/routes"
	"project-mini-e-commerce/internal/service"
)

type UserModel struct {
	routes routes.Route
}

func NewUserModel() *UserModel {
	userRepo := repository.NewMemoryUserRepository()

	userSer := service.NewUserService(userRepo)

	userHand := handler.NewUserHandler(userSer)

	userRoutes := routes.NewUserRoutes(userHand)

	return &UserModel{routes: userRoutes}
}

func (m *UserModel) Routes() routes.Route {
	return m.routes
}
