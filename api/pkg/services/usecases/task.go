package usecases

import "github.com/bovinxx/code-processor/api/model"

type TaskUsecase interface {
	CreateTask(task model.Task) (string, error)
	GetStatus(id string) (string, error)
	GetResult(id string) (string, error)
}
