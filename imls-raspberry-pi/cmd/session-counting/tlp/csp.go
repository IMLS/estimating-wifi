package tlp

import (
	"log"
	"sync"
)

type Ping struct {
}

// In an infinite loop, we read in from the input channel, and
// in parallel, write out the value to the two output channels.
func ParDelta(ch_kill <-chan Ping, chs_reset ...chan Ping) {
	// Block waiting for a message
	// It will be the zeroth channel in the group.

	for {
		select {
		case <-ch_kill:
			log.Println("Exiting ParDelta")
			return
		case <-chs_reset[0]:
			var wg sync.WaitGroup
			// Launch two goroutines.
			wg.Add(len(chs_reset) - 1)
			// Don't send to the input channel!
			for ndx := 1; ndx < len(chs_reset); ndx++ {
				go func(i int) {
					chs_reset[i] <- Ping{}
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
