package mail

import (
	"fmt"
	"project-mini-e-commerce/internal/utils"
)

type ProviderType string

const (
	ProviderMailtrap ProviderType = "mailtrap"
)

type ProviderFactory interface {
	CreateProvider(config *MailConfig) (EmailProviderService, error)
}

type MailTrapProviderFactory struct {
}

func (f *MailTrapProviderFactory) CreateProvider(config *MailConfig) (EmailProviderService, error) {
	return NewMailTrapProvider(config)
}

func NewProviderFactory(providerType ProviderType) (ProviderFactory, error) {
	switch providerType {
	case ProviderMailtrap:
		return &MailTrapProviderFactory{}, nil
	default:
		return nil, utils.NewError(fmt.Sprintf("Unsuported provider type: %s", utils.ErrorCode(providerType)), utils.ErrCodeInternal)
	}
}
