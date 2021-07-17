package tlp

import (
	"github.com/robfig/cron/v3"
	"gsa.gov/18f/internal/logwrapper"
)

func TockEveryMinute(ka *Keepalive, kb *KillBroker, out chan<- bool) {
	lw := logwrapper.NewLogger(nil)
	lw.Debug("starting TockEveryMinute")
	var ping, pong chan interface{} = nil, nil
	var chKill chan interface{} = nil
	if kb != nil {
		chKill = kb.Subscribe()
	} else {
		ping, pong = ka.Subscribe("TockEveryMinute", 5)
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

	for {
		select {
		case <-ping:
			pong <- "TockEveryMinute"
		case <-chKill:
			lw.Debug("Exiting TockEveryN")
			return
		}
	}
}

// ******* WARNING
// This is only used in testing. It lets us drive a variable clock.

/* PROCESS tockEveryN
 * consumes a tag (for logging purposes) as well as
 * a driving `tick` on `in`. Every `n` ticks, it outputs
 * a boolean `tock` on the channel `out`.
 * When `in` is every second, and `n` is 60, it turns
 * a stream of second ticks into minute `tocks`.
 */
func TockEveryN(ka *Keepalive, kb *KillBroker, n int, in <-chan bool, out chan<- bool) {
	lw := logwrapper.NewLogger(nil)
	lw.Debug("starting TockEveryN")
	var ping, pong chan interface{} = nil, nil
	var chKill chan interface{} = nil
	if kb != nil {
		chKill = kb.Subscribe()
	} else {
		ping, pong = ka.Subscribe("TockEveryN", 5)
	}

	var counter int = 0
	for {
		select {
		case <-ping:
			pong <- "tock"
		case <-chKill:
			lw.Debug("exiting TockEveryN")
			return

		case <-in:
			counter = counter + 1
			if counter == n {
				//lw.Info("tickN", counter)
				counter = 0
				out <- true
			}
		}
	}
}
