package postmark

import (
	client "github.com/keighl/postmark"
	"github.com/xn3cr0nx/email-service/pkg/model"
)

type PostmarkClient struct {
	*client.Client
}

func NewClient(serverToken, accountToken string) *PostmarkClient {
	c := client.NewClient(serverToken, accountToken)
	return &PostmarkClient{c}
}

func (p *PostmarkClient) Send(email model.Email) error {
	_, err := p.SendEmail(modelToEmail(email))
	if err != nil {
		return err
	}
	return nil
}

func (p *PostmarkClient) SendBatch(email model.Email, recipients []string) error {
	clientEmails := make([]client.Email, len(recipients))
	for i, recipient := range recipients {
		model := modelToEmail(email)
		model.To = recipient
		clientEmails[i] = model
	}
	_, err := p.SendEmailBatch(clientEmails)
	if err != nil {
		return err
	}
	return nil
}

func modelToEmail(email model.Email) client.Email {
	return client.Email{
		From:     email.From,
		To:       email.To,
		Subject:  email.Subject,
		HtmlBody: email.HtmlBody,
		TextBody: email.TextBody,
		Tag:      email.Tag,
	}
}
