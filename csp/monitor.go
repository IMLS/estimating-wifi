package csp

import (
	"log"
	"time"
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
	publishCh chan interface{}
	subCh     chan resp
}

func NewKeepalive() *Keepalive {
	return &Keepalive{
		publishCh: make(chan interface{}, 1),
		subCh:     make(chan resp, 1),
	}
}

func (b *Keepalive) Start() {
	// Internal state for the broker.
	procs := make(map[chan interface{}]resp)

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
					// We'll log an error and die, so that we can be restarted.
					select {
					case <-procs[c].pongCh:
						// log.Printf("Pong from %v", procs[c].id)
					case <-time.After(procs[c].timeout):
						log.Fatalf("TIMEOUT [%v :: %v]", procs[c].id, procs[c].timeout)
					}
				}(ch)
			}
		}
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
