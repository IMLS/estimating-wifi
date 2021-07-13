package tlp

// A small broker for kill messages.
// This lets us distribute a shutdown to the process network.
// Primarily used in testing (with accelerated time).
// Must be threaded through TLP regardless.

type Broker struct {
	ch_pub chan interface{}
	ch_sub chan chan interface{}
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
		ch_pub: make(chan interface{}, 1),
		ch_sub: make(chan chan interface{}, 1),
	}
}

func (b *Broker) Start() {
	subs := map[chan interface{}]struct{}{}
	for {
		select {
		case ch_msg := <-b.ch_sub:
			subs[ch_msg] = struct{}{}
		case msg := <-b.ch_pub:
			for ch_msg := range subs {
				select {
				case ch_msg <- msg:
				default:
				}
			}
		}
	}
}

func (b *Broker) Subscribe() chan interface{} {
	ch_msg := make(chan interface{}, 5)
	b.ch_sub <- ch_msg
	return ch_msg
}

func (b *Broker) Publish(msg interface{}) {
	b.ch_pub <- msg
}
