package main

import (
	"project-mini-e-commerce/internal/app"
	"project-mini-e-commerce/internal/config"
)

func main() {
	configFile := config.NewConfig()
	application := app.NewApplication(configFile)

	if err := application.Run(); err != nil {
		panic(err)
	}
}
