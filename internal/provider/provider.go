package provider

import (
	"context"

	"github.com/xn3cr0nx/email-service/pkg/model"
)

type Mailer interface {
	Send(context.Context, model.Email) error
	SendBatch(context.Context, model.Email, []string) error
}
