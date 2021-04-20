package task

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hibiken/asynq"
	"github.com/xn3cr0nx/email-service/internal/mailer"
	"github.com/xn3cr0nx/email-service/internal/template"
	"github.com/xn3cr0nx/email-service/pkg/model"
)

type WelcomeEmailTask struct {
	From    string                 `json:"from,omitempty"`
	To      string                 `json:"to,omitempty"`
	Subject string                 `json:"subject,omitempty"`
	Params  WelcomeEmailTaskParams `json:"params,omitempty"`
}

type WelcomeEmailTaskParams struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

func validateWelcomeEmailTask(emailTask WelcomeEmailTask) error {
	if emailTask.From == "" {
		return errInvalidFrom
	}
	if emailTask.To == "" {
		return errInvalidTo
	}
	if emailTask.Subject == "" {
		return errInvalidSubject
	}

	if emailTask.Params.Name == "" || len(emailTask.Params.Name) > 200 {
		return errInvalidName
	}
	if _, err := url.Parse(emailTask.Params.URL); err != nil {
		return errInvalidURL
	}
	return nil
}

type WelcomeEmailHandler struct {
	Mailer mailer.Service
}

func (h WelcomeEmailHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	bytes, err := t.Payload.MarshalJSON()
	if err != nil {
		return err
	}
	var emailTask WelcomeEmailTask
	if err := json.Unmarshal(bytes, &emailTask); err != nil {
		return err
	}

	if err := validateWelcomeEmailTask(emailTask); err != nil {
		return err
	}

	path := template.PathByType(t.Type)
	if path == "" {
		return errTemplateNotFound
	}

	// cache is a singleton, so it is already initialized
	cache, err := template.NewTemplateCache(nil)
	if err != nil {
		return err
	}
	html := string(cache.Get(path))
	filledHtml := fmt.Sprintf(html, emailTask.Params.Name, emailTask.Params.URL)

	if err := h.Mailer.Send(model.Email{
		From:     emailTask.From,
		To:       emailTask.To,
		Subject:  emailTask.Subject,
		HtmlBody: template.FillLayout(filledHtml),
	}); err != nil {
		return err
	}
	return nil
}
