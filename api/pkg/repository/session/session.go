package session

type Session interface {
	InitSession() (string, error)
	CheckSession(sid string) (string, error)
	DeleteSession(sid string) error
	SessionGC()
	UpdateSession(sid string) error
}
