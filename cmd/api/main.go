package main

import (
	"project-mini-e-commerce/internal/app"
	"project-mini-e-commerce/internal/config"
)

func main() {
	config := config.NewConfig()
	application := app.NewApplication(config)

	if err := application.Run(); err != nil {
		panic(err)
	}
}
