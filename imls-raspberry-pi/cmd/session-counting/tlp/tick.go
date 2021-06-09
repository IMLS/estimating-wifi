package tlp

import (
	"log"
	"time"
)

/* PROCESS tick
 * communicates out on the channel `ch` once
 * per second.
 */
func Tick(ka *Keepalive, ch chan<- bool) {
	log.Println("Starting tick")
	ping, pong := ka.Subscribe("tick", 2)
	// What is the best way to drive a 1-second tick?
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-ping:
			pong <- "tick"

		// FIXME: This drifts?
		//case <-time.After(1 * time.Second):
		// MCJ: Is this better?
		case <-ticker.C:
			ch <- true
		}
	}
}
