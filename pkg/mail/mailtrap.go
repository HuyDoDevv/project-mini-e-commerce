package mail

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"project-mini-e-commerce/internal/utils"
	"project-mini-e-commerce/pkg/logger"
	"time"

	"github.com/rs/zerolog"
)

type MailtrapConfig struct {
	MailSender     string
	MailSenderName string
	APIKey         string
	URL            string
}
type MailtrapProvider struct {
	client *http.Client
	config *MailtrapConfig
	logger *zerolog.Logger
}

func NewMailtrapProvider(config *MailConfig) (EmailProviderService, error) {
	mailtrapCfg, ok := config.ProviderConfigs["mailtrap"].(map[string]any)
	if !ok {
		return nil, utils.NewError("Invalid or missing Mailtrap configuaration", utils.ErrCodeInternal)
	}

	return &MailtrapProvider{
		client: &http.Client{Timeout: config.Timeout},
		config: &MailtrapConfig{
			MailSender:     mailtrapCfg["mail_sender"].(string),
			MailSenderName: mailtrapCfg["mail_sender_name"].(string),
			APIKey:         mailtrapCfg["api_key"].(string),
			URL:            mailtrapCfg["url"].(string),
		},
		logger: config.Logger,
	}, nil
}

func (m *MailtrapProvider) SendEmail(ctx context.Context, email *Email) error {
	traceId := logger.GetTraceId(ctx)
	startTime := time.Now()

	email.Form = Address{
		Email: m.config.MailSender,
		Name:  m.config.MailSenderName,
	}

	payload, err := json.Marshal(email)
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to marshal email payload")
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", m.config.URL, bytes.NewReader(payload))
	if err != nil {
		m.logger.Error().Err(err).Msg("Failed to create HTTP request")
		return err
	}
	req.Header.Add("Authorization", "Bearer "+m.config.APIKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		m.logger.Error().Str("trace_id", traceId).Dur("duration", time.Since(startTime)).Err(err).Msg("Failed to send email via Mailtrap")
		return utils.WrapError(err, "Failed to send email via Mailtrap", utils.ErrCodeInternal)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		m.logger.Error().Str("trace_id", traceId).Dur("duration", time.Since(startTime)).Int("status", resp.StatusCode).Str("response_body", string(body)).Msg("Mailtrap API returned non-OK status")
		return utils.NewError("Mailtrap API error: "+string(body), utils.ErrCodeInternal)
	}

	return nil
}
