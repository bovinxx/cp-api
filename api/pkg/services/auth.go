package services

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"strings"

	"github.com/bovinxx/code-processor/api/model"
	"github.com/bovinxx/code-processor/api/pkg/repository"
)

const (
	SALT = "31#kop3wkeop32k3##(!(3p12o321o-32eo32-=o1="
)

type Auth struct {
	repository *repository.Repository
}

func NewAuth(repository *repository.Repository) *Auth {
	return &Auth{
		repository: repository,
	}
}

func (s *Auth) CreateUser(login, password string) (string, error) {
	if !isValidLogin(login) || !isValidPassword(password) {
		return "", errors.New("invalid login or password (len >= 5 and no spaces)")
	}
	var id string
	var err error
	ok, _ := (*s.repository).User.CheckLogin(login)
	if ok {
		return "", errors.New("login already exists")
	}

	user := model.User{Username: login, Password: hashPassword(password)}
	id, err = (*s.repository).User.CreateUser(user)
	if err != nil {
		return "", err
	}
	return id, nil
}

func isValidLogin(login string) bool {
	if len(login) >= 5 && !strings.Contains(login, " ") {
		return true
	}
	return false
}

func isValidPassword(password string) bool {
	if len(password) >= 5 && !strings.Contains(password, " ") {
		return true
	}
	return false
}

func hashPassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(SALT)))
}

func (s *Auth) LoginUser(login, password string) (string, error) {
	if !isValidLogin(login) || !isValidPassword(password) {
		return "", errors.New("invalid login or password (len >= 5 and no spaces)")
	}
	id, err := (*s.repository).User.CheckUser(login, hashPassword(password))
	if err != nil {
		return "", errors.New("incorrect login and password")
	}
	if id == "" {
		return "", errors.New("incorrect password")
	}
	sid, err := (*s.repository).Session.InitSession()
	if err != nil {
		return "", err
	}
	return sid, nil
}

func (s *Auth) Authorization(token string) (bool, error) {
	_, err := (*s.repository).Session.CheckSession(token)
	if err != nil {
		return false, err
	}
	return true, nil
}
