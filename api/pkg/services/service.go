package services

import (
	"github.com/bovinxx/code-processor/api/pkg/broker"
	"github.com/bovinxx/code-processor/api/pkg/repository"
)

type Services struct {
	Auth *Auth
	Task *Task
}

func NewServices(repository *repository.Repository, broker broker.Broker) *Services {
	return &Services{
		Auth: NewAuth(repository),
		Task: NewTask(repository, broker),
	}
}
