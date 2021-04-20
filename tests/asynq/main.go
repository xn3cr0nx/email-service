package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/xn3cr0nx/email-service/internal/task"
	"github.com/xn3cr0nx/email-service/internal/template"
)

func main() {
	r := asynq.RedisClientOpt{Addr: "127.0.0.1:6379", DB: 2}
	c := asynq.NewClient(r)
	defer c.Close()

	for i := 0; i < 20; i++ {
		w := task.WelcomeEmailTask{
			From:    RandomString(),
			To:      RandomString(),
			Subject: RandomString(),
			Params: task.WelcomeEmailTaskParams{
				Name: RandomString(),
				URL:  RandomEmail(),
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

		// t, err := NewEmailDeliveryTask("Me", "You", "Test", "You", "http://test.com")
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

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const digits = "0123456789"

// RandomString returns a random generated string with fixed length
func RandomString() string {
	return RandomStringWithLen(20)
}

// RandomStringWithLen returns a random string with length
func RandomStringWithLen(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// RandomEmail returns a randomly generated email address
func RandomEmail() string {
	localPart := fmt.Sprintf("%s.%s", RandomStringWithLen(5), RandomStringWithLen(7))
	domain := RandomDomain()
	return strings.ToLower(fmt.Sprintf("%s@%s", localPart, domain))
}

// RandomDomain returns a randomly generated domain
func RandomDomain() string {
	domain := fmt.Sprintf("%s.%s", RandomStringWithLen(10), RandomStringWithLen(3))
	return strings.ToLower(domain)
}
