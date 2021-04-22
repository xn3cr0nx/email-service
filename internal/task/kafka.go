package task

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/xn3cr0nx/email-service/internal/email"
	"github.com/xn3cr0nx/email-service/internal/mailer"
	"github.com/xn3cr0nx/email-service/internal/template"
	"github.com/xn3cr0nx/email-service/pkg/logger"
)

type KafkaEmailConsumer struct {
	Reader *kafka.Reader
	Mailer mailer.Service
}

func NewKafkaEmailConsumer(k *kafka.Reader, m mailer.Service) *KafkaEmailConsumer {
	return &KafkaEmailConsumer{k, m}
}

func (k *KafkaEmailConsumer) Run(ctx context.Context) {
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := k.Reader.ReadMessage(ctx)
		if err != nil {
			logger.Error("Email Service Kafka", fmt.Errorf("Could not read message: %v", err), logger.Params{})
			continue
		}
		fmt.Println("received: ", string(msg.Value))
		logger.Info("Email Service Kafka", fmt.Sprintf("Received: %s", string(msg.Value)), logger.Params{})

		switch string(msg.Key) {
		case template.WelcomeEmail:
			var emailTask email.WelcomeEmailBody
			if err = json.Unmarshal(msg.Value, &emailTask); err != nil {
				logger.Error("Email Service Kafka", fmt.Errorf("Could not unmarshal message: %v", err), logger.Params{})
			}

			if err = emailTask.Process(k.Mailer); err != nil {
				logger.Error("Email Service Kafka", fmt.Errorf("Could not process message: %v", err), logger.Params{})
			}
		}
	}
}
