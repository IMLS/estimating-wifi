package state

import (
	"log"

	"gsa.gov/18f/internal/interfaces"
)

func InitializeSessionID(durationsDB interfaces.Database) int {
	sid := 0
	if durationsDB.CheckTableExists("durations") {
		var sessionID int
		err := durationsDB.GetPtr().Get(&sessionID, "SELECT MAX(session_id) FROM durations")
		if err != nil {
			log.Println("error in finding max session id; returning 0")
			log.Println(err.Error())
		}
		sid = sessionID
	} else {
		log.Println("durations table did not exist; returning session id 0")
		sid = 0
	}
	return sid + 1
}

func (dc *databaseConfig) GetCurrentSessionID() int {
	return dc.sessionID
}

func (dc *databaseConfig) IncrementSessionID() int {
	dc.sessionID++
	return dc.sessionID
}

func (dc *databaseConfig) GetPreviousSessionID() int {
	return dc.sessionID - 1
}
