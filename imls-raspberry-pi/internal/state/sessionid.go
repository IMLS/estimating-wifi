package state

import (
	"log"
	"sync"
)

var sid_once sync.Once

func (cfg *CFG) InitializeSessionId() {
	sid_once.Do(func() {
		sid := 0
		tdb := NewSqliteDB(theConfig.Databases.DurationsPath)
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
		theConfig.SessionId = sid + 1
	})
}

func (cfg *CFG) GetCurrentSessionID() int {
	return theConfig.SessionId
}

func (cfg *CFG) IncrementSessionID() int {
	theConfig.SessionId = theConfig.SessionId + 1
	return theConfig.SessionId
}

func (cfg *CFG) GetPreviousSessionID() int {
	return theConfig.SessionId - 1
}
