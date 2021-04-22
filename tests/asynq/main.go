package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/xn3cr0nx/email-service/internal/email"
	"github.com/xn3cr0nx/email-service/internal/template"
	"github.com/xn3cr0nx/email-service/tests/utils"
)

const (
	address = "127.0.0.1:6379"
	db      = 2
)

func main() {
	r := asynq.RedisClientOpt{Addr: address, DB: db}
	c := asynq.NewClient(r)
	defer c.Close()

	for i := 0; i < 20; i++ {
		w := email.WelcomeEmailBody{
			From:    utils.RandomEmail(),
			To:      utils.RandomEmail(),
			Subject: utils.RandomString(),
			Params: email.WelcomeEmailBodyParams{
				Name: utils.RandomString(),
				URL:  utils.RandomEmail(),
			}}
		bytes, err := json.Marshal(w)
		if err != nil {
			return
		}
		var payload map[string]interface{}
		if err := json.Unmarshal(bytes, &payload); err != nil {
			return
		}
		t := asynq.NewTask(template.WelcomeEmail, payload)
		if err != nil {
			log.Fatalf("could not enqueue task: %v", err)
		}
		res, err := c.Enqueue(t)
		if err != nil {
			log.Fatalf("could not enqueue task: %v", err)
		}
		fmt.Printf("Enqueued Result: %+v\n", res)
	}
}
