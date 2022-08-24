package provider

import "github.com/xn3cr0nx/email-service/pkg/model"

type Mailer interface {
	Send(model.Email) error
	SendBatch(model.Email, []string) error
}
