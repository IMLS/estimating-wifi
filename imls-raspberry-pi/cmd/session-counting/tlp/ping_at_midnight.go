package tlp

import (
	"log"

	"github.com/robfig/cron/v3"
	"gsa.gov/18f/config"
)

func PingAtMidnight(ka *Keepalive, cfg *config.Config, ch_reset chan<- Ping, ch_kill <-chan Ping) {
	if config.Verbose {
		log.Println("Starting PingAtMidnight")
	}

	// ch_kill will be nil in production
	var ping, pong chan interface{} = nil, nil
	if ch_kill == nil {
		ping, pong = ka.Subscribe("PingAtMidnight", 30)
	}

	// Use the cron library to send out the pings.
	// How to kill this in testing? Perhaps we don't...
	c := cron.New()
	_, err := c.AddFunc(cfg.Local.Crontab, func() {
		ch_reset <- Ping{}
	})
	if err != nil {
		log.Println("cron: could not set up crontab entry")
		log.Fatal(err.Error())
	}
	c.Start()

	for {
		select {
		case <-ch_kill:
			if config.Verbose {
				log.Println("Exiting PingAtMidnight")
			}
			// Stop the cron scheduler.
			c.Stop()
			return
		case <-ping:
			pong <- "PingAtMidnight"
		}

	}
}
