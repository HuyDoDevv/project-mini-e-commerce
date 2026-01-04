package config

import (
	"fmt"
	"project-mini-e-commerce/internal/utils"
)

type DatabaseConfig struct {
	Host       string
	Port       string
	DBName     string
	User       string
	Password   string
	DbModeless string
}

type Config struct {
	ServerAddress string
	DB            DatabaseConfig
}

func NewConfig() *Config {
	return &Config{
		ServerAddress: utils.GetEnv("SERVER_ADDRESS", ":8080"),
		DB: DatabaseConfig{
			Host:       utils.GetEnv("DB_HOST", "localhost"),
			Port:       utils.GetEnv("DB_PORT", "5432"),
			DBName:     utils.GetEnv("DB_NAME", "myapp"),
			User:       utils.GetEnv("DB_USER", "postgres"),
			Password:   utils.GetEnv("DB_PASSWORD", "postgres"),
			DbModeless: utils.GetEnv("DB_MODELESS", "disable"),
		},
	}
}
func (c *Config) DNS() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.DBName, c.DB.DbModeless)
}
