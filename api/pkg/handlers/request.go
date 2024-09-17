package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bovinxx/code-processor/api/model"
	"github.com/go-chi/chi"
)

type Request struct{}

func (req *Request) ParseAuthToken(r *http.Request) string {
	return r.Header.Get("Authorization")
}

func (req *Request) ParseTaskID(r *http.Request) string {
	return chi.URLParam(r, "task_id")
}

func (req *Request) ParseTask(r *http.Request) (*model.Task, error) {
	var task model.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (req *Request) ParseUser(r *http.Request) (*model.User, error) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
