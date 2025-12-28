package config

import (
	"log"
	"project-mini-e-commerce/internal/utils"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress string
}

func NewConfig() *Config {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	serverAddress := utils.GetEnv("SERVER_PORT", ":8080")

	return &Config{
		ServerAddress: serverAddress,
	}
}
