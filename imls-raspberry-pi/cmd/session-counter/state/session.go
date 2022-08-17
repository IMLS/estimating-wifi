package state

import (
	"time"
)

type sessionId struct {
	id int64
}

var (
	// singleton pattern. not thread safe.
	currentSession = sessionId{0}
)

func InitializeSession() sessionId {
	id := NewSessionID()
	currentSession := sessionId{id}
	return currentSession
}

func NewSessionID() int64 {
	return GetClock().Now().In(time.Local).Unix()
}

func GetCurrentSessionID() int64 {
	return currentSession.id
}

func IncrementSessionID() int64 {
	currentSession.id = NewSessionID()
	return currentSession.id
}
