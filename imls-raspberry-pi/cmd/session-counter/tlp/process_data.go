package tlp

import (
	"github.com/rs/zerolog/log"
	"gsa.gov/18f/cmd/session-counter/state"
)

func ProcessData(dDB *state.DurationsDB, sq *state.Queue[int64]) bool {
	// Queue up what needs to be sent still.
	session := state.GetCurrentSessionID()

	log.Debug().
		Int64("session", session).
		Msg("queueing to images and send")

	if session >= 0 {
		sq.Enqueue(session)
	}

	macs := state.GetMACs()
	dDB.InsertMany(session, macs)
	return true
}
