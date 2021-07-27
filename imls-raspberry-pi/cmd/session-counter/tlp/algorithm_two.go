package tlp

import (
	"gsa.gov/18f/cmd/session-counter/model"
	"gsa.gov/18f/internal/logwrapper"
	"gsa.gov/18f/internal/state"
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

func AlgorithmTwo(ka *Keepalive, rb *ResetBroker, kb *KillBroker, in <-chan []string, out chan<- map[string]int) {
	cfg := state.GetConfig()
	lw := logwrapper.NewLogger(nil)
	lw.Debug("starting AlgorithmTwo")

	// This is our "tracking database"
	umdb := model.NewUMDB(cfg)

	// The reset broker manages comms for when we should
	// reset our internal structures
	chReset := rb.Subscribe()
	var ping, pong chan interface{} = nil, nil
	var chKill chan interface{} = nil
	if kb != nil {
		chKill = kb.Subscribe()
	} else {
		ping, pong = ka.Subscribe("AlgorithmTwo", 5)
	}

	for {
		select {
		case <-ping:
			pong <- "AlgorithmTwo"
		case <-chKill:
			lw.Debug("exiting AlgorithmTwo")
			return
		case <-chReset:
			// Tell our mapping "db" to wipe itself.
			// This clears all counters, etc., and essentially
			// resets the algorithm as if we had just launched the whole process.
			lw.Debug("wiping mfg/patron mapping DB")
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
			umdb.RemoveOldMappings(cfg.GetUniquenessWindow())
			// Get the mappings as UserMappings, and send them out
			um := umdb.AsUserMappings()
			lw.Debug("# user mappings [", len(um), "]")
			out <- um
		}
	}

}
