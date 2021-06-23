package tlp

import (
	"log"

	"gsa.gov/18f/config"
	"gsa.gov/18f/session-counter/model"
)

/*
 * 1. Gather the MAC addresses that were seen.
 * 2. Map those to UIDs of mfg:counter
 * 3. Store what was seen into a DB with a timestamp.
 *    3a. If the UID exists, set the timestamp to now.
 *    3b. If the UID does not, insert it.
 * 4. If anything is older than the disconnection_window, remove it.
 * 5. WRONG Report this UID:timestamp pairing.
 */

func AlgorithmTwo(ka *Keepalive, cfg *config.Config, in <-chan []string, out chan<- map[string]int, reset <-chan Ping, ch_kill <-chan Ping) {
	if config.Verbose {
		log.Println("Starting AlgorithmTwo")
	}
	// This is our "tracking database"
	umdb := model.NewUMDB(cfg)

	// If we are running live, the kill channel is `nil`.
	// When we are live, THEN init the ping/pong.
	testing := true
	if ch_kill == nil {
		testing = false
	}
	var ping chan interface{} = nil
	var pong chan interface{} = nil
	if !testing {
		ping, pong = ka.Subscribe("AlgorithmTwo", 5)
		if config.Verbose {
			log.Println("a2: initialized keepalive")
		}
	}
	for {
		select {
		case <-ch_kill:
			if config.Verbose {
				log.Println("a2: exiting")
			}
			return
		case <-ping:
			pong <- "AlgorithmTwo"
		case <-reset:
			// Tell our mapping "db" to wipe itself.
			// This clears all counters, etc., and essentially
			// resets the algorithm as if we had just launched the whole process.
			umdb.WipeDB()

		case arr := <-in:

			// If we consider every message a "tick" of the clock, we need to advance time.
			umdb.AdvanceTime()

			// We get in a list of MAC addresses. Create mappings.
			// Timestamp everything as we see it, new or old.
			for _, mac := range arr {
				umdb.UpdateMapping(mac)
			}
			// Now, filter old things out
			umdb.RemoveOldMappings(cfg.Monitoring.UniquenessWindow)
			// Get the mappings as UserMappings, and send them out
			out <- umdb.AsUserMappings()
		}
	}

}
