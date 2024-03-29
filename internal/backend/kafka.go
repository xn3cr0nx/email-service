package backend

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/xn3cr0nx/email-service/internal/email"
	"github.com/xn3cr0nx/email-service/internal/provider"
	"github.com/xn3cr0nx/email-service/internal/template"
	"github.com/xn3cr0nx/email-service/pkg/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/trace"
)

type KafkaEmailConsumer struct {
	Reader *kafka.Reader
	Mailer provider.Mailer
	tracer trace.Tracer
	meter  metric.Meter
}

func NewKafkaEmailConsumer(k *kafka.Reader, m provider.Mailer, tracer trace.Tracer, meter metric.Meter) *KafkaEmailConsumer {
	return &KafkaEmailConsumer{k, m, tracer, meter}
}

func (k *KafkaEmailConsumer) Run(ctx context.Context) error {
	emailCounterLock := new(sync.RWMutex)
	var emailCounter syncfloat64.Counter
	welcomeEmailCounterLock := new(sync.RWMutex)
	var welcomeEmailCounter syncfloat64.Counter
	if k.meter != nil {
		var err error
		emailCounter, err = metric.NewNoopMeter().SyncFloat64().Counter("kafka.emails")
		if err != nil {
			return err
		}
		welcomeEmailCounter, err = metric.NewNoopMeter().SyncFloat64().Counter("kafka.emails.welcome")
		if err != nil {
			return err
		}
	}

	var spanContext context.Context
	var span trace.Span
	if k.tracer != nil {
		spanContext, span = (k.tracer).Start(ctx, "email")
		defer span.End()
	} else {
		spanContext = context.WithValue(ctx, "email", "")
	}

	// TODO: consider if adding unique UUID to message keys

	for {
		// the `FetchMessage` method blocks until we receive the next event, and the message needs to
		// be commited in order to update offset
		msg, err := k.Reader.FetchMessage(spanContext)
		if err != nil {
			logger.Error("Email Service Kafka", fmt.Errorf("could not read message: %v", err), logger.Params{})
			continue
		}
		logger.Info("Email Service Kafka", fmt.Sprintf("Received: %s", string(msg.Value)), logger.Params{})

		var emailSpanContext context.Context
		var emailSpan trace.Span
		if k.tracer != nil {
			emailSpanContext, emailSpan = (k.tracer).Start(spanContext, string(msg.Key))
			emailSpan.SetAttributes(attribute.Key(msg.Key).String(string(msg.Value)))
		} else {
			emailSpanContext = context.WithValue(spanContext, string(msg.Key), string(msg.Value))
		}

		switch string(msg.Key) {
		case template.WelcomeEmail:
			var emailTask email.WelcomeEmailBody
			if err = json.Unmarshal(msg.Value, &emailTask); err != nil {
				logger.Error("Email Service Kafka", fmt.Errorf("could not unmarshal message: %v", err), logger.Params{})
				continue
			}

			if err = emailTask.Process(emailSpanContext, k.Mailer); err != nil {
				logger.Error("Email Service Kafka", fmt.Errorf("could not process message: %v", err), logger.Params{})
				continue
			}
			if k.meter != nil {
				(*welcomeEmailCounterLock).Lock()
				welcomeEmailCounter.Add(emailSpanContext, 1)
				(*welcomeEmailCounterLock).Unlock()
			}

		case template.ResetEmail:
			var emailTask email.ResetEmailBody
			if err = json.Unmarshal(msg.Value, &emailTask); err != nil {
				logger.Error("Email Service Kafka", fmt.Errorf("could not unmarshal message: %v", err), logger.Params{})
				continue
			}

			if err = emailTask.Process(emailSpanContext, k.Mailer); err != nil {
				logger.Error("Email Service Kafka", fmt.Errorf("could not process message: %v", err), logger.Params{})
				continue
			}
		default:
			logger.Error("Email Service Kafka", errors.New("unmatched case"), logger.Params{})
			continue
		}
		if k.meter != nil {
			(*emailCounterLock).Lock()
			emailCounter.Add(emailSpanContext, 1)
			(*emailCounterLock).Unlock()
		}

		if err := k.Reader.CommitMessages(emailSpanContext, msg); err != nil {
			logger.Error("Email Service Kafka", fmt.Errorf("could not commit message: %v. Retrying", err), logger.Params{"message": string(msg.Value)})
			time.Sleep(2 * time.Second)
			if err := k.Reader.CommitMessages(emailSpanContext, msg); err != nil {
				logger.Error("Email Service Kafka", fmt.Errorf("could not commit message second time: %v", err), logger.Params{"message": string(msg.Value)})
				if k.tracer != nil {
					emailSpan.AddEvent("Could not commit message", trace.WithAttributes(attribute.Int("timestamp", int(time.Now().Unix()))))
				}
			}
		}
	}
}
