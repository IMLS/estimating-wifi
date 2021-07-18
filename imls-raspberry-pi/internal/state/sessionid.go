package state

import (
	"log"
	"sync"
)

var sid_once sync.Once

func (cfg *CFG) InitializeSessionId() {
	sid_once.Do(func() {
		sid := 0
		tdb := NewSqliteDB(the_config.Databases.Durations)
		if tdb.CheckTableExists("durations") {
			tdb.Open()
			defer tdb.Close()
			var sessionId int
			err := tdb.Ptr.Get(&sessionId, "SELECT MAX(session_id) FROM durations")
			if err != nil {
				log.Println("error in finding max session id; returning 0")
				log.Println(err.Error())
			}
			sid = sessionId
		} else {
			log.Println("durations table did not exist; returning session id 0")
			sid = 0
		}
		the_config.SessionId = sid
	})
}

func (cfg *CFG) GetCurrentSessionID() int {
	return the_config.SessionID
}

func (cfg *CFG) IncrementSessionID() int {
	the_config.SessionId = the_config.SessionId + 1
	return the_config.SessionID
}

func (cfg *CFG) GetPreviousSessionID() int {
	return the_config.SessionID - 1
}
