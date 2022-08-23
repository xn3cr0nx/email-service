package main

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/xn3cr0nx/email-service/internal/email"
	"github.com/xn3cr0nx/email-service/pkg/random"
)

type NatsMessage struct {
	Key   string
	Value []byte
}

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)

	msg := NatsMessage{
		Key:   "email:welcome",
		Value: []byte{},
	}

	for i := 0; i < 20; i++ {
		log.Println("Baking message")
		w := email.WelcomeEmailBody{
			From:    random.Email(),
			To:      random.Email(),
			Subject: random.String(),
			Params: email.WelcomeEmailBodyParams{
				Name: random.String(),
				URL:  random.Email(),
			}}
		value, err := json.Marshal(w)
		if err != nil {
			return
		}
		msg.Value = value
		bytes, err := json.Marshal(msg)
		if err != nil {
			return
		}

		log.Println("Publishing message")
		if err = nc.Publish("emails", bytes); err != nil {
			log.Fatalf("could not enqueue task: %v", err)
		}
	}
}
