package csp

/* PROC delta (int)
 * A `delta` process copies the input from a channel
 * and outputs it to two other channels.
 * This is a PAR delta, because it spawns two anonymous
 * goroutines to do the sends on the two outputs
 * "at the same time."
 */
func delta_int(in chan int, o1 chan int, o2 chan int) {
	for {
		v := <-in
		// PAR DELTA
		go func() { o1 <- v }()
		go func() { o2 <- v }()
	}
}

/* PROC delta (map) */
func delta_map(in chan map[string]int, o1 chan map[string]int, o2 chan map[string]int) {
	for {
		v := <-in
		// PAR DELTA
		go func() { o1 <- v }()
		go func() { o2 <- v }()
	}
}

func Delta_par_three(in <-chan map[string]int, o1 chan<- map[string]int, o2 chan<- map[string]int, o3 chan<- map[string]int) {
	for {
		v := <-in
		go func() { o1 <- v }()
		go func() { o2 <- v }()
		go func() { o3 <- v }()
	}
}
