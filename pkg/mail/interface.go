package mail

import "context"

type EmailProviderService interface {
	SendEmail(ctx context.Context, email *Email) error
}
