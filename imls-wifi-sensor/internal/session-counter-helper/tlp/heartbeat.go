package tlp

import (
	"github.com/rs/zerolog/log"
	"gsa.gov/18f/internal/session-counter-helper/api"
)

func HeartBeat() {
	log.Debug().Msg("Running heartbeat")
	err := api.PostHeartBeat()
	if err != nil {
		log.Error().
			Err(err).
			Msg("could not provide heartbeat")
	}
}
