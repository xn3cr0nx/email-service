package sendgrid

import (
	"context"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/xn3cr0nx/email-service/pkg/logger"
	"github.com/xn3cr0nx/email-service/pkg/model"
)

type SendgridClient struct {
	client *sendgrid.Client
}

func NewClient(apiKey string) *SendgridClient {
	c := sendgrid.NewSendClient(apiKey)
	return &SendgridClient{client: c}
}

func (p *SendgridClient) Send(ctx context.Context, email model.Email) error {
	_, err := p.client.SendWithContext(ctx, modelToEmail(email))
	return err
}

func (p *SendgridClient) SendBatch(ctx context.Context, email model.Email, recipients []string) error {
	for _, recipient := range recipients {
		email.To = recipient
		if err := p.Send(ctx, email); err != nil {
			logger.Error("Mailgun provider", err, logger.Params{})
			continue
		}
		time.Sleep(10 * time.Millisecond)
	}
	return nil
}

func modelToEmail(email model.Email) *mail.SGMailV3 {
	from := mail.NewEmail("From", email.From)
	to := mail.NewEmail("To", email.To)
	return mail.NewSingleEmail(from, email.Subject, to, email.TextBody, email.HtmlBody)
}
