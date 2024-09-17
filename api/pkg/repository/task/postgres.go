package task

import (
	"database/sql"
	"fmt"

	"github.com/bovinxx/code-processor/api/model"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type postgres struct {
	db *sql.DB
}

func NewPostgres(user string, password string, host string, dbName string) (*postgres, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, dbName)
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &postgres{
		db: db,
	}, nil
}

func (p *postgres) Close() error {
	return p.db.Close()
}

func (p *postgres) CreateTask(task model.Task) (string, error) {
	id := uuid.New().String()
	task.Id = id
	_, err := p.db.Query(`INSERT INTO tasks (task_id) VALUES ($1)`, task.Id)
	if err != nil {
		return "", fmt.Errorf("failed create a new task: %v", err)
	}
	return id, nil
}

func (p *postgres) GetTask(id string) (model.Task, error) {
	var task model.Task
	err := p.db.QueryRow(`SELECT * FROM tasks WHERE task_id=$1`, id).Scan(&task.Id, &task.Result)
	if err != nil {
		if err == sql.ErrNoRows {
			return model.Task{}, fmt.Errorf("there's no such task")
		}
		return model.Task{}, fmt.Errorf("failed get task: %v", err)
	}
	return task, nil
}

func (p *postgres) DeleteTask(id string) error {
	_, err := p.db.Exec(`DELETE FROM tasks WHERE task_id=$1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}
	return nil
}

func (p *postgres) UpdateTask(id string, task model.Task) error {
	_, err := p.db.Exec(`UPDATE tasks SET result=$1 WHERE task_id=$2`, task.Result, id)
	if err != nil {
		return fmt.Errorf("failed to update task: %v", err)
	}
	return nil
}
