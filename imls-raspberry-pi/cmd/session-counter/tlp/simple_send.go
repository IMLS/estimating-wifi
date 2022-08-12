package tlp

import (
	"strconv"

	"github.com/rs/zerolog/log"
	"gsa.gov/18f/cmd/session-counter/api"
	"gsa.gov/18f/cmd/session-counter/state"
)

func SimpleSend(db *state.DurationsDB) {
	log.Debug().
		Msg("starting batch send")

	// This only comes in on reset...
	sq := state.NewQueue[int64]("sent")
	sessionsToSend := sq.AsList()

	for _, nextSessionIDToSend := range sessionsToSend {

		durations := db.GetSession(nextSessionIDToSend)

		if len(durations) == 0 {
			log.Debug().
				Str("session", strconv.FormatInt(nextSessionIDToSend, 10)).
				Msg("found zero durations")
			sq.Remove(nextSessionIDToSend)
		} else {

			log.Debug().
				Int("duration", len(durations)).
				Str("session", strconv.FormatInt(nextSessionIDToSend, 10)).
				Msg("sending durations to API")

			err := api.PostDurations(durations)

			if err != nil {
				log.Error().
					Str("session", strconv.FormatInt(nextSessionIDToSend, 10)).
					Err(err).
					Msg("could not send; data left on queue")
			} else {
				// If we successfully sent the data remotely, we can now mark it is as sent.
				sq.Remove(nextSessionIDToSend)
			}
		}
	}
}
