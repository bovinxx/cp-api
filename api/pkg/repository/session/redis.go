package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/bovinxx/code-processor/api/model"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

const (
	SESSION_EXPIRATION = 3 * time.Minute
)

var (
	ctx = context.TODO()
)

type session model.Session

func (s session) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func NewRedis(password string, host string) (*Redis, error) {
	url := fmt.Sprintf("redis://%s@%s", password, host)
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	if err = client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &Redis{
		client: client,
	}, nil
}

func (r *Redis) SessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (r *Redis) InitSession() (string, error) {
	sid := r.SessionID()
	session := session{}
	session.TimeAccessed = time.Now()
	if err := r.client.Set(ctx, sid, session, SESSION_EXPIRATION).Err(); err != nil {
		return "", err
	}
	return sid, nil
}

func (r *Redis) CheckSession(sid string) (string, error) {
	cmd := r.client.Get(ctx, sid)
	err := cmd.Err()
	if err != nil {
		return "", errors.New("there's no such session")
	}
	return sid, nil
}

func (r *Redis) DeleteSession(sid string) error {
	return r.client.Del(ctx, sid).Err()
}

func (r *Redis) SessionGC() {

}

func (r *Redis) UpdateSession(sid string) error {
	var session session
	session.TimeAccessed = time.Now()
	if err := r.client.Set(ctx, sid, session, SESSION_EXPIRATION).Err(); err != nil {
		return errors.New("error sid")
	}
	return nil
}
