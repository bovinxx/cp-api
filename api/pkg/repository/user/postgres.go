package user

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

func (p *postgres) CreateUser(user model.User) (string, error) {
	id := uuid.New().String()
	user.Id = id
	_, err := p.db.Query(`INSERT INTO users (user_id, username, password_hash) VALUES ($1, $2, $3)`, user.Id, user.Username, user.Password)
	if err != nil {
		return "", fmt.Errorf("failed create new user: %v", err)
	}
	return id, nil
}

func (p *postgres) CheckUser(login, password string) (string, error) {
	var userId string
	err := p.db.QueryRow(`SELECT user_id FROM users WHERE username=$1 AND password_hash=$2`, login, password).Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found")
		}
		return "", fmt.Errorf("failed to check user: %w", err)
	}
	return userId, nil
}

func (p *postgres) CheckLogin(login string) (bool, error) {
	var userId string
	err := p.db.QueryRow(`SELECT user_id FROM users WHERE username=$1`, login).Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check login: %w", err)
	}
	return true, nil
}
