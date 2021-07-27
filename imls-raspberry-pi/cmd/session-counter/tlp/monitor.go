package tlp

import (
	"fmt"
	"os"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"gsa.gov/18f/cmd/session-counter/constants"
	"gsa.gov/18f/internal/logwrapper"
	"gsa.gov/18f/internal/state"
)

// Inspired by the broker pattern found here:
// https://stackoverflow.com/questions/36417199/how-to-broadcast-message-using-channel

type resp struct {
	pingCh  chan interface{}
	pongCh  chan interface{}
	id      string
	timeout time.Duration
}

type Keepalive struct {
	publishCh   chan interface{}
	subCh       chan resp
	eventLogger *logwrapper.StandardLogger
}

func NewKeepalive() *Keepalive {
	lw := logwrapper.NewLogger(nil)

	return &Keepalive{
		publishCh:   make(chan interface{}, 1),
		subCh:       make(chan resp, 1),
		eventLogger: lw,
	}
}

func (b *Keepalive) Start() {

	// Internal state for the broker.
	procs := make(map[chan interface{}]resp)
	// We'll check once every 5 seconds to see if we've died.
	interval := 5 * time.Second

	// Notify that we are ready.
	// https://vincent.bernat.ch/en/blog/2017-systemd-golang
	daemon.SdNotify(false, daemon.SdNotifyReady)
	// Even though the watchdog does not seem to be working,
	// this does tell systemd we're here when running as a simple
	// service. So, we need to keep it.

	processTimedOut := false

	for {
		select {
		// Likewise, also init the response map.
		case r := <-b.subCh:
			procs[r.pingCh] = r
		// When a message is published...
		case msg := <-b.publishCh:
			for ch := range procs {
				// In parallel...
				go func(c chan interface{}) {
					// Ping every process that is subscribed.
					// The channel has a single slot.
					c <- msg
					// Now, wait for the response, or timeout.
					// If we timeout, that means someone didn't reply.
					// We'll log an error and stop notifying systemd.
					// This way, systemd will restart us.
					select {
					case <-procs[c].pongCh:
						// WARNING: This could get very noisy in the log.
						b.eventLogger.Debug("pong from ", procs[c].id)
					case <-state.GetClock().After(procs[c].timeout):
						b.eventLogger.Debug(fmt.Sprintf("TIMEOUT [%v :: %v]\n", procs[c].id, procs[c].timeout))
						processTimedOut = true
					}
				}(ch)
			}
		// Lets check
		case <-state.GetClock().After(interval):
			if processTimedOut {
				// If we timed out, exit. Hope systemd restarts us.
				b.eventLogger.Error(fmt.Sprintf("exiting after %v seconds. Hopefully someone will restart us!", interval))
				os.Exit(constants.ExitProcessTimeout)
			}
		} // end select
	}
}

func (b *Keepalive) Subscribe(id string, timeout int) (chan interface{}, chan interface{}) {
	pingCh := make(chan interface{})
	pongCh := make(chan interface{})
	b.subCh <- resp{pingCh, pongCh, id, time.Duration(time.Duration(timeout) * time.Second)}
	return pingCh, pongCh
}

func (b *Keepalive) Publish(msg interface{}) {
	b.publishCh <- msg

}

// TOP LEVEL PROCESS

func StayinAlive(ka *Keepalive) {
	lw := logwrapper.NewLogger(nil)
	lw.Info("starting keepalive")
	ka.Start()

	var counter int64 = 0
	for {
		// Ping every 30 seconds.
		state.GetClock().Sleep(time.Duration(30) * time.Second)
		ka.Publish(counter)
		counter = counter + 1
	}
}
