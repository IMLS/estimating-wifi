package tlp

import (
	"log"
	"time"

	"gsa.gov/18f/session-counter/api"
	"gsa.gov/18f/session-counter/config"
	"gsa.gov/18f/session-counter/csp"
	"gsa.gov/18f/session-counter/model"
)

type uniqueMapping struct {
	lastid    int
	uid       map[string]int
	mfg       map[string]string
	timestamp map[string]time.Time
}

func (um uniqueMapping) init() {
	um.lastid = 0
	um.uid = make(map[string]int)
	um.mfg = make(map[string]string)
	um.timestamp = make(map[string]time.Time)
}

//
func (um uniqueMapping) updateMapping(cfg *config.Config, mac string) {
	_, ok := um.mfg[mac]
	// If we didn't find it, then we need to add it.
	if !ok {
		// Assign the next id.
		log.Println("lastid")
		um.uid[mac] = um.lastid
		// Increment for the next found address.
		um.lastid += 1
		// Grab a manufacturer for this MAC
		log.Println("mfg")
		um.mfg[mac] = api.Mac_to_mfg(cfg, mac)
		// Say when we saw it.
		um.timestamp[mac] = time.Now()
	} else {
		// If this address is already known, update
		// when we last saw it.
		um.timestamp[mac] = time.Now()
	}
}

func (um uniqueMapping) removeOldMappings(window int) {
	n := time.Now()
	remove := make([]string, 0)
	// Find everything we need to remove.
	for _, mac := range um.mfg {
		diff := n.Sub(um.timestamp[mac])
		// Is it further in the past than our window (in minutes)?
		if int(diff.Minutes()) > window {
			remove = append(remove, mac)
		}
	}
	// Remove everything that's old.
	for _, mac := range remove {
		delete(um.uid, mac)
		delete(um.mfg, mac)
		delete(um.timestamp, mac)
	}
}

func (um uniqueMapping) asUserMappings() map[model.UserMapping]int {
	h := make(map[model.UserMapping]int)
	n := time.Now()

	for _, mac := range um.mfg {
		diff := n.Sub(um.timestamp[mac])
		userm := model.UserMapping{}
		userm.Id = um.uid[mac]
		userm.Mfg = um.mfg[mac]
		h[userm] = int(diff.Minutes())
	}

	return h
}

/*
 * 1. Gather the MAC addresses that were seen.
 * 2. Map those to UIDs of mfg:counter
 * 3. Store what was seen into a DB with a timestamp.
 *    3a. If the UID exists, set the timestamp to now.
 *    3b. If the UID does not, insert it.
 * 4. If anything is older than the disconnection_window, remove it.
 * 5. Report this UID:timestamp pairing.
 */

func AlgorithmTwo(ka *csp.Keepalive, cfg *config.Config, in <-chan []string, out chan<- map[model.UserMapping]int, kill <-chan bool) {
	log.Println("Starting AlgorithmTwo")
	// This is our "tracking database"
	um := &uniqueMapping{}
	log.Println(um)
	um.init()
	log.Println(um)
	if um.mfg == nil {
		log.Println("um.mfg is nil")
	}
	// If we are running live, the kill channel is `nil`.
	// When we are live, THEN init the ping/pong.
	testing := true
	if kill == nil {
		testing = false
	}
	var ping chan interface{} = nil
	var pong chan interface{} = nil
	if !testing {
		ping, pong = ka.Subscribe("AlgorithmTwo", 5)
		log.Println("a2: initialized keepalive")
	}
	for {
		select {
		case <-kill:
			log.Println("a2: exiting")
			return
		case <-ping:
			pong <- "AlgorithmTwo"
		case arr := <-in:
			// We get in a list of MAC addresses. Create mappings.
			// Timestamp everything as we see it, new or old.
			for _, mac := range arr {
				um.updateMapping(cfg, mac)
			}
			// Now, filter old things out
			um.removeOldMappings(cfg.Monitoring.UniquenessWindow)
			// Get the mappings as UserMappings, and send them out
			out <- um.asUserMappings()
		}
	}

}
