package config

type SessionID interface {
	GetSessionId() string
	IncrementSessionId() string
	PreviousSessionId() string
}
