package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/bovinxx/code-processor/api/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	QUEUE_NAME = "queue_for_tasks"
)

type Rabbitmq struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitmq() (Rabbitmq, error) {
	uri := fmt.Sprintf("amqp://%s:%s@%s/", os.Getenv("RMQ_USER"), os.Getenv("RMQ_PASSWORD"), os.Getenv("RMQ_HOST"))
	conn, err := amqp.Dial(uri)
	if err != nil {
		return Rabbitmq{}, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %v", err)
	}
	_, err = ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to declare a queue: %v", err)
	}

	return Rabbitmq{conn: conn, ch: ch}, nil
}

func (r Rabbitmq) PublishMessage(task model.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	taskJson, err := json.Marshal(task)
	if err != nil {
		log.Fatalf("failed to parse task json: %v", err)
		return err
	}

	corrId := randomString(32)
	err = r.ch.PublishWithContext(ctx,
		"",
		"rpc_queue",
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       QUEUE_NAME,
			Body:          taskJson,
		})
	if err != nil {
		log.Fatalf("failed to publish a message: %v", err)
		return err
	}
	return nil
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
