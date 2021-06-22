package tlp

import (
	"log"

	"github.com/robfig/cron/v3"
	"gsa.gov/18f/config"
)

func TockEveryMinute(ka *Keepalive, out chan<- bool, ch_kill <-chan Ping) {
	log.Println("Starting tick")
	ping, pong := ka.Subscribe("TockEveryMinute", 2)
	// What is the best way to drive a 1-second tick?

	c := cron.New()
	_, err := c.AddFunc("*/1 * * * *", func() {
		out <- true
	})
	if err != nil {
		log.Println("cron: could not set up crontab entry")
		log.Fatal(err.Error())
	}
	c.Start()

	for {
		select {
		case <-ping:
			pong <- "tick"
		case <-ch_kill:
			log.Println("Exiting TockEveryN")
			return
		}
	}
}

/* PROCESS tockEveryN
 * consumes a tag (for logging purposes) as well as
 * a driving `tick` on `in`. Every `n` ticks, it outputs
 * a boolean `tock` on the channel `out`.
 * When `in` is every second, and `n` is 60, it turns
 * a stream of second ticks into minute `tocks`.
 */
func TockEveryN(ka *Keepalive, n int, in <-chan bool, out chan<- bool, ch_kill <-chan Ping) {
	log.Println("Starting tockEveryN")
	// We timeout one second beyond the number of ticks we're waiting for

	// ch_kill will be nil in production
	var ping, pong chan interface{} = nil, nil
	if ch_kill == nil {
		ping, pong = ka.Subscribe("tock", 2)
	}

	var counter int = 0
	for {
		select {
		case <-ping:
			pong <- "tock"
		case <-ch_kill:
			log.Println("Exiting TockEveryN")
			return

		case <-in:
			counter = counter + 1
			if counter == n {
				if config.Verbose {
					log.Println("tickN", counter)
				}
				counter = 0
				out <- true
			}
		}
	}
}
