package session

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"sync"
	"time"

	"github.com/bovinxx/code-processor/api/model"
	"github.com/pkg/errors"
)

type ht struct {
	session_table map[string]model.Session
	lock          sync.Mutex
	maxlifetime   int64
}

func NewHt() *ht {
	ht := &ht{
		session_table: map[string]model.Session{},
		lock:          sync.Mutex{},
		maxlifetime:   1000,
	}
	go ht.SessionGC()
	return ht
}

func (ht *ht) SessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func (ht *ht) InitSession() (string, error) {
	ht.lock.Lock()
	defer ht.lock.Unlock()
	var session model.Session
	sid := ht.SessionID()
	if sid == "" {
		return "", errors.New("session error")
	}
	session.TimeAccessed = time.Now()
	ht.session_table[sid] = session
	return sid, nil
}

func (ht *ht) CheckSession(sid string) (string, error) {
	ht.lock.Lock()
	defer ht.lock.Unlock()
	if _, ok := ht.session_table[sid]; ok {
		return sid, nil
	}
	return "", errors.New("there's no such session")
}

func (ht *ht) DeleteSession(sid string) error {
	delete(ht.session_table, sid)
	return nil
}

func (ht *ht) SessionGC() {
	ht.lock.Lock()
	defer ht.lock.Unlock()

	for sid, session := range ht.session_table {
		if time.Now().Unix()-session.TimeAccessed.Unix() > ht.maxlifetime {
			ht.DeleteSession(sid)
		}
	}
	time.AfterFunc(time.Duration(ht.maxlifetime), func() { ht.SessionGC() })
}

func (ht *ht) UpdateSession(sid string) error {
	ht.lock.Lock()
	defer ht.lock.Unlock()
	if session, ok := ht.session_table[sid]; ok {
		session.TimeAccessed = time.Now()
		ht.session_table[sid] = session
	} else {
		return errors.New("error sid")
	}
	return nil
}
