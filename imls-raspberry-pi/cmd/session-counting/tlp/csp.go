package tlp

import (
	"sync"

	"gsa.gov/18f/logwrapper"
)

type Ping struct {
}

// In an infinite loop, we read in from the input channel, and
// in parallel, write out the value to the two output channels.
func ParDelta(ch_kill <-chan Ping, ch_reset_in chan Ping, chs_reset_out ...chan Ping) {
	// Block waiting for a message
	// It will be the zeroth channel in the group.
	lw := logwrapper.NewLogger(nil)
	lw.Debug("starting ParDelta")

	for {
		select {
		case <-ch_kill:
			lw.Debug("exiting ParDelta")
			return
		case <-ch_reset_in:
			var wg sync.WaitGroup
			// Launch two goroutines.
			wg.Add(len(chs_reset_out))
			// Don't send to the input channel!
			for ndx := 0; ndx < len(chs_reset_out); ndx++ {
				go func(i int) {
					chs_reset_out[i] <- Ping{}
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
