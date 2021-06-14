package tlp

import (
	"log"

	"gsa.gov/18f/config"
	"gsa.gov/18f/session-counter/api"
	"gsa.gov/18f/session-counter/model"
)

/* PROC mac_to_mfg
 * Takes in a hashmap of MAC addresses and counts, and passes on a hashmap
 * of manufacturer IDs and counts.
 * Uses "unknown" for all unknown manufacturers.
 */
func MacToEntry(ka *Keepalive, cfg *config.Config, macmap <-chan map[string]int, mfgmap chan<- map[string]model.Entry, ch_kill <-chan Ping) {
	log.Println("Starting macToEntry")

	// ch_kill will be nil in production
	var ping, pong chan interface{} = nil, nil
	if ch_kill == nil {
		ping, pong = ka.Subscribe("macToEntry", 5)
	}

	for {
		select {
		case <-ping:
			pong <- "macToEntry"
		case <-ch_kill:
			log.Println("Exiting MacToEntry")
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
