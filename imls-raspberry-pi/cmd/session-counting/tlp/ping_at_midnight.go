package tlp

import (
	"github.com/robfig/cron/v3"
	"gsa.gov/18f/config"
	"gsa.gov/18f/logwrapper"
)

func PingAtMidnight(ka *Keepalive, cfg *config.Config, rb *Broker, kb *Broker) {
	lw := logwrapper.NewLogger(nil)
	lw.Debug("starting PingAtMidnight")
	var ping, pong chan interface{} = nil, nil
	var ch_kill chan interface{} = nil
	if kb != nil {
		ch_kill = kb.Subscribe()
	} else {
		ping, pong = ka.Subscribe("PingAtMidnight", 30)
	}

	// Use the cron library to send out the pings.
	// Publish a message on the reset broker.
	c := cron.New()
	_, err := c.AddFunc(cfg.Local.Crontab, func() {
		rb.Publish(Ping{})
	})
	if err != nil {
		lw.Error("could not set up crontab entry")
		lw.Fatal(err.Error())
	}
	c.Start()

	for {
		select {
		case <-ch_kill:
			lw.Debug("exiting PingAtMidnight")
			// Stop the cron scheduler.
			c.Stop()
			return
		case <-ping:
			pong <- "PingAtMidnight"
		}

	}
}
