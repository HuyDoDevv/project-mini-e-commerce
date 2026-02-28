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
	ServerAddress      string
	DB                 DatabaseConfig
	MailProviderType   string
	MailProviderConfig map[string]any
}

func NewConfig() *Config {
	mailProviderConfig := make(map[string]any)

	mailProviderType := utils.GetEnv("MAIL_PROVIDER_TYPE", "mailtrap")
	if mailProviderType == "mailtrap" {
		mailtrapConfig := map[string]any{
			"mail_sender":      utils.GetEnv("MAILTRAP_MAIL_SENDER", "admin@admin.com.vn"),
			"mail_sender_name": utils.GetEnv("MAILTRAP_MAIL_SENDER_NAME", "Admin"),
			"api_key":          utils.GetEnv("MAILTRAP_API_KEY", "f10809b7f6c12394e1f51fc874a21bc0"),
			"url":              utils.GetEnv("MAILTRAP_URL", "https://sandbox.api.mailtrap.io/api/send/4418837"),
		}
		mailProviderConfig[mailProviderType] = mailtrapConfig
	}

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
		MailProviderType:   mailProviderType,
		MailProviderConfig: mailProviderConfig,
	}
}
func (c *Config) DNS() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.DBName, c.DB.DbModeless)
}
