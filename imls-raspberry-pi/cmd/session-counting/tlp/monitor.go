package tlp

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"gsa.gov/18f/config"
	"gsa.gov/18f/http"
	"gsa.gov/18f/session-counter/constants"
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
	eventLogger *http.EventLogger
}

func NewKeepalive(cfg *config.Config) *Keepalive {
	el := http.NewEventLogger(cfg)

	return &Keepalive{
		publishCh:   make(chan interface{}, 1),
		subCh:       make(chan resp, 1),
		eventLogger: el,
	}
}

// UPDATE 20210402 MCJ
// We're just going to exit(-1) if things go bad.
// systemd notify does not seem to be working.
// Given that the github repos has build failures...
// I'm going to stop trying to debug this.
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
						// log.Printf("Pong from %v", procs[c].id)
					case <-time.After(procs[c].timeout):
						b.eventLogger.Log("keepalive_timeout",
							map[string]string{
								"process_id":      fmt.Sprint(procs[c].id),
								"process_timeout": fmt.Sprint(procs[c].timeout)})

						if config.Verbose {
							log.Printf("TIMEOUT [%v :: %v]\n", procs[c].id, procs[c].timeout)
						}
						processTimedOut = true
					}
				}(ch)
			}
		// Lets check
		case <-time.After(interval):
			if processTimedOut {
				// If we timed out, exit. Hope systemd restarts us.
				if config.Verbose {
					log.Println("Exiting after", interval, "seconds. Bye!")
				}
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
