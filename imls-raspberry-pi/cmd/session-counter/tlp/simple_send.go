package tlp

import (
	"strconv"

	"github.com/rs/zerolog/log"
	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/cmd/session-counter/api"
	"gsa.gov/18f/cmd/session-counter/state"
	"gsa.gov/18f/cmd/session-counter/structs"
)

func SimpleSend(db *state.DurationsDB) {
	log.Debug().
		Msg("starting batch send")

	// This only comes in on reset...
	sq := state.NewQueue[int64]("sent")
	sessionsToSend := sq.AsList()

	for _, nextSessionIDToSend := range sessionsToSend {
		durations := []*structs.Duration{}
		// FIXME: Leaky Abstraction
		// err := db.GetPtr().Select(&durations, "SELECT * FROM durations WHERE session_id=?", nextSessionIDToSend)
		durations = db.GetSession(nextSessionIDToSend)

		if len(durations) == 0 {
			log.Debug().
				Str("session", strconv.FormatInt(nextSessionIDToSend, 10)).
				Msg("found zero durations")
			sq.Remove(nextSessionIDToSend)
		} else {
			log.Debug().
				Int("durations", len(durations)).
				Str("session", strconv.FormatInt(nextSessionIDToSend, 10)).
				Msg("preparing to send durations to API")

			// convert []Duration to an array of map[string]interface{}
			data := make([]map[string]interface{}, 0)
			for _, duration := range durations {
				data = append(data, duration.AsMap())
			}

			// After writing images, we come back and try and send the data remotely.
			log.Debug().
				Int("duration", len(data)).
				Str("session", strconv.FormatInt(nextSessionIDToSend, 10)).
				Msg("sending durations to API")

			err := api.PostJSON(config.GetDurationsURI(), data)
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
