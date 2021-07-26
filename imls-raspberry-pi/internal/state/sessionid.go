package state

import (
	"log"
	"sync"
)

var sid_once sync.Once

func (cfg *CFG) InitializeSessionID() {
	sid_once.Do(func() {
		sid := 0
		tdb := theConfig.Databases.DurationsDB
		if tdb != nil && tdb.CheckTableExists("durations") {
			var sessionId int
			err := tdb.GetPtr().Get(&sessionId, "SELECT MAX(session_id) FROM durations")
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
