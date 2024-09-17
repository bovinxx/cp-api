package user

import "github.com/bovinxx/code-processor/api/model"

type User interface {
	CreateUser(user model.User) (string, error)
	CheckUser(login, password string) (string, error)
	CheckLogin(login string) (bool, error)
}
