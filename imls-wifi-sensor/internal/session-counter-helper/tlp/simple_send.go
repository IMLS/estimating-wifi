package tlp

import (
	"fmt"
	"strconv"

	"github.com/rs/zerolog/log"
	"gsa.gov/18f/internal/session-counter-helper/api"
	"gsa.gov/18f/internal/session-counter-helper/state"
)

func SimpleSend(db *state.DurationsDB, sq *state.Queue[int64]) {
	log.Debug().
		Msg("starting batch send")

	// This only comes in on reset...
	//sq := state.NewQueue[int64]("sent")
	sessionsToSend := sq.AsList()
	log.Debug().
		Str("sessionsToSend", fmt.Sprint(sessionsToSend)).
		Msg("sessions in queue to be sent")

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
