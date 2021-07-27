package tlp

import (
	"time"

	"github.com/robfig/cron/v3"
	"gsa.gov/18f/internal/logwrapper"
)

func TickEveryMinute(kb *KillBroker, out chan<- bool) {
	lw := logwrapper.NewLogger(nil)
	lw.Debug("STARTING")
	var chKill chan interface{} = nil
	if kb != nil {
		chKill = kb.Subscribe()
	}

	c := cron.New()
	_, err := c.AddFunc("*/1 * * * *", func() {
		// lw := logwrapper.NewLogger(nil)
		//lw.Debug("tock every minute")
		out <- true
	})
	if err != nil {
		lw.Info("cron: could not set up crontab entry")
		lw.Fatal(err.Error())
	}
	c.Start()

	// Wait for us to be killed, shut down the ticker, and return.
	for {
		<-chKill
		c.Stop()
		lw.Debug("EXITING")
		return
	}
}

// ******* WARNING
// This is only used in testing. It lets us drive a variable clock.
func TickConstantly(kb *KillBroker, delay_millis int, out chan<- Ping) {
	lw := logwrapper.NewLogger(nil)
	lw.Debug("STARTING")
	var chKill chan interface{} = nil
	if kb != nil {
		chKill = kb.Subscribe()
	}

	for {
		select {
		case <-chKill:
			lw.Debug("EXITING")
			return
		default:
			if delay_millis != 0 {
				time.Sleep(time.Duration(delay_millis) * time.Millisecond)
			}
			out <- Ping{}
		}
	}
}
