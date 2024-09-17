package model

import "time"

type Session struct {
	TimeAccessed time.Time `json:"TimeAccessed"`
}
