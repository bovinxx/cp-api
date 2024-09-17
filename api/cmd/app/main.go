package main

import (
	"log"

	"github.com/bovinxx/code-processor/api/pkg/broker"
	"github.com/bovinxx/code-processor/api/pkg/handlers"
	"github.com/bovinxx/code-processor/api/pkg/repository"
	"github.com/bovinxx/code-processor/api/pkg/router"
	"github.com/bovinxx/code-processor/api/pkg/services"
	"github.com/bovinxx/code-processor/api/server"
)

// @title Code processor
// @version 1.0
// @description

// @host localhost:8000
// @BasePath /
func main() {
	db := repository.NewRepository()
	rabbitmq, err := broker.NewRabbitmq()
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq: %s", err.Error())
	}
	service := services.NewServices(&db, rabbitmq)
	handler := handlers.NewHandler(service)
	router := router.NewRouter(*handler)
	server := server.NewServer(router.Router)
	server.Run()
}
