package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hibiken/asynq"
	"github.com/xn3cr0nx/email-service/internal/email"
	"github.com/xn3cr0nx/email-service/internal/mailer"
	"github.com/xn3cr0nx/email-service/internal/template"
	"github.com/xn3cr0nx/email-service/pkg/logger"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/trace"
)

type EmailHandler struct {
	Mailer mailer.Service
	tracer trace.Tracer

	meter                   metric.Meter
	emailCounter            syncfloat64.Counter
	emailCounterLock        *sync.RWMutex
	welcomeEmailCounter     syncfloat64.Counter
	welcomeEmailCounterLock *sync.RWMutex
}

func NewEmailHandler(m mailer.Service, tracer trace.Tracer, meter metric.Meter) *EmailHandler {
	emailCounterLock := new(sync.RWMutex)
	var emailCounter syncfloat64.Counter
	welcomeEmailCounterLock := new(sync.RWMutex)
	var welcomeEmailCounter syncfloat64.Counter
	if meter != nil {
		var err error
		emailCounter, err = metric.NewNoopMeter().SyncFloat64().Counter("kafka.emails")
		if err != nil {
			return nil
		}
		welcomeEmailCounter, err = metric.NewNoopMeter().SyncFloat64().Counter("kafka.emails.welcome")
		if err != nil {
			return nil
		}
	}

	return &EmailHandler{m, tracer, meter, emailCounter, emailCounterLock, welcomeEmailCounter, welcomeEmailCounterLock}
}

func (h EmailHandler) ProcessTask(ctx context.Context, t *asynq.Task) (err error) {
	start := time.Now()
	logger.Info("Email Service Queue", "Start processing", logger.Params{"type": t.Type})
	defer (func() {
		logger.Info("Email Service Queue", fmt.Sprintf("Finished processing. Elapsed Time = %v", time.Since(start)), logger.Params{"type": t.Type})
	})()

	switch t.Type() {
	case template.WelcomeEmail:
		var emailTask email.WelcomeEmailBody
		if err = json.Unmarshal(t.Payload(), &emailTask); err != nil {
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

	case template.VerificationEmail:
		var emailTask email.VerificationEmailBody
		if err = json.Unmarshal(t.Payload(), &emailTask); err != nil {
			logger.Error("Email Service Queue", err, logger.Params{"type": t.Type()})
			return
		}

		if err = emailTask.Process(ctx, h.Mailer); err != nil {
			logger.Error("Email Service Queue", err, logger.Params{"type": t.Type()})
			return
		}

	case template.ResetEmail:
		var emailTask email.ResetEmailBody
		if err = json.Unmarshal(t.Payload(), &emailTask); err != nil {
			logger.Error("Email Service Queue", err, logger.Params{"type": t.Type()})
			return
		}

		if err = emailTask.Process(ctx, h.Mailer); err != nil {
			logger.Error("Email Service Queue", err, logger.Params{"type": t.Type()})
			return
		}

	default:
		return errors.New("unmatched case")
	}
	if h.meter != nil {
		(*h.emailCounterLock).Lock()
		h.emailCounter.Add(ctx, 1)
		(*h.emailCounterLock).Unlock()
	}
	return
}
