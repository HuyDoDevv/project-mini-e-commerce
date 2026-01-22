package app

import (
	v1handler "project-mini-e-commerce/internal/handler/v1"
	"project-mini-e-commerce/internal/repository"
	"project-mini-e-commerce/internal/routes"
	v1routes "project-mini-e-commerce/internal/routes/v1"
	v1service "project-mini-e-commerce/internal/service/v1"
)

type UserModel struct {
	routes routes.Route
}

func NewUserModel(moduleCtx *ModuleContext) *UserModel {
	userRepo := repository.NewQueryUserRepository(moduleCtx.DB)

	userSer := v1service.NewUserService(userRepo, moduleCtx.Redis)

	userHand := v1handler.NewUserHandler(userSer)

	userRoutes := v1routes.NewUserRoutes(userHand)

	return &UserModel{routes: userRoutes}
}

func (m *UserModel) Routes() routes.Route {
	return m.routes
}
