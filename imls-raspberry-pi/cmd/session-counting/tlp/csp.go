package tlp

import (
	"sync"

	"gsa.gov/18f/logwrapper"
	"gsa.gov/18f/state"
)

type Generic interface {
	Noop() bool
}

type Ping struct {
}

// In an infinite loop, we read in from the input channel, and
// in parallel, write out the value to the two output channels.
func ParDeltaTempDB(kb *KillBroker,
	ch_reset_in <-chan *state.TempDB,
	chs_reset_out ...chan *state.TempDB) {
	// Block waiting for a message
	// It will be the zeroth channel in the group.
	lw := logwrapper.NewLogger(nil)
	lw.Debug("starting ParDelta")

	var ch_kill chan interface{} = nil
	if kb != nil {
		ch_kill = kb.Subscribe()
	}

	for {
		select {
		case <-ch_kill:
			lw.Debug("exiting ParDelta")
			return
		case v := <-ch_reset_in:
			var wg sync.WaitGroup
			// Launch two goroutines.
			wg.Add(len(chs_reset_out))
			// Don't send to the input channel!
			for ndx := 0; ndx < len(chs_reset_out); ndx++ {
				go func(i int) {
					chs_reset_out[i] <- v
					wg.Done()
				}(ndx)
			}
			// Wait for both to complete
			wg.Wait()
		}
	}
}

// In an infinite loop, we read in from the input channel, and
// in parallel, write out the value to the two output channels.
func ParDeltaPing(kb *Broker, ch_reset_in <-chan Ping, chs_reset_out ...chan Ping) {
	// Block waiting for a message
	// It will be the zeroth channel in the group.
	lw := logwrapper.NewLogger(nil)
	lw.Debug("starting ParDelta")

	var ch_kill chan interface{} = nil
	if kb != nil {
		ch_kill = kb.Subscribe()
	}

	for {
		select {
		case <-ch_kill:
			lw.Debug("exiting ParDelta")
			return
		case v := <-ch_reset_in:
			var wg sync.WaitGroup
			// Launch two goroutines.
			wg.Add(len(chs_reset_out))
			// Don't send to the input channel!
			for ndx := 0; ndx < len(chs_reset_out); ndx++ {
				go func(i int) {
					chs_reset_out[i] <- v
					wg.Done()
				}(ndx)
			}
			// Wait for both to complete
			wg.Wait()
		}
	}
}

func Blackhole(ch <-chan Ping) {
	for {
		<-ch
	}
}
