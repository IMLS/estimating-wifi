package csp

import (
	"log"
	"time"
)

// https://stackoverflow.com/questions/36417199/how-to-broadcast-message-using-channel

type resp struct {
	uid     chan interface{}
	respCh  chan interface{}
	id      string
	timeout time.Duration
}

type Keepalive struct {
	publishCh chan interface{}
	respCh    chan resp
	subCh     chan chan interface{}
}

func NewKeepalive() *Keepalive {
	return &Keepalive{
		publishCh: make(chan interface{}, 1),
		respCh:    make(chan resp, 1),
		subCh:     make(chan chan interface{}, 1),
	}
}

func (b *Keepalive) Start() {
	pings := map[chan interface{}]struct{}{}
	pongs := make(map[chan interface{}]resp)

	for {
		select {
		case msgCh := <-b.subCh:
			pings[msgCh] = struct{}{}
		case r := <-b.respCh:
			pongs[r.uid] = r
		case msg := <-b.publishCh:
			// Send all the pings
			for ch := range pings {
				select {
				case ch <- msg:
				default:
				}
			}

			// Listen for all of the pings with their respective timeouts.
			// If any fail, then exit. Wait for all of them.
			for ch := range pings {
				go func(cid chan interface{}) {
					select {
					case id := <-pongs[cid].respCh:
						log.Printf("Pong from %v", id)
					case <-time.After(pongs[cid].timeout):
						log.Printf("Timeout: %v after %v", pongs[cid].id, pongs[cid].timeout)
						log.Fatalf("FAIL BECAUSE OF %v", pongs[cid].id)
					}
				}(ch)
			}
		}
	}
}

func (b *Keepalive) Subscribe(id string, timeout int) (chan interface{}, chan interface{}) {
	// The message on which we send a ping has a one-slot buffer
	pingCh := make(chan interface{}, 1)
	b.subCh <- pingCh
	pongCh := make(chan interface{}, 1)
	b.respCh <- resp{pingCh, pongCh, id, time.Duration(time.Duration(timeout) * time.Second)}
	return pingCh, pongCh
}

func (b *Keepalive) Publish(msg interface{}) {
	b.publishCh <- msg

}
