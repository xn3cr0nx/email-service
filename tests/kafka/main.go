package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/xn3cr0nx/email-service/internal/email"
	"github.com/xn3cr0nx/email-service/internal/template"
	"github.com/xn3cr0nx/email-service/tests/utils"
)

// the topic and broker address are initialized as constants
const (
	topic         = "emails"
	brokerAddress = "localhost:9092"
	// broker1Address = "localhost:19092"
	// broker2Address = "localhost:29092"
	// broker3Address = "localhost:39092"
	// specify the number of produced messages
	times = 10
)

func main() {
	// create a new context
	ctx := context.Background()
	// produce messages in a new go routine, since
	// both the produce and consume functions are
	// blocking
	produce(ctx)
}

func produce(ctx context.Context) {
	// intialize the writer with the broker addresses, and the topic
	w := kafka.NewWriter(kafka.WriterConfig{
		// Brokers: []string{broker1Address, broker2Address, broker3Address},
		Brokers: []string{brokerAddress},
		Topic:   topic,
	})

	for i := 0; i < times; i++ {
		message := email.WelcomeEmailBody{
			From:    utils.RandomEmail(),
			To:      utils.RandomEmail(),
			Subject: utils.RandomString(),
			Params: email.WelcomeEmailBodyParams{
				Name: utils.RandomString(),
				URL:  utils.RandomEmail(),
			}}

		bytes, err := json.Marshal(message)
		if err != nil {
			return
		}
		if err := w.WriteMessages(ctx, kafka.Message{
			Key:   []byte(template.WelcomeEmail),
			Value: bytes,
		}); err != nil {
			fmt.Print("could not write message " + err.Error())
			continue
		}

		fmt.Println("Email", message.To, "pushed")
	}
}
