package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
	"gsa.gov/18f/session-counter/api"
	"gsa.gov/18f/session-counter/constants"
	"gsa.gov/18f/session-counter/csp"
	"gsa.gov/18f/session-counter/model"
	"gsa.gov/18f/session-counter/tshark"
)

/* PROCESS tick
 * communicates out on the channel `ch` once
 * per second.
 */
func tick(ka *csp.Keepalive, ch chan<- bool) {
	log.Println("Starting tick")
	ping, pong := ka.Subscribe("tick", 2)

	for {
		select {
		case <-ping:
			pong <- "tick"

		case <-time.After(1 * time.Second):
			// Drive the 1 second ticker
			ch <- true
		}
	}
}

/* PROCESS tock_every_n
 * consumes a tag (for logging purposes) as well as
 * a driving `tick` on `in`. Every `n` ticks, it outputs
 * a boolean `tock` on the channel `out`.
 * When `in` is every second, and `n` is 60, it turns
 * a stream of second ticks into minute `tocks`.
 */
func tock_every_n(ka *csp.Keepalive, n int, in chan<- bool, out <-chan bool) {
	log.Println("Starting tock_every_n")
	// We timeout one second beyond the number of ticks we're waiting for
	ping, pong := ka.Subscribe("tock", 2)

	var counter int = 0
	for {
		select {
		case <-ping:
			pong <- "tock"

		case <-in:
			counter = counter + 1
			if counter == n {
				counter = 0
				out <- true
			}
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
func run_wireshark(ka *csp.Keepalive, cfg model.Config, in <-chan bool, out chan<- map[string]int) {
	log.Println("Starting run_wireshark")
	// If we have to wait twice the monitor duration, something broke.
	ping, pong := ka.Subscribe("run_wireshark", cfg.Wireshark.Duration*2)

	for {
		select {

		case <-ping:
			// We ping faster than this process can reply. However, we have a long
			// enough timeout that we will *eventually* catch up with all of the pings.
			pong <- "wireshark"

		case <-in:
			// This will block for [cfg.Wireshark.Duration] seconds.
			macmap := tshark.Tshark(cfg)
			// Mark and remove too-short MAC addresses
			// for removal from the tshark findings.
			var to_remove []string
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

/* PROC mac_to_mfg
 * Takes in a hashmap of MAC addresses and counts, and passes on a hashmap
 * of manufacturer IDs and counts.
 * Uses "unknown" for all unknown manufacturers.
 */
func mac_to_Entry(ka *csp.Keepalive, cfg model.Config, macmap <-chan map[string]int, mfgmap chan<- map[string]model.Entry) {
	log.Println("Starting mac_to_Entry")
	ping, pong := ka.Subscribe("mac_to_Entry", 5)

	for {
		select {
		case <-ping:
			pong <- "mac_to_Entry"

		case mm := <-macmap:
			mfgs := make(map[string]model.Entry)
			for mac, count := range mm {
				mfg := api.Mac_to_mfg(cfg, mac)
				mfgs[mac] = model.Entry{MAC: mac, Mfg: mfg, Count: count}
			}
			mfgmap <- mfgs
		}
	}
}

/* PROC report_map
 * Takes a hashmap of [mfg id : count] and POSTs
 * each one to the server individually. We have no bulk insert.
 */
func report_map(ka *csp.Keepalive, cfg model.Config, mfgs <-chan map[string]model.Entry) {
	log.Println("Starting report_map")
	ping, pong := ka.Subscribe("report_map", 5)

	var count int64 = 0
	http_error_count := 0

	for {
		select {
		case <-ping:
			// Every [cfg.Monitoring.HTTPErrorIntervalMins] this value
			// is reset to zero. If we get too many errors in that number of
			// minutes, then we should fail the next pong request. This will
			// kill the program, and we'll restart.
			if http_error_count < cfg.Monitoring.MaxHTTPErrorCount {
				pong <- "report_map"
			} else {
				log.Printf("report: http_error_count threshold of %d reached\n", http_error_count)
			}

		case m := <-mfgs:
			count = count + 1
			log.Println("reporting: ", count)
			tok, err := api.Get_token(cfg)
			if err != nil {
				log.Println("report: error in token fetch")
				log.Println(err)
				http_error_count = http_error_count + 1
			} else {
				for _, entry := range m {
					go func(entry model.Entry) {
						// Smear the requests out in time.
						time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
						err := api.Report_mfg(cfg, tok, entry)
						if err != nil {
							log.Println("report: results POST failure")
							log.Println(err)
							http_error_count = http_error_count + 1
						}
					}(entry)
				}
				err := api.Report_telemetry(cfg, tok)
				if err != nil {
					log.Println("report: error in telemetry POST")
					log.Println(err)
					http_error_count = http_error_count + 1
				}
			}

		case <-time.After(time.Duration(cfg.Monitoring.HTTPErrorIntervalMins) * time.Minute):
			// If this much time has gone by, go ahead and reset the error count.
			if http_error_count != 0 {
				log.Printf("report: resetting http_error_count from %d to 0\n", http_error_count)
				http_error_count = 0
			}
		}
	}
}

func ring_buffer(ka *csp.Keepalive, cfg model.Config, in <-chan map[string]int, out chan<- map[string]int) {
	log.Println("Starting ring_buffer")
	ping, pong := ka.Subscribe("ring_buffer", 3)

	// Nothing in the buffer, capacity = number of rounds
	buffer := make([]map[string]int, cfg.Wireshark.Rounds)
	for ndx := 0; ndx < cap(buffer); ndx++ {
		buffer[ndx] = nil
	}
	// Circular index.
	ring_ndx := 0

	for {
		select {
		case <-ping:
			pong <- "ring_buffer"

		case buffer[ring_ndx] = <-in:
			// Read in to the most recent buffer index.
			// Zero out a map for counting how many times
			// MAC addresses appear.
			total := make(map[string]int)

			// Count everything in the ring. The ring is right-sized
			// to the window we're interested in.
			filled_slots := 0
			for _, m := range buffer {
				if m != nil {
					filled_slots += 1
					for mac, _ := range m {
						cnt, ok := total[mac]
						if ok {
							total[mac] = cnt + 1
						} else {
							total[mac] = 1
						}
					}
				}
			}

			// If we have filled enough slots to be "countable,"
			// we should go through and see which MAC addresses appeared
			// enough times to be "worth reporting."
			if filled_slots == cfg.Wireshark.Rounds {
				// Filter out the ones that don't make the cut.
				var filter []string
				for mac, count := range total {
					if count < cfg.Wireshark.Threshold {
						filter = append(filter, mac)
					}
				}
				for _, f := range filter {
					delete(total, f)
				}
				// These are the MAC addresses that passed our
				// threshold of `threshold` in `rounds` cycles.
				out <- total
			}

			// Bump the index. Overwrite old values.
			// Then, wait for the next hash to come in.
			ring_ndx = (ring_ndx + 1) % cfg.Wireshark.Rounds
		}
	}
}

func read_config(cfgPtr string) model.Config {

	f, err := os.Open(cfgPtr)
	if err != nil {
		log.Fatal("Could not open configuration file. Exiting.")
	}
	defer f.Close()
	var cfg model.Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatalf("Could not decode YAML:\n%v\n", err)
	}
	return cfg
}

func run(ka *csp.Keepalive, cfg model.Config) {
	log.Println("Starting run")
	// Create channels for process network
	ch_sec := make(chan bool)
	ch_nsec := make(chan bool)
	ch_macs := make(chan map[string]int)
	ch_macs_counted := make(chan map[string]int)
	mfg := make(chan map[string]model.Entry)

	// Run the process network.
	// Driven by a 1s `tick` process.
	// Thread the keepalive through the network
	go tick(ka, ch_sec)
	go tock_every_n(ka, 60, ch_sec, ch_nsec)
	go run_wireshark(ka, cfg, ch_nsec, ch_macs)
	go ring_buffer(ka, cfg, ch_macs, ch_macs_counted)
	go mac_to_Entry(ka, cfg, ch_macs_counted, mfg)
	go report_map(ka, cfg, mfg)
}

func keepalive(ka *csp.Keepalive, cfg model.Config) {
	log.Println("Starting keepalive")
	var counter int64 = 0
	for {
		time.Sleep(time.Duration(cfg.Monitoring.PingInterval) * time.Second)
		ka.Publish(counter)
		counter = counter + 1
	}
}

func main() {
	check_env_vars()
	// FIXME consider turning this into an env var
	cfgPtr := flag.String("config", "config.yaml", "config file")
	flag.Parse()
	cfg := read_config(*cfgPtr)

	ka := csp.NewKeepalive()
	go ka.Start()
	go keepalive(ka, cfg)
	go run(ka, cfg)

	// Wait forever.
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
