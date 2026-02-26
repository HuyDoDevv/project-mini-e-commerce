package mail

import "log"

type MailtrapProvider struct {
	config *MailConfig
	logger *log.Logger
}

type MailtrapConfig struct{}

func NewMailtrapProvider(config *MailConfig) (EmailProviderService, error) {
	return nil, nil
}
