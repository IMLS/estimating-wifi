package tlp

// A small broker for kill messages.
// This lets us distribute a shutdown to the process network.
// Primarily used in testing (with accelerated time).
// Must be threaded through TLP regardless.

type Broker struct {
	chPub chan interface{}
	chSub chan chan interface{}
}

type ResetBroker struct {
	*Broker
}

type KillBroker struct {
	*Broker
}

func NewKillBroker() *KillBroker {
	return &KillBroker{NewBroker()}
}

func NewResetBroker() *ResetBroker {
	return &ResetBroker{NewBroker()}
}

func NewBroker() *Broker {
	return &Broker{
		chPub: make(chan interface{}, 1),
		chSub: make(chan chan interface{}, 1),
	}
}

func (b *Broker) Start() {
	subs := map[chan interface{}]struct{}{}
	for {
		select {
		case chMsg := <-b.chSub:
			subs[chMsg] = struct{}{}
		case msg := <-b.chPub:
			for chMsg := range subs {
				select {
				case chMsg <- msg:
				default:
				}
			}
		}
	}
}

func (b *Broker) Subscribe() chan interface{} {
	chMsg := make(chan interface{}, 5)
	b.chSub <- chMsg
	return chMsg
}

func (b *Broker) Publish(msg interface{}) {
	b.chPub <- msg
}
