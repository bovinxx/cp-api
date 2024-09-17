package task

import (
	"errors"

	"github.com/bovinxx/code-processor/api/model"
	"github.com/google/uuid"
)

type ht struct {
	task_table map[string]model.Task
}

func NewHt() *ht {
	return &ht{
		task_table: map[string]model.Task{},
	}
}

func (ht *ht) CreateTask(task model.Task) (string, error) {
	id := uuid.New().String()
	task.Id = id
	ht.task_table[id] = task
	return id, nil
}

func (ht *ht) GetTask(id string) (model.Task, error) {
	t, ok := ht.task_table[id]
	if !ok {
		return model.Task{}, errors.New("there's no such task")
	}
	return t, nil
}

func (ht *ht) DeleteTask(id string) error {
	delete(ht.task_table, id)
	return nil
}

func (ht *ht) UpdateTask(id string, task model.Task) error {
	ht.task_table[id] = task
	return nil
}
