package services

import (
	"fmt"

	"github.com/bovinxx/code-processor/api/model"
	"github.com/bovinxx/code-processor/api/pkg/broker"
	"github.com/bovinxx/code-processor/api/pkg/repository"
)

const (
	IN_PROGRESS = "in_progress"
	READY       = "ready"
)

type Task struct {
	repository *repository.Repository
	broker     broker.Broker
}

func NewTask(repository *repository.Repository, broker broker.Broker) *Task {
	return &Task{repository: repository, broker: broker}
}

func (s *Task) CreateTask(task model.Task) (string, error) {
	var id string
	var err error
	if id, err = (*s.repository).Task.CreateTask(task); err != nil {
		return "", err
	}
	task.Id = id
	err = s.broker.PublishMessage(task)
	if err != nil {
		task.Result = ""
		_ = (*s.repository).Task.UpdateTask(id, task)
		return id, fmt.Errorf("failed send task to a queue: %v", err)
	}
	return id, nil
}

func (s *Task) GetStatus(id string) (string, error) {
	task, err := (*s.repository).Task.GetTask(id)
	if err != nil {
		return "", err
	}
	if task.Result == "" {
		return IN_PROGRESS, nil
	}
	return READY, nil
}

func (s *Task) GetResult(id string) (string, error) {
	task, err := (*s.repository).Task.GetTask(id)
	if err != nil {
		return "", err
	}
	return task.Result, nil
}
