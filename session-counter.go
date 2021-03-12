package main

import ( 
    "fmt"
	"flag"
	"os"
    "time"
	"sync"
	"gsa.gov/18f/session-counter/tshark"
	"gsa.gov/18f/session-counter/constants"
) 

func tick (ch chan bool) {
	for {
		time.Sleep(1 * time.Second)
		ch <- true
	}
}

func tock_every_n (n int, in chan bool, out chan bool) {
	var counter int = 0
	for {
		<- in
		counter = counter + 1
		if (counter == n) {
			fmt.Println("tock")
			counter = 0
			out <- true
		}
	}
}

func get_ip_addrs (in chan bool, out chan map[string]int) {
	for {
		<-in
		macmap := tshark.Tshark("wlan1", 10)
		out <- macmap
	}
}

func show_map (in chan map[string]int) {
	for {
		macmap := <- in
		for k, v := range macmap {
			fmt.Println(k, " <- ", v)
		}
	}
}

func check_env_vars () {
	if os.Getenv(constants.EnvUsername) == "" {
		fmt.Printf("%s must be set in the env!\n", constants.EnvUsername)
		os.Exit(constants.ExitNoUsername)
	}
	if os.Getenv(constants.EnvPassword) == "" {
		fmt.Printf("%s must be set in the env!\n", constants.EnvPassword)
		os.Exit(constants.ExitNoPassword)
	}
}

func main() {
	check_env_vars()
	
	obsecPtr := flag.Int("observe", 10, "seconds to observe wifi traffic")
	//reptMinsPtr := flag.Int("reptMins", 10, "minutes to observe")
	//appearThreshPtr := flag.Int("appearThresh", 3, "number of times MAC must appear")
	flag.Parse()
	
	ch_sec := make(chan bool)
	ch_nsec := make(chan bool)
	ch_macs := make(chan map[string]int)

	go tick(ch_sec)
	go tock_every_n(*obsecPtr, ch_sec, ch_nsec)
	go get_ip_addrs(ch_nsec, ch_macs)
	go show_map(ch_macs)

	// Wait forever.
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
