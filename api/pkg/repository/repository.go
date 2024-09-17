package repository

import (
	"log"
	"os"

	"github.com/bovinxx/code-processor/api/pkg/repository/session"
	"github.com/bovinxx/code-processor/api/pkg/repository/task"
	"github.com/bovinxx/code-processor/api/pkg/repository/user"
)

type Repository struct {
	Task    task.Task
	User    user.User
	Session session.Session
}

func NewRepository() Repository {
	userDB, err := user.NewPostgres(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DB"))
	if err != nil {
		log.Fatalf("can't connect to postgres: %v", err)
	}
	taskDB, err := task.NewPostgres(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DB"))
	if err != nil {
		log.Fatalf("can't connect to postgres: %v", err)
	}
	sessionDB, err := session.NewRedis(os.Getenv("REDIS_PASSWORD"), os.Getenv("REDIS_HOST"))
	if err != nil {
		log.Fatalf("can't connect to redis: %v", err)
	}
	return Repository{
		Task:    taskDB,
		User:    userDB,
		Session: sessionDB,
	}
}
