package tlp

import (
	"log"

	"github.com/robfig/cron/v3"
	"gsa.gov/18f/config"
)

func PingAtMidnight(ka *Keepalive, cfg *config.Config, ch_reset chan<- Ping) {
	log.Println("Starting PingAtMidnight")
	ping, pong := ka.Subscribe("PingAtMidnight", 30)
	// For event logging
	// el := http.NewEventLogger(cfg)

	// Use the cron library to send out the pings.
	c := cron.New()
	c.AddFunc(cfg.Local.Crontab, func() {
		ch_reset <- Ping{}
	})
	c.Start()

	for {
		<-ping
		pong <- "PingAtMidnight"
	}
}
