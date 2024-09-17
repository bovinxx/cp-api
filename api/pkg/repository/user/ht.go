package user

import (
	"errors"

	"github.com/bovinxx/code-processor/api/model"
	"github.com/google/uuid"
)

type ht struct {
	user_table map[string]model.User
}

func NewHt() *ht {
	return &ht{
		user_table: map[string]model.User{},
	}
}

func (ht *ht) CreateUser(user model.User) (string, error) {
	id := uuid.New().String()
	user.Id = id
	ht.user_table[id] = user
	return id, nil
}

func (ht *ht) CheckUser(login, password string) (string, error) {
	for _, user := range ht.user_table {
		if user.Username == login {
			if user.Password == password {
				return user.Id, nil
			}
			return "", errors.New("wrong password")
		}
	}
	return "", errors.New("such user doesn't exist")
}
func (ht *ht) CheckLogin(login string) (bool, error) {
	for _, user := range ht.user_table {
		if user.Username == login {
			return true, nil
		}
	}
	return false, nil
}
