package tlp

import (
	"github.com/rs/zerolog/log"
	"gsa.gov/18f/cmd/session-counter/state"
	"gsa.gov/18f/internal/http"
	"gsa.gov/18f/internal/interfaces"
	"gsa.gov/18f/internal/structs"
)

func SimpleSend(db interfaces.Database) {
	log.Debug().
		Msg("starting batch send")

	// This only comes in on reset...
	sq := state.NewQueue[int64]("sent")
	sessionsToSend := sq.AsList()

	for _, nextSessionIDToSend := range sessionsToSend {
		durations := []structs.Duration{}
		// FIXME: Leaky Abstraction
		err := db.GetPtr().Select(&durations, "SELECT * FROM durations WHERE session_id=?", nextSessionIDToSend)

		if err != nil {
			log.Error().
				Err(err).
				Str("session", nextSessionIDToSend).
				Msg("could not extract durations")
		}

		if len(durations) == 0 {
			log.Debug().
				Str("session", nextSessionIDToSend).
				Msg("found zero durations")
			sq.Remove(nextSessionIDToSend)
		} else if state.IsStoringToAPI() {
			log.Debug().
				Int("durations", len(durations)).
				Str("session", nextSessionIDToSend).
				Msg("preparing to send durations to API")

			// convert []Duration to an array of map[string]interface{}
			data := make([]map[string]interface{}, 0)
			for _, duration := range durations {
				data = append(data, duration.AsMap())
			}

			// After writing images, we come back and try and send the data remotely.
			log.Debug().
				Int("duration", len(data)).
				Str("session", nextSessionIDToSend).
				Msg("sending durations to API")

			err = http.PostJSON(state.GetDurationsURI(), data)
			if err != nil {
				log.Error().
					Str("session", nextSessionIDToSend).
					Err(err).
					Msg("could not send; data left on queue")
			} else {
				// If we successfully sent the data remotely, we can now mark it is as sent.
				sq.Remove(nextSessionIDToSend)
			}
		} else {
			// Always dequeue. We're storing locally "for free" into the
			// durations table before trying to do the send.
			log.Info().
				Msg("not in API mode, not sending data")
			sq.Remove(nextSessionIDToSend)
		}
	}

}
