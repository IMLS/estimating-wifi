package state

import (
	"fmt"
	"path/filepath"

	"gsa.gov/18f/internal/config"
)

type SessionId struct {
	name    string
	cfg     *config.Config
	counter *Counter
}

func NewSessionId(cfg *config.Config) *SessionId {
	fullpath := filepath.Join(cfg.Local.WebDirectory, DURATIONSDB)
	tdb := NewSqliteDB(DURATIONSDB, fullpath)
	// If the table doesn't exist, create a counter and return it.
	//log.Println("table does not exist: ", tdb.CheckTableDoesNotExist("sessionid"))
	if tdb.CheckTableDoesNotExist("sessionid") {
		counter := NewCounter(cfg, "sessionid")
		counter.Reset()
		return &SessionId{name: "sessionid", cfg: cfg, counter: counter}
	} else {
		// log.Println("getting counter...")
		ctr := GetCounter(cfg, "sessionid")
		sid := &SessionId{name: "sessionid", cfg: cfg, counter: ctr}
		// log.Println("old session id value", sid.GetSessionId())
		sid.IncrementSessionId()
		//log.Println("new session id value", sid.GetSessionId())
		return sid
	}
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
