// Package tlp provides sequential communicating processes.
package tlp

import (
	"gsa.gov/18f/cmd/session-counter/api"
	"gsa.gov/18f/cmd/session-counter/model"
	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/logwrapper"
)

// MacToEntry takes in a hashmap of MAC addresses and counts, and passes on a
// hashmap of manufacturer IDs and counts. Uses "unknown" for all unknown
// manufacturers.
func MacToEntry(ka *Keepalive, cfg *config.Config, macmap <-chan map[string]int, mfgmap chan<- map[string]model.Entry, chKill <-chan Ping) {
	lw := logwrapper.NewLogger(nil)
	lw.Debug("starting MacToEntry")

	// chKill will be nil in production
	var ping, pong chan interface{} = nil, nil
	if chKill == nil {
		ping, pong = ka.Subscribe("macToEntry", 5)
	}

	for {
		select {
		case <-ping:
			pong <- "macToEntry"
		case <-chKill:
			lw.Debug("exiting MacToEntry")
			return

		case mm := <-macmap:
			mfgs := make(map[string]model.Entry)
			for mac, count := range mm {
				mfg := api.MacToMfg(cfg, mac)
				mfgs[mac] = model.Entry{MAC: mac, Mfg: mfg, Count: count}
			}
			mfgmap <- mfgs
		}
	}
}
