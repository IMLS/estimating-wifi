package tlp

import (
	"fmt"
	"log"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"gsa.gov/18f/session-counter/api"
	"gsa.gov/18f/session-counter/config"
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
	eventLogger *api.EventLogger
}

func NewKeepalive(cfg *config.Config) *Keepalive {
	el := api.NewEventLogger(cfg)

	return &Keepalive{
		publishCh:   make(chan interface{}, 1),
		subCh:       make(chan resp, 1),
		eventLogger: el,
	}
}

func (b *Keepalive) Start() {

	// Internal state for the broker.
	procs := make(map[chan interface{}]resp)

	// Notify that we are ready.
	// https://vincent.bernat.ch/en/blog/2017-systemd-golang
	daemon.SdNotify(false, daemon.SdNotifyReady)
	interval, err := daemon.SdWatchdogEnabled(false)
	if err != nil {
		log.Println("monitor: unable to enable watchdog")
		log.Fatal(err)
	}
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
					// log.Printf("Pinging %v\n", procs[c].id)
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

						log.Printf("TIMEOUT [%v :: %v]\n", procs[c].id, procs[c].timeout)
						processTimedOut = true
					}
				}(ch)
			}
		case <-time.After(interval / 3):
			if !processTimedOut {
				// log.Println("... kicking the dog")
				daemon.SdNotify(false, daemon.SdNotifyWatchdog)
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
