package tlp

import (
	"log"

	"github.com/robfig/cron/v3"
)

/* PROCESS tick
 * communicates out on the channel `ch` once
 * per second.
 */
func Tick(ka *Keepalive, ch chan<- bool) {
	log.Println("Starting tick")
	ping, pong := ka.Subscribe("tick", 2)
	// What is the best way to drive a 1-second tick?

	c := cron.New()
	_, err := c.AddFunc("*/1 * * * *", func() {
		ch <- true
	})
	if err != nil {
		log.Println("cron: could not set up crontab entry")
		log.Fatal(err.Error())
	}
	c.Start()

	for {
		<-ping
		pong <- "tick"
	}
}
