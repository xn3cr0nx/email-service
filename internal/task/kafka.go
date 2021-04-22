package task

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/segmentio/kafka-go"
	"github.com/xn3cr0nx/email-service/internal/email"
	"github.com/xn3cr0nx/email-service/internal/mailer"
	"github.com/xn3cr0nx/email-service/internal/template"
	"github.com/xn3cr0nx/email-service/pkg/logger"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type KafkaEmailConsumer struct {
	Reader *kafka.Reader
	Mailer mailer.Service
	tracer *trace.Tracer
	meter  *metric.Meter
}

func NewKafkaEmailConsumer(k *kafka.Reader, m mailer.Service, tracer *trace.Tracer, meter *metric.Meter) *KafkaEmailConsumer {
	return &KafkaEmailConsumer{k, m, tracer, meter}
}

func (k *KafkaEmailConsumer) Run(ctx context.Context) {
	emailCounterLock := new(sync.RWMutex)
	var emailCounter metric.Int64Counter
	welcomeEmailCounterLock := new(sync.RWMutex)
	var welcomeEmailCounter metric.Int64Counter
	if k.meter != nil {
		emailCounter = metric.Must(*k.meter).NewInt64Counter("kafka.emails")
		welcomeEmailCounter = metric.Must(*k.meter).NewInt64Counter("kafka.emails.welcome")
	}

	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := k.Reader.ReadMessage(ctx)
		if err != nil {
			logger.Error("Email Service Kafka", fmt.Errorf("Could not read message: %v", err), logger.Params{})
			continue
		}
		logger.Info("Email Service Kafka", fmt.Sprintf("Received: %s", string(msg.Value)), logger.Params{})

		// var span trace.Span
		// if k.tracer != nil {
		// 	_, span = (*k.tracer).Start(ctx, "email")
		// 	span.SetAttributes(attribute.Key(msg.Key).String(string(msg.Value)))
		// 	defer span.End()
		// }

		switch string(msg.Key) {
		case template.WelcomeEmail:
			var emailTask email.WelcomeEmailBody
			if err = json.Unmarshal(msg.Value, &emailTask); err != nil {
				logger.Error("Email Service Kafka", fmt.Errorf("Could not unmarshal message: %v", err), logger.Params{})
				continue
			}

			if err = emailTask.Process(ctx, k.Mailer); err != nil {
				logger.Error("Email Service Kafka", fmt.Errorf("Could not process message: %v", err), logger.Params{})
				continue
			}
			if k.meter != nil {
				(*welcomeEmailCounterLock).Lock()
				welcomeEmailCounter.Add(ctx, 1)
				(*welcomeEmailCounterLock).Unlock()
			}
		}
		if k.meter != nil {
			(*emailCounterLock).Lock()
			emailCounter.Add(ctx, 1)
			(*emailCounterLock).Unlock()
		}
	}
}
