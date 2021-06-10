package tlp

import "sync"

type Ping struct {
}

// In an infinite loop, we read in from the input channel, and
// in parallel, write out the value to the two output channels.
func ParDeltaBool(chs_reset ...chan Ping) {
	for {
		var wg sync.WaitGroup
		// Block waiting for a message
		// It will be the zeroth channel in the group.
		val := <-chs_reset[0]
		// Launch two goroutines.
		wg.Add(len(chs_reset) - 1)
		// Don't send to the input channel!
		for ndx := 1; ndx < len(chs_reset); ndx++ {
			go func(i int) {
				chs_reset[i] <- val
				wg.Done()
			}(ndx)
		}
		// Wait for both to complete
		wg.Wait()
	}
}

func Blackhole(ch <-chan Ping) {
	for {
		<-ch
	}
}
