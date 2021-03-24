package tlp

import (
	"fmt"
	"log"
	"time"

	"gsa.gov/18f/session-counter/api"
	"gsa.gov/18f/session-counter/config"
	"gsa.gov/18f/session-counter/csp"
)

type uniqueMappingDB struct {
	lastid    int
	uid       map[string]int
	mfg       map[string]string
	timestamp map[string]string
}

func newUMDB() *uniqueMappingDB {
	// var umdb uniqueMappingDB
	// umdb.lastid = 0
	// umdb.uid = make(map[string]int)
	// umdb.mfg = make(map[string]string)
	// umdb.timestamp = make(map[string]time.Time)
	umdb := &uniqueMappingDB{
		lastid:    0,
		uid:       make(map[string]int),
		mfg:       make(map[string]string),
		timestamp: make(map[string]string)}
	return umdb
}

func (umdb uniqueMappingDB) updateMapping(cfg *config.Config, mac string) {
	_, found := umdb.mfg[mac]

	// If we didn't find it, then we need to add it.
	if !found {
		// Assign the next id.
		umdb.uid[mac] = umdb.lastid
		// Increment for the next found address.
		umdb.lastid += 1
		// Grab a manufacturer for this MAC
		umdb.mfg[mac] = api.Mac_to_mfg(cfg, mac)
		// Say when we saw it.
		now := time.Now().Format(time.RFC3339)
		umdb.timestamp[mac] = now
	} else {
		// If this address is already known, update
		// when we last saw it.
		umdb.timestamp[mac] = time.Now().Format(time.RFC3339)
	}
}

func (umdb uniqueMappingDB) removeOldMappings(window int) {
	now := time.Now()
	remove := make([]string, 0)
	// Find everything we need to remove.
	for mac, _ := range umdb.mfg {
		storedtime, _ := time.Parse(time.RFC3339, umdb.timestamp[mac])
		diff := now.Sub(storedtime)
		// Is it further in the past than our window (in minutes)?
		if int(diff.Minutes()) > window {
			log.Println(mac, "is old. removing. diff:", diff.Minutes())
			remove = append(remove, mac)
		}
	}
	// Remove everything that's old.
	for _, mac := range remove {
		delete(umdb.uid, mac)
		delete(umdb.mfg, mac)
		delete(umdb.timestamp, mac)
	}
}

func (umdb uniqueMappingDB) asUserMappings() map[string]int {
	h := make(map[string]int)
	n := time.Now()

	for mac, _ := range umdb.mfg {
		userm := fmt.Sprintf("%v:%d", umdb.mfg[mac], umdb.uid[mac])
		storedtime, _ := time.Parse(time.RFC3339, umdb.timestamp[mac])
		diff := n.Sub(storedtime)
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

func AlgorithmTwo(ka *csp.Keepalive, cfg *config.Config, in <-chan []string, out chan<- map[string]int, kill <-chan bool) {
	log.Println("Starting AlgorithmTwo")
	// This is our "tracking database"
	umdb := newUMDB()

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
				umdb.updateMapping(cfg, mac)
			}
			// Now, filter old things out
			umdb.removeOldMappings(cfg.Monitoring.UniquenessWindow)
			// Get the mappings as UserMappings, and send them out
			out <- umdb.asUserMappings()
		}
	}

}
