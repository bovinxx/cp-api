package task

import (
	"github.com/bovinxx/code-processor/api/model"
)

type Task interface {
	CreateTask(task model.Task) (string, error)
	GetTask(id string) (model.Task, error)
	DeleteTask(id string) error
	UpdateTask(id string, task model.Task) error
}
