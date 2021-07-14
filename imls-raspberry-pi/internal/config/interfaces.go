package config

type SessionId interface {
	GetSessionId() string
	IncrementSessionId() string
}
