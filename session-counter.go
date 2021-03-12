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

/* PROCESS run_wireshark
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

/* PROCESS show_map
 * Walks a hashmap and prints it.
 */
func show_map(in chan map[string]int) {
	for {
		macmap := <-in
		for k, v := range macmap {
			fmt.Println(k, " <- ", v)
		}
	}
}

/* FUNC check_env_vars
 * Checks to see if the username and password for
 * working with Directus is in memory.
 * If not, it quits.
 */
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

/* PROC mac_to_mfg
 * Takes in a hashmap of MAC addresses and counts, and passes on a hashmap
 * of manufacturer IDs and counts.
 * Uses "unknown" for all unknown manufacturers.
 */
func mac_to_Entry(cfg model.Config, macmap chan map[string]int, mfgmap chan map[string]model.Entry) {
	for {
		mfgs := make(map[string]model.Entry)
		for mac, count := range <-macmap {
			mfg := api.Mac_to_mfg(cfg, mac)
			mfgs[mac] = model.Entry{MAC: mac, Mfg: mfg, Count: count}
		}
		mfgmap <- mfgs
	}
}

/* PROC report_map
 * Takes a hashmap of [mfg id : count] and POSTs
 * each one to the server individually. We have no bulk insert.
 */
func report_map(cfg model.Config, mfgs chan map[string]model.Entry) {
	for {
		m := <-mfgs
		for _, entry := range m {
			api.Report_mfg(cfg, entry)
		}
	}
}

func main() {
	check_env_vars()
	cfgPtr := flag.String("config", "config.yaml", "config file")

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
	mfg := make(chan map[string]model.Entry)

	go tick(ch_sec)
	go tock_every_n("min", 10, ch_sec, ch_nsec)
	go run_wireshark(cfg, ch_nsec, ch_macs)
	go delta_map(ch_macs, ch_m1, ch_m2)
	go mac_to_Entry(cfg, ch_m1, mfg)
	go report_map(cfg, mfg)
	go show_map(ch_m2)

	// Wait forever.
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
