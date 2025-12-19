package app

import (
	"log"
	"project-mini-e-commerce/internal/config"
	"project-mini-e-commerce/internal/routes"
	"project-mini-e-commerce/internal/validation"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Modules interface {
	Routes() routes.Route
}

type Application struct {
	config *config.Config
	router *gin.Engine
}

func NewApplication(config *config.Config) *Application {
	LoadEnv()
	if err := validation.InitValidator(); err != nil {
		log.Fatalf("Validation init faild %v", err)
	}
	r := gin.Default()

	modules := []Modules{
		NewUserModel(),
	}
	routes.RegisterRoutes(r, getModuleRoutes(modules)...)

	return &Application{
		config: config,
		router: r,
	}
}

func getModuleRoutes(modules []Modules) []routes.Route {
	routesList := make([]routes.Route, len(modules))
	for i, module := range modules {
		routesList[i] = module.Routes()
	}
	return routesList
}

func (a *Application) Run() error {
	return a.router.Run(a.config.ServerAddress)
}
func LoadEnv() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("env not found")
	}
}
