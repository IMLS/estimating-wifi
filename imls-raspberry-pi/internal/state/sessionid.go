package state

import (
	"fmt"

	"gsa.gov/18f/internal/config"
)

type SessionId struct {
	name    string
	cfg     *config.Config
	counter *Counter
}

func NewSessionId(cfg *config.Config) *SessionId {
	ctr := GetCounter(cfg, "sessionid")
	if ctr.db.Ptr != nil {
		if ctr.db.CheckTableExists("sessionid") {
			return &SessionId{name: "sessionid", cfg: cfg, counter: ctr}
		}
	}
	// implicit else...
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
