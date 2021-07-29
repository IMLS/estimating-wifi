package tlp

import (
	"log"

	"gsa.gov/18f/internal/interfaces"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"
)

func ProcessData(db interfaces.Database, sq *state.Queue, iq *state.Queue) bool {
	var eds []structs.EphemeralDuration
	db.GetPtr().Select(&eds, "SELECT * FROM ephemeraldurations")
	log.Println(eds)
	return true
}
