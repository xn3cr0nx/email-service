package email

import (
	"context"
	"fmt"
	"net/url"

	"github.com/xn3cr0nx/email-service/internal/environment"
	"github.com/xn3cr0nx/email-service/internal/provider"
	"github.com/xn3cr0nx/email-service/internal/template"
	"github.com/xn3cr0nx/email-service/pkg/model"
)

type VerificationEmailBody struct {
	From    string                  `json:"from,omitempty"`
	To      string                  `json:"to,omitempty"`
	Subject string                  `json:"subject,omitempty"`
	Params  VerificationEmailParams `json:"params,omitempty"`
}

type VerificationEmailParams struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

func (b *VerificationEmailBody) ValidateBody() error {
	if b.From == "" {
		b.From = environment.Get().Sender
	}
	if b.To == "" {
		return errInvalidTo
	}
	if b.Subject == "" {
		return errInvalidSubject
	}

	if b.Params.Name == "" || len(b.Params.Name) > 200 {
		return errInvalidName
	}
	if _, err := url.Parse(b.Params.URL); err != nil {
		return errInvalidURL
	}
	return nil
}

func (b *VerificationEmailBody) Process(ctx context.Context, m provider.Mailer) error {
	if err := b.ValidateBody(); err != nil {
		return err
	}

	path := template.PathByType(template.VerificationEmail)
	if path == "" {
		return errTemplateNotFound
	}

	// cache is a singleton, so it is already initialized
	cache, err := template.NewTemplateCache(nil)
	if err != nil {
		return err
	}
	html := string(cache.Get(path))
	filledHtml := fmt.Sprintf(html, b.Params.Name, b.Params.URL, b.Params.URL)

	if err := m.Send(ctx, model.Email{
		From:     b.From,
		To:       b.To,
		Subject:  b.Subject,
		HtmlBody: template.FillLayout(filledHtml),
	}); err != nil {
		return err
	}
	return nil
}
