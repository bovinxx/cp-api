package main

import (
	"log"
	"net/http"

	"github.com/bovinxx/code-processor/processor/broker"
	docker "github.com/bovinxx/code-processor/processor/docker_client"
	"github.com/bovinxx/code-processor/processor/repository"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	db := initPostgres()
	client := initDocker()
	rabbitmq := initRabbitMQ(db)

	go func() {
		if err := rabbitmq.FetchMessages(client); err != nil {
			log.Fatalf("failed to fetch messages: %v", err)
		}
	}()

	startMetrics()
	var forever chan struct{}
	<-forever
}

func initPostgres() repository.Repository {
	db, err := repository.NewPostgres()
	checkErr(err, "failed to connect to postgres")
	log.Println("Connected to Postgres successfully")
	return db
}

func initDocker() *docker.Client {
	client, err := docker.NewClient("basic_image")
	checkErr(err, "failed to create a docker client")
	log.Println("Docker client created successfully")
	return client
}

func initRabbitMQ(db repository.Repository) *broker.Rabbitmq {
	rabbitmq, err := broker.NewRabbitmq(db)
	checkErr(err, "failed to connect to RabbitMQ")
	log.Println("Connected to RabbitMQ successfully")
	return rabbitmq
}

func startMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":3831", nil); err != nil {
			log.Fatalf("failed to start the HTTP server with metrics: %v", err)
		}
	}()
	log.Println("The built-in HTTP server with metrics is running")
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}
