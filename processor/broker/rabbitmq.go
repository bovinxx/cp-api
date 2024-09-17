package broker

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bovinxx/code-processor/api/model"
	docker "github.com/bovinxx/code-processor/processor/docker_client"
	"github.com/bovinxx/code-processor/processor/metrics"
	"github.com/bovinxx/code-processor/processor/repository"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Rabbitmq struct {
	conn       *amqp.Connection
	ch         *amqp.Channel
	queue      amqp.Queue
	repository repository.Repository
}

var (
	metric = metrics.NewMetrics()
)

func NewRabbitmq(repository repository.Repository) (*Rabbitmq, error) {
	uri := fmt.Sprintf("amqp://%s:%s@%s/", os.Getenv("RMQ_USER"), os.Getenv("RMQ_PASSWORD"), os.Getenv("RMQ_HOST"))
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	q, err := ch.QueueDeclare(
		"rpc_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	err = ch.Qos(
		1,
		0,
		false,
	)
	if err != nil {
		return nil, err
	}

	return &Rabbitmq{repository: repository, ch: ch, queue: q, conn: conn}, nil
}

func (r *Rabbitmq) Close() error {
	if r.ch != nil {
		if err := r.ch.Close(); err != nil {
			return err
		}
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

func (r *Rabbitmq) FetchMessages(client *docker.Client) error {
	if r.ch == nil || r.conn == nil {
		return fmt.Errorf("channel or connection is closed")
	}
	msgs, err := r.ch.Consume(
		r.queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("fuck it all")
		return err
	}

	go func() {
		log.Println("we are here")
		for d := range msgs {
			if err := r.ProcessMessage(client, &d); err != nil {
				log.Fatalf("failed to process msg: %v", err)
			}

		}
	}()

	return nil
}

func (r *Rabbitmq) ProcessMessage(client *docker.Client, msg *amqp.Delivery) error {
	start := time.Now()
	var data model.Task
	data.Translator = ""
	defer metric.Record(start, data.Translator)
	if err := json.Unmarshal(msg.Body, &data); err != nil {
		return err
	}
	metric.TranslatorsUsed.WithLabelValues(data.Translator).Inc()
	res, err := client.ExecCode(data.Translator, []byte(data.Code))
	if err != nil {
		return err
	}
	if err := r.repository.PutResult(data.Id, string(res)); err != nil {
		log.Fatalf("failed to pust result: %v", err)
	}
	msg.Ack(false)
	return nil
}
