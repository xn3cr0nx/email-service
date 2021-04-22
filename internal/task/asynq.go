package task

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/xn3cr0nx/email-service/internal/email"
	"github.com/xn3cr0nx/email-service/internal/mailer"
	"github.com/xn3cr0nx/email-service/internal/template"
	"github.com/xn3cr0nx/email-service/pkg/logger"
)

type EmailHandler struct {
	Mailer mailer.Service
}

func NewEmailHandler(m mailer.Service) *EmailHandler {
	return &EmailHandler{m}
}

func (h EmailHandler) ProcessTask(ctx context.Context, t *asynq.Task) (err error) {
	start := time.Now()
	logger.Info("Email Service Queue", fmt.Sprintf("Start processing"), logger.Params{"type": t.Type})
	defer (func() {
		logger.Info("Email Service Queue", fmt.Sprintf("Finished processing. Elapsed Time = %v", time.Since(start)), logger.Params{"type": t.Type})
	})()

	switch t.Type {
	case template.WelcomeEmail:
		bytes, e := t.Payload.MarshalJSON()
		if e != nil {
			logger.Error("Email Service Queue", err, logger.Params{"type": t.Type})
			return e
		}
		var emailTask email.WelcomeEmailBody
		if err = json.Unmarshal(bytes, &emailTask); err != nil {
			logger.Error("Email Service Queue", err, logger.Params{"type": t.Type})
			return
		}

		if err = emailTask.Process(h.Mailer); err != nil {
			logger.Error("Email Service Queue", err, logger.Params{"type": t.Type})
			return
		}
	}
	return
}
