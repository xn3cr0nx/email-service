package mailgun

import (
	"context"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/xn3cr0nx/email-service/pkg/model"
)

type MailgunClient struct {
	client *mailgun.MailgunImpl
}

func NewClient(domain, apiKey string) *MailgunClient {
	c := mailgun.NewMailgun(domain, apiKey)
	return &MailgunClient{client: c}
}

func (m *MailgunClient) modelToEmail(email model.Email) *mailgun.Message {
	msg := m.client.NewMessage(email.From, email.Subject, email.TextBody, email.To)
	msg.SetHtml(email.HtmlBody)
	return msg
}

func (m *MailgunClient) Send(ctx context.Context, email model.Email) error {
	_, _, err := m.client.Send(ctx, m.modelToEmail(email))
	return err
}

func (m *MailgunClient) SendBatch(ctx context.Context, email model.Email, recipients []string) error {
	msg := m.client.NewMessage(email.From, email.Subject, email.TextBody)
	for _, recipient := range recipients {
		msg.AddRecipient(recipient)
	}
	_, _, err := m.client.Send(ctx, m.modelToEmail(email))
	return err
}
