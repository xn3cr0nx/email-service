package mailer

import "github.com/xn3cr0nx/email-service/pkg/model"

type Service interface {
	Send(model.Email) error
	SendBatch([]model.Email) error
}
