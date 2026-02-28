package mail

import (
	"context"
	"project-mini-e-commerce/internal/config"
	"project-mini-e-commerce/internal/utils"
	"project-mini-e-commerce/pkg/logger"
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

type MailService struct {
	config   *MailConfig
	provider EmailProviderService
	logger   *zerolog.Logger
}

func NewMailService(config *config.Config, logger *zerolog.Logger, providerFactory ProviderFactory) (EmailProviderService, error) {
	cfg := &MailConfig{
		ProviderConfigs: config.MailProviderConfig,
		ProviderType:    ProviderType(config.MailProviderType),
		MaxRetry:        3,
		Timeout:         10 * time.Second,
		Logger:          logger,
	}

	provider, err := providerFactory.CreateProvider(cfg)
	if err != nil {
		return nil, utils.WrapError(err, "Failed to create email provider", utils.ErrCodeInternal)
	}

	return &MailService{
		config:   cfg,
		provider: provider,
		logger:   logger,
	}, nil
}

func (s *MailService) SendEmail(ctx context.Context, email *Email) error {
	traceId := logger.GetTraceId(ctx)
	startTime := time.Now()

	var lastErr error
	for attempt := 1; attempt <= s.config.MaxRetry; attempt++ {
		startAttemptTime := time.Now()
		err := s.provider.SendEmail(ctx, email)

		if err == nil {
			s.logger.Error().Str("trace_id", traceId).Dur("duration", time.Since(startAttemptTime)).Str("operation", "send_email").Interface("to", email.To).Str("subject", email.Subject).Str("category", email.Category).Msg("Email sent successfully")
			return nil
		}

		lastErr = err
		s.logger.Warn().Err(err).Str("trace_id", traceId).Dur("duration", time.Since(startAttemptTime)).Int("attempt", attempt).Msg("Failed to send email, retrying...")
		time.Sleep(time.Duration(attempt) * time.Second)
	}

	s.logger.Error().Str("trace_id", traceId).Err(lastErr).Dur("duration", time.Since(startTime)).Int("attempts", s.config.MaxRetry).Msg("Failed to send email after retries")
	return utils.WrapError(lastErr, "Failed to send email after retries", utils.ErrCodeInternal)
}
