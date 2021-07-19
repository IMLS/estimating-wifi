package tlp

import (
	"github.com/robfig/cron/v3"
	"gsa.gov/18f/internal/logwrapper"
	"gsa.gov/18f/internal/state"
)

func PingAtMidnight(ka *Keepalive, rb *ResetBroker, kb *KillBroker) {
	cfg := state.GetConfig()
	lw := logwrapper.NewLogger(nil)
	lw.Debug("starting PingAtMidnight")
	var ping, pong chan interface{} = nil, nil
	var chKill chan interface{} = nil
	if kb != nil {
		chKill = kb.Subscribe()
	} else {
		ping, pong = ka.Subscribe("PingAtMidnight", 30)
	}

	// Use the cron library to send out the pings.
	// Publish a message on the reset broker.
	c := cron.New()
	_, err := c.AddFunc(cfg.Monitoring.ResetCron, func() {
		rb.Publish(Ping{})
	})
	if err != nil {
		lw.Error("could not set up crontab entry")
		lw.Fatal(err.Error())
	}
	c.Start()

	for {
		select {
		case <-chKill:
			lw.Debug("exiting PingAtMidnight")
			// Stop the cron scheduler.
			c.Stop()
			return
		case <-ping:
			pong <- "PingAtMidnight"
		}

	}
}
