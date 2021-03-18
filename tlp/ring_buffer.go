package tlp

import (
	"log"

	"gsa.gov/18f/session-counter/csp"
	"gsa.gov/18f/session-counter/model"
)

func RingBuffer(ka *csp.Keepalive, cfg *model.Config, in <-chan map[string]int, out chan<- map[string]int) {
	log.Println("Starting ringBuffer")
	ping, pong := ka.Subscribe("ringBuffer", 3)

	// Nothing in the buffer, capacity = number of rounds
	buffer := make([]map[string]int, cfg.Wireshark.Rounds)
	for ndx := 0; ndx < cap(buffer); ndx++ {
		buffer[ndx] = nil
	}
	// Circular index.
	ring_ndx := 0

	for {
		select {
		case <-ping:
			pong <- "ringBuffer"

		case buffer[ring_ndx] = <-in:
			// Read in to the most recent buffer index.
			// Zero out a map for counting how many times
			// MAC addresses appear.
			total := make(map[string]int)

			// Count everything in the ring. The ring is right-sized
			// to the window we're interested in.
			filled_slots := 0
			for _, m := range buffer {
				if m != nil {
					filled_slots += 1
					for mac := range m {
						cnt, ok := total[mac]
						if ok {
							total[mac] = cnt + 1
						} else {
							total[mac] = 1
						}
					}
				}
			}

			// If we have filled enough slots to be "countable,"
			// we should go through and see which MAC addresses appeared
			// enough times to be "worth reporting."
			if filled_slots == cfg.Wireshark.Rounds {
				// Filter out the ones that don't make the cut.
				var filter []string
				for mac, count := range total {
					if count < cfg.Wireshark.Threshold {
						filter = append(filter, mac)
					}
				}
				for _, f := range filter {
					delete(total, f)
				}
				// These are the MAC addresses that passed our
				// threshold of `threshold` in `rounds` cycles.
				out <- total
			}

			// Bump the index. Overwrite old values.
			// Then, wait for the next hash to come in.
			ring_ndx = (ring_ndx + 1) % cfg.Wireshark.Rounds
		}
	}
}
