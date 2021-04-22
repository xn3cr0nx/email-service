package task

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/hibiken/asynq"
	"github.com/xn3cr0nx/email-service/internal/email"
	"github.com/xn3cr0nx/email-service/internal/mailer"
	"github.com/xn3cr0nx/email-service/internal/template"
	"github.com/xn3cr0nx/email-service/pkg/logger"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type EmailHandler struct {
	Mailer mailer.Service
	tracer *trace.Tracer

	meter                   *metric.Meter
	emailCounter            metric.Int64Counter
	emailCounterLock        *sync.RWMutex
	welcomeEmailCounter     metric.Int64Counter
	welcomeEmailCounterLock *sync.RWMutex
}

func NewEmailHandler(m mailer.Service, tracer *trace.Tracer, meter *metric.Meter) *EmailHandler {
	emailCounterLock := new(sync.RWMutex)
	var emailCounter metric.Int64Counter
	welcomeEmailCounterLock := new(sync.RWMutex)
	var welcomeEmailCounter metric.Int64Counter
	if meter != nil {
		emailCounter = metric.Must(*meter).NewInt64Counter("kafka.emails")
		welcomeEmailCounter = metric.Must(*meter).NewInt64Counter("kafka.emails.welcome")
	}

	return &EmailHandler{m, tracer, meter, emailCounter, emailCounterLock, welcomeEmailCounter, welcomeEmailCounterLock}
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

		if err = emailTask.Process(ctx, h.Mailer); err != nil {
			logger.Error("Email Service Queue", err, logger.Params{"type": t.Type})
			return
		}
		if h.meter != nil {
			(*h.welcomeEmailCounterLock).Lock()
			h.welcomeEmailCounter.Add(ctx, 1)
			(*h.welcomeEmailCounterLock).Unlock()
		}
	}
	if h.meter != nil {
		(*h.emailCounterLock).Lock()
		h.emailCounter.Add(ctx, 1)
		(*h.emailCounterLock).Unlock()
	}
	return
}
