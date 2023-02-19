package tlp

import (
	"github.com/rs/zerolog/log"
	"gsa.gov/18f/internal/session-counter-helper/api"
	"gsa.gov/18f/internal/wifi-hardware-search/netadapter"
)

func HeartBeat() {
	log.Debug().Msg("Running heartbeat")

	// through exhaustive testing we've determined the adapater is not stable
	// in Windows over long periods. As a prophylactic measure, bounce the
	// adapter at every heartbeat (typically hourly)
	netadapter.RestartNetAdapter()

	err := api.PostHeartBeat()
	if err != nil {
		log.Error().
			Err(err).
			Msg("could not provide heartbeat")
	}
}
