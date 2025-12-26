package app

import (
	"log"
	"project-mini-e-commerce/internal/config"
	"project-mini-e-commerce/internal/routes"
	"project-mini-e-commerce/internal/validation"

	"github.com/gin-gonic/gin"
)

type Module interface {
	Routes() routes.Route
}

type Application struct {
	config *config.Config
	router *gin.Engine
}

func NewApplication(cfg *config.Config) *Application {
	if err := validation.InitValidator(); err != nil {
		log.Fatalf("Failed to initialize validator: %v", err)
	}

	r := gin.Default()

	app := &Application{
		config: cfg,
		router: r,
	}

	app.registerModules()

	return app
}

func (a *Application) registerModules() {
	modules := []Module{
		NewUserModel(),
	}

	var moduleRoutes []routes.Route
	for _, m := range modules {
		moduleRoutes = append(moduleRoutes, m.Routes())
	}

	routes.RegisterRoutes(a.router, moduleRoutes...)
}

func (a *Application) Run() error {
	log.Printf("Server is starting on %s", a.config.ServerAddress)
	return a.router.Run(a.config.ServerAddress)
}
