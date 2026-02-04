package main

import (
	"log"
	"os"
	"path/filepath"
	"project-mini-e-commerce/internal/app"
	"project-mini-e-commerce/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	configFile := config.NewConfig()
	application := app.NewApplication(configFile)

	if err := application.Run(); err != nil {
		panic(err)
	}
}

func loadEnv() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln("unable to get working dir: ", err)
	}
	envPath := filepath.Join(cwd, "/.env")
	if err = godotenv.Load(envPath); err != nil {
		log.Println("No .env file found")
	}
}
