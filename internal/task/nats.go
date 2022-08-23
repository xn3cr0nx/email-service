package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/xn3cr0nx/email-service/internal/email"
	"github.com/xn3cr0nx/email-service/internal/environment"
	"github.com/xn3cr0nx/email-service/internal/mailer"
	"github.com/xn3cr0nx/email-service/internal/template"
	"github.com/xn3cr0nx/email-service/pkg/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/trace"
)

type NatsEmailConsumer struct {
	Subscriber *nats.Conn
	Mailer     mailer.Service
	tracer     trace.Tracer
	meter      metric.Meter
}

type NatsMessage struct {
	Key   string
	Value []byte
}

func NewNatsEmailConsumer(n *nats.Conn, m mailer.Service, tracer trace.Tracer, meter metric.Meter) *NatsEmailConsumer {
	return &NatsEmailConsumer{n, m, tracer, meter}
}

func (n *NatsEmailConsumer) Run(ctx context.Context) error {
	emailCounterLock := new(sync.RWMutex)
	var emailCounter syncfloat64.Counter
	welcomeEmailCounterLock := new(sync.RWMutex)
	var welcomeEmailCounter syncfloat64.Counter
	if n.meter != nil {
		var err error
		emailCounter, err = metric.NewNoopMeter().SyncFloat64().Counter("NATS.emails")
		if err != nil {
			return err
		}
		welcomeEmailCounter, err = metric.NewNoopMeter().SyncFloat64().Counter("NATS.emails.welcome")
		if err != nil {
			return err
		}
	}

	var spanContext context.Context
	var span trace.Span
	if n.tracer != nil {
		spanContext, span = (n.tracer).Start(ctx, "email")
		defer span.End()
	} else {
		spanContext = context.WithValue(ctx, "email", "")
	}

	// TODO: inject the subscribe subject from env
	sub, err := n.Subscriber.SubscribeSync(environment.Get().NatsSubject)
	if err != nil {
		return err
	}

	for {
		var msg NatsMessage
		m, err := sub.NextMsg(5 * time.Second)
		if err != nil {
			logger.Debug("Email Service NATS", err.Error(), logger.Params{})
			continue
		}
		logger.Info("Email Service NATS", "Received NATS msg", logger.Params{})
		if err := m.Ack(); err != nil {
			logger.Debug("Email Service NATS", err.Error(), logger.Params{})
		}
		if err := json.Unmarshal(m.Data, &msg); err != nil {
			logger.Error("Email Service NATS", err, logger.Params{})
		}

		var emailSpanContext context.Context
		var emailSpan trace.Span
		if n.tracer != nil {
			emailSpanContext, emailSpan = (n.tracer).Start(spanContext, string(msg.Key))
			emailSpan.SetAttributes(attribute.Key(msg.Key).String(string(msg.Value)))
		} else {
			emailSpanContext = context.WithValue(spanContext, string(msg.Key), string(msg.Value))
		}

		switch string(msg.Key) {
		case template.WelcomeEmail:
			var emailTask email.WelcomeEmailBody
			if err := json.Unmarshal(msg.Value, &emailTask); err != nil {
				logger.Error("Email Service NATS", fmt.Errorf("could not unmarshal message: %v", err), logger.Params{})
				continue
			}

			logger.Info("Email Service NATS", "Received WelcomeEmail", logger.Params{"subject": emailTask.Subject, "to": emailTask.To, "from": emailTask.From})

			if err := emailTask.Process(emailSpanContext, n.Mailer); err != nil {
				logger.Error("Email Service NATS", fmt.Errorf("could not process message: %v", err), logger.Params{})
				continue
			}
			if n.meter != nil {
				(*welcomeEmailCounterLock).Lock()
				welcomeEmailCounter.Add(emailSpanContext, 1)
				(*welcomeEmailCounterLock).Unlock()
			}

		case template.ResetEmail:
			var emailTask email.ResetEmailBody
			if err := json.Unmarshal(msg.Value, &emailTask); err != nil {
				logger.Error("Email Service NATS", fmt.Errorf("could not unmarshal message: %v", err), logger.Params{})
				continue
			}

			if err := emailTask.Process(emailSpanContext, n.Mailer); err != nil {
				logger.Error("Email Service NATS", fmt.Errorf("could not process message: %v", err), logger.Params{})
				continue
			}

		default:
			logger.Error("Email Service NATS", errors.New("unmatched case"), logger.Params{})
		}

		if n.meter != nil {
			(*emailCounterLock).Lock()
			emailCounter.Add(emailSpanContext, 1)
			(*emailCounterLock).Unlock()
		}
	}
}
