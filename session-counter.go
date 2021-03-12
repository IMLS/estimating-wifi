package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
	"gsa.gov/18f/session-counter/api"
	"gsa.gov/18f/session-counter/constants"
	"gsa.gov/18f/session-counter/model"
	"gsa.gov/18f/session-counter/tshark"
)

/* PROCESS tick
 * communicates out on the channel `ch` once
 * per second.
 */
func tick(ch chan bool) {
	for {
		time.Sleep(1 * time.Second)
		ch <- true
	}
}

/* PROCESS tock_every_n
 * consumes a tag (for logging purposes) as well as
 * a driving `tick` on `in`. Every `n` ticks, it outputs
 * a boolean `tock` on the channel `out`.
 * When `in` is every second, and `n` is 60, it turns
 * a stream of second ticks into minute `tocks`.
 */
func tock_every_n(tag string, n int, in chan bool, out chan bool) {
	var counter int = 0
	for {
		<-in
		counter = counter + 1
		if counter == n {
			fmt.Printf("tock: %s %d\n", tag, n)
			counter = 0
			out <- true
		}
	}
}

/* PROCESS get_ip_addrs
 * Runs a subprocess for a duration of OBSERVE_SECONDS.
 * Therefore, this process effectively blocks for that time.
 * Gathers a hashmap of [MAC -> count] values. This hashmap
 * is then communicated out.
 * Empty MAC addresses are filtered out.
 */
func run_wireshark(cfg model.Config, in chan bool, out chan map[string]int) {
	for {
		<-in
		fmt.Println("getting mac addresses")

		macmap := tshark.Tshark(cfg)

		var to_remove []string
		// Mark and remove too-short MAC addresses
		// for removal from the tshark findings.
		for k, _ := range macmap {
			if len(k) < constants.MACLENGTH {
				to_remove = append(to_remove, k)
			}
		}
		for _, s := range to_remove {
			delete(macmap, s)
		}
		// Report out the cleaned MACmap.
		out <- macmap
	}
}

func show_map(in chan map[string]int) {
	for {
		macmap := <-in
		for k, v := range macmap {
			fmt.Println(k, " <- ", v)
		}
	}
}

func check_env_vars() {
	if os.Getenv(constants.EnvUsername) == "" {
		fmt.Printf("%s must be set in the env!\n", constants.EnvUsername)
		os.Exit(constants.ExitNoUsername)
	}
	if os.Getenv(constants.EnvPassword) == "" {
		fmt.Printf("%s must be set in the env!\n", constants.EnvPassword)
		os.Exit(constants.ExitNoPassword)
	}
}

func delta_int(in chan int, o1 chan int, o2 chan int) {
	for {
		v := <-in
		// PAR DELTA
		go func() { o1 <- v }()
		go func() { o2 <- v }()
	}
}

func delta_map(in chan map[string]int, o1 chan map[string]int, o2 chan map[string]int) {
	for {
		v := <-in
		// PAR DELTA
		go func() { o1 <- v }()
		go func() { o2 <- v }()
	}
}

func mac_to_mfg(cfg model.Config, macmap chan map[string]int, mfgmap chan map[string]int) {
	for {
		mfgs := make(map[string]int)
		for mac, count := range <-macmap {
			mfg := api.Mac_to_mfg(cfg, mac)
			mfgs[mfg] = count
		}
		mfgmap <- mfgs
	}
}

func report_map(mfgs chan map[string]int) {
	for {
		m := <-mfgs
		fmt.Println(m)
	}
}

func main() {
	check_env_vars()
	cfgPtr := flag.String("config", "config.yaml", "config file")

	// adapterPtr := flag.String("adapter", "wlan1", "adapter to monitor")
	// windowPtr := flag.Int("window", 10, "window size in minutes")
	// mfgPtr := flag.String("manufacturers", "", "manufacturerer sqlite database")
	//appearThreshPtr := flag.Int("appearThresh", 3, "number of times MAC must appear")
	flag.Parse()
	// FIXME: handle errors
	f, _ := os.Open(*cfgPtr)
	defer f.Close()
	var cfg model.Config
	decoder := yaml.NewDecoder(f)
	_ = decoder.Decode(&cfg)

	ch_sec := make(chan bool)
	ch_nsec := make(chan bool)
	ch_macs := make(chan map[string]int)
	ch_m1 := make(chan map[string]int)
	ch_m2 := make(chan map[string]int)
	mfg := make(chan map[string]int)

	go tick(ch_sec)
	go tock_every_n("min", 10, ch_sec, ch_nsec)
	go run_wireshark(cfg, ch_nsec, ch_macs)
	go delta_map(ch_macs, ch_m1, ch_m2)
	go mac_to_mfg(cfg, ch_m1, mfg)
	go report_map(mfg)
	go show_map(ch_m2)

	// Wait forever.
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
