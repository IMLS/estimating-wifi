package tlp

import (
	"github.com/rs/zerolog/log"
	"gsa.gov/18f/internal/session-counter-helper/state"
)

func ProcessData(dDB *state.DurationsDB, sq *state.Queue[int64]) bool {
	// Queue up what needs to be sent still.
	session := state.GetCurrentSessionID()

	log.Debug().
		Int64("session", session).
		Msg("queueing to send")

	if session >= 0 {
		sq.Enqueue(session)
	}

	macs := state.GetMACs()
	// FIGURE OUT FORMAT OF GetMACs return value
	for mac, element := range macs {
		log.Debug().
			Int64("start", element.Start).
			Int64("end", element.End).
			Msg(mac)
	}
	dDB.InsertMany(session, macs)
	return true
}
