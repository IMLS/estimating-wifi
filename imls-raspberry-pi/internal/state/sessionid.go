package state

import (
	"fmt"

	"gsa.gov/18f/config"
)

type SessionId struct {
	name    string
	cfg     *config.Config
	counter *Counter
}

func NewSessionId(cfg *config.Config) *SessionId {
	counter := NewCounter(cfg, "sessionid")
	counter.Reset()
	return &SessionId{name: "sessionid", cfg: cfg, counter: counter}
}

func (sid *SessionId) GetSessionId() string {
	return fmt.Sprint(sid.counter.Value())
}

func (sid *SessionId) IncrementSessionId() string {
	return fmt.Sprint(sid.counter.Increment())
}

func (sid *SessionId) PreviousSessionId() string {
	return fmt.Sprint(sid.counter.PrevValue())
}
