package tlp

import (
	"log"
	"sort"
	"strings"

	"gsa.gov/18f/session-counter/api"
	"gsa.gov/18f/session-counter/config"
	"gsa.gov/18f/session-counter/csp"
	"gsa.gov/18f/session-counter/model"
)

func RawToUids(ka *csp.Keepalive, cfg *config.Config, in <-chan map[string]int, out chan<- map[model.UserMapping]int, kill <-chan bool) {
	log.Println("Starting rawToUids")

	// If we are running live, the kill channel is `nil`.
	// When we are live, THEN init the ping/pong.
	testing := true
	if kill == nil {
		testing = false
	}
	var ping chan interface{} = nil
	var pong chan interface{} = nil
	if !testing {
		ping, pong = ka.Subscribe("rawToUids", 5)
		log.Println("rtu: initialized keepalive")
	}

	macToNdx := make(map[string]int)
	ndxToMac := make(map[int]string)

	// The ndx, or nextId, needs to be maintained as a "monotonically increasing"
	// value for the life of a session-counter run.
	nextId := 0
	// To track who has overstayed their disconnection window.
	uniq := make(map[model.UserMapping]int)
	disco := make(map[model.UserMapping]int)

	for {
		select {
		case <-kill:
			log.Println("rtu: exiting")
			return
		case <-ping:
			pong <- "rawToUids"
		case m := <-in:
			if testing {
				log.Println("rtu: received map: ", m)
			}
			// For each incoming address, decide if it is already in our map.
			// If it is, do nothing. If not, give that mac address a new id.
			// FIXME: Traversing in order, primarily for consistency in testing.
			// This should not be a big enough performance issue to matter for production.
			sortedaddrs := make([]string, 0)
			for k := range m {
				sortedaddrs = append(sortedaddrs, k)
			}
			sort.Strings(sortedaddrs)
			for _, addr := range sortedaddrs {
				addr = strings.ToLower(addr)
				_, found := macToNdx[addr]
				if testing {
					log.Printf("rtu: [%v :: %v]\n", addr, found)
				}
				if !found {
					if testing {
						log.Printf("rtu: adding newmap [%v <- %v]\n", addr, nextId)
					}
					macToNdx[addr] = nextId
					nextId += 1
				}
			}
			// Now, build a new mapping for sending down the pipeline.
			// That mapping will be their user id and the device manufacturer
			// and we will keep the "count" of the number of packets that WS saw.
			// That number is probably not meaningful, but we'll hold it for a moment.
			newMapping := make(map[model.UserMapping]int)

			for oldaddr, v := range m {
				mfg := api.Mac_to_mfg(cfg, oldaddr)
				um := model.UserMapping{Mfg: mfg, Id: macToNdx[oldaddr]}
				// log.Println("rtu: newmap ", um, " to ", v)
				newMapping[um] = v
				// If you just arrived to be mapped, you by
				// definition have a 0 uniqueness window ticker
				// and a 0 disco ticker.
				uniq[um] = 0
				disco[um] = 0
			}

			// Everyone we do *not* see has their time bumped.
			for k := range uniq {
				_, here := newMapping[k]
				if !here {
					uniq[k] = uniq[k] + 1
					disco[k] = disco[k] + 1
				}
			}

			// Now, we have to do some munging. A new map is needed.
			// We begin with everyone who is still in the uniqueness window.
			sendmap := make(map[model.UserMapping]int)
			for k, v := range uniq {
				sendmap[k] = v
			}

			// Now look at everyone we are uniquely tracking.
			// We'll clean up our sendmap based on this set.
			for k, v := range uniq {
				// Anyone who disconnected should not be communicated in the sendmap.
				if v >= cfg.Monitoring.DisconnectionWindow {
					delete(sendmap, k)
				}

				// Anyone who overstays our unique tracking window gets
				// removed *completely*. If they show up again, they are considered
				// a new device.
				if v >= cfg.Monitoring.UniquenessWindow {
					if testing {
						log.Printf("%v [%v] no longer uniq\n", k, v)
					}
					delete(sendmap, k)
					delete(disco, k)
					delete(uniq, k)
					// And, clean up the other state we have lying around
					addr := ndxToMac[k.Id]
					delete(ndxToMac, k.Id)
					delete(macToNdx, addr)
				}
			}

			out <- sendmap

		}
	}
}
