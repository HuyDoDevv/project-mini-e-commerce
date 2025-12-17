package main

import (
	"project-mini-e-commerce/internal/config"
	"project-mini-e-commerce/internal/handler"
	"project-mini-e-commerce/internal/repository"
	"project-mini-e-commerce/internal/routes"
	"project-mini-e-commerce/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	config := config.NewConfig()

	userRepo := repository.NewMemoryUserRepository()

	userSer := service.NewUserService(userRepo)

	userHand := handler.NewUserHandler(userSer)

	userRoutes := routes.NewUserRoutes(userHand)

	r := gin.Default()

	routes.RegisterRoutes(r, userRoutes)
	if err := r.Run(config.ServerAddress); err != nil {
		panic("Cannot run server")
	}
}
