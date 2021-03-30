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

	for {
		select {
		case <-ping:
			pong <- "tick"

		case <-time.After(1 * time.Second):
			// Drive the 1 second ticker
			ch <- true
		}
	}
}
