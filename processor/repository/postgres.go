package repository

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type postgres struct {
	db *sql.DB
}

func NewPostgres() (postgres, error) {
	url := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DB"))
	db, err := sql.Open("postgres", url)
	if err != nil {
		return postgres{}, err
	}
	if err = db.Ping(); err != nil {
		return postgres{}, err
	}
	return postgres{db: db}, nil
}

func (p postgres) PutResult(taskId string, result string) error {
	_, err := p.db.Exec(`UPDATE tasks SET result=$1 WHERE task_id=$2`, result, taskId)
	if err != nil {
		return fmt.Errorf("failed to update task: %v", err)
	}
	return nil
}
