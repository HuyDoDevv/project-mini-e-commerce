package mail

import (
	"project-mini-e-commerce/internal/config"
	"time"

	"github.com/rs/zerolog"
)

type Email struct {
	Form     Address   `json:"form"`
	To       []Address `json:"to"`
	Subject  string    `json:"subject"`
	Text     string    `json:"text"`
	Category string    `json:"category"`
}

type Address struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type MailConfig struct {
	ProviderConfigs map[string]any
	ProviderType    ProviderType
	MaxRetry        int
	Timeout         time.Duration
	Logger          *zerolog.Logger
}

func NewMailService(config *config.Config, logger *zerolog.Logger, providerFactory ProviderFactory) (EmailProviderService, error) {
	return nil, nil
}
