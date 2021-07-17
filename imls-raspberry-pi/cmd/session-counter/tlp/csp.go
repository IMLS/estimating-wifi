package tlp

import (
	"sync"

	"gsa.gov/18f/internal/logwrapper"
	"gsa.gov/18f/internal/state"
)

type Generic interface {
	Noop() bool
}

type Ping struct {
}

// ParDeltaTempDB reads in from the input channel, and in parallel, writes out
// the value to the two output channels.
func ParDeltaTempDB(kb *KillBroker,
	chResetIn <-chan *state.TempDB,
	chsResetOut ...chan *state.TempDB) {
	// Block waiting for a message
	// It will be the zeroth channel in the group.
	lw := logwrapper.NewLogger(nil)
	lw.Debug("starting ParDelta")

	var chKill chan interface{} = nil
	if kb != nil {
		chKill = kb.Subscribe()
	}

	for {
		select {
		case <-chKill:
			lw.Debug("exiting ParDelta")
			return
		case v := <-chResetIn:
			var wg sync.WaitGroup
			// Launch two goroutines.
			wg.Add(len(chsResetOut))
			// Don't send to the input channel!
			for ndx := 0; ndx < len(chsResetOut); ndx++ {
				go func(i int) {
					chsResetOut[i] <- v
					wg.Done()
				}(ndx)
			}
			// Wait for both to complete
			wg.Wait()
		}
	}
}

// ParDeltaPing reads in from the input channel, and in parallel, writes out the
// value to the two output channels.
func ParDeltaPing(kb *Broker, chResetIn <-chan Ping, chsResetOut ...chan Ping) {
	// Block waiting for a message
	// It will be the zeroth channel in the group.
	lw := logwrapper.NewLogger(nil)
	lw.Debug("starting ParDelta")

	var chKill chan interface{} = nil
	if kb != nil {
		chKill = kb.Subscribe()
	}

	for {
		select {
		case <-chKill:
			lw.Debug("exiting ParDelta")
			return
		case v := <-chResetIn:
			var wg sync.WaitGroup
			// Launch two goroutines.
			wg.Add(len(chsResetOut))
			// Don't send to the input channel!
			for ndx := 0; ndx < len(chsResetOut); ndx++ {
				go func(i int) {
					chsResetOut[i] <- v
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
