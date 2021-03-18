package main

import (
	"crypto/sha256"
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

/* PROCESS tockEveryN
 * consumes a tag (for logging purposes) as well as
 * a driving `tick` on `in`. Every `n` ticks, it outputs
 * a boolean `tock` on the channel `out`.
 * When `in` is every second, and `n` is 60, it turns
 * a stream of second ticks into minute `tocks`.
 */
func tockEveryN(ka *csp.Keepalive, n int, in <-chan bool, out chan<- bool) {
	log.Println("Starting tockEveryN")
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

/* PROCESS runWireshark
 * Runs a subprocess for a duration of OBSERVE_SECONDS.
 * Therefore, this process effectively blocks for that time.
 * Gathers a hashmap of [MAC -> count] values. This hashmap
 * is then communicated out.
 * Empty MAC addresses are filtered out.
 */
func runWireshark(ka *csp.Keepalive, cfg *model.Config, in <-chan bool, out chan<- map[string]int) {
	log.Println("Starting runWireshark")
	// If we have to wait twice the monitor duration, something broke.
	ping, pong := ka.Subscribe("runWireshark", cfg.Wireshark.Duration*2)

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
			// for `k, _ :=` is the same as `for k :=`
			for k := range macmap {
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

/* PROC mac_to_mfg
 * Takes in a hashmap of MAC addresses and counts, and passes on a hashmap
 * of manufacturer IDs and counts.
 * Uses "unknown" for all unknown manufacturers.
 */
func macToEntry(ka *csp.Keepalive, cfg *model.Config, macmap <-chan map[string]int, mfgmap chan<- map[string]model.Entry) {
	log.Println("Starting macToEntry")
	ping, pong := ka.Subscribe("macToEntry", 5)

	for {
		select {
		case <-ping:
			pong <- "macToEntry"

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

/* PROC reportMap
 * Takes a hashmap of [mfg id : count] and POSTs
 * each one to the server individually. We have no bulk insert.
 */
func reportMap(ka *csp.Keepalive, cfg *model.Config, mfgs <-chan map[string]model.Entry) {
	log.Println("Starting reportMap")
	ping, pong := ka.Subscribe("reportMap", 5)

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
				pong <- "reportMap"
			} else {
				log.Printf("report: http_error_count threshold of %d reached\n", http_error_count)
			}

		case m := <-mfgs:
			count = count + 1
			log.Println("reporting: ", count)
			// Try and grab the token from the OS Env.
			// It would have been set if we found it in a global config file.
			accessKey := os.Getenv(constants.AuthTokenKey)
			tok := &model.Token{}
			if len(accessKey) > 1 {
				tok.Data.AccessToken = accessKey
			} else {
				// If the token is too short/empty, we should try and get a token
				// via username/password in the env. This should have failed long ago
				// if the username/password are not in the env.
				apiTok, err := api.Get_token(cfg)
				tok = apiTok
				if err != nil {
					log.Println("report: error in token fetch")
					log.Println(err)
					http_error_count = http_error_count + 1
				}
			}

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

		case <-time.After(time.Duration(cfg.Monitoring.HTTPErrorIntervalMins) * time.Minute):
			// If this much time has gone by, go ahead and reset the error count.
			if http_error_count != 0 {
				log.Printf("report: resetting http_error_count from %d to 0\n", http_error_count)
				http_error_count = 0
			}
		}
	}
}

func ringBuffer(ka *csp.Keepalive, cfg *model.Config, in <-chan map[string]int, out chan<- map[string]int) {
	log.Println("Starting ringBuffer")
	ping, pong := ka.Subscribe("ringBuffer", 3)

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
			pong <- "ringBuffer"

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
					for mac := range m {
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

func rawToUids(ka *csp.Keepalive, cfg *model.Config, in <-chan map[string]int, out chan<- map[model.UserMapping]int, kill <-chan bool) {
	log.Println("Starting rawToUids")

	// This is an odd construct to facilitate unit testing.
	var ping chan interface{} = nil
	var pong chan interface{} = nil
	if kill == nil {
		ping, pong = ka.Subscribe("rawToUids", 5)
	}
	log.Println("rtu: initialized keepalive")

	macToNdx := make(map[string]int)
	ndxToMac := make(map[int]string)

	// The ndx, or nextId, needs to be maintained as a "monotonically increasing"
	// value for the life of a session-counter run.
	nextId := 0
	// To track who has overstayed their disconnection window.
	uniq := make(map[model.UserMapping]int)
	disco := make(map[model.UserMapping]int)

	for {
		select {
		case <-kill:
			log.Println("rtu: exiting")
			return
		case <-ping:
			pong <- "rawToUids"
		case m := <-in:
			log.Println("rtu: received map: ", m)
			// For each incoming address, decide if it is already in our map.
			// If it is, do nothing. If not, give that mac address a new id.
			for addr := range m {
				_, found := macToNdx[addr]
				log.Printf("rtu: [%v :: %v]\n", addr, found)
				if !found {
					log.Printf("rtu: adding [%v] as [%v]\n", addr, nextId)
					macToNdx[addr] = nextId
					nextId += 1
				}
			}
			// Now, build a new mapping for sending down the pipeline.
			// That mapping will be their user id and the device manufacturer
			// and we will keep the "count" of the number of packets that WS saw.
			// That number is probably not meaningful, but we'll hold it for a moment.
			newMapping := make(map[model.UserMapping]int)

			for oldaddr, v := range m {
				mfg := api.Mac_to_mfg(cfg, oldaddr)
				um := model.UserMapping{Mfg: mfg, Id: macToNdx[oldaddr]}
				log.Println("rtu: newmap ", um, " to ", v)
				newMapping[um] = v
				// If you just arrived to be mapped, you by
				// definition have a 0 uniqueness window ticker
				// and a 0 disco ticker.
				uniq[um] = 0
				disco[um] = 0
			}

			// Everyone we do *not* see has their time bumped.
			// Everyone we see has their uniqueness timeout set to 0.
			// And, their disconnect timeout must necessarily be reset as well.
			log.Println("newmap is ", newMapping)
			for k := range uniq {
				log.Println("looking for ", k, " in uniq")
				_, here := newMapping[k]
				if !here {
					uniq[k] = uniq[k] + 1
					disco[k] = disco[k] + 1
				}
			}

			// Now, we have to do some munging. A new map is needed.
			sendmap := make(map[model.UserMapping]int)
			for k, v := range uniq {
				sendmap[k] = v
			}

			// Anyone who timed out should not be communicated in the sendmap.
			for k, v := range disco {
				if v > cfg.Monitoring.DisconnectionWindow {
					delete(sendmap, k)
				}
			}

			// And, if anyone overstays their uniqueness,
			// complettely remove them.
			for k, v := range uniq {
				if v > cfg.Monitoring.UniquenessWindow {
					delete(sendmap, k)
					delete(disco, k)
					delete(uniq, k)
					// And, clean up the other state we have lying around
					addr := ndxToMac[k.Id]
					delete(ndxToMac, k.Id)
					delete(macToNdx, addr)
				}
			}

			out <- sendmap

		}
	}
}

/* FUNC checkEnvVars
 * Checks to see if the username and password for
 * working with Directus is in memory.
 * If not, it quits.
 */
func checkEnvVars() {
	if os.Getenv(constants.EnvUsername) == "" {
		fmt.Printf("%s must be set in the env!\n", constants.EnvUsername)
		os.Exit(constants.ExitNoUsername)
	}
	if os.Getenv(constants.EnvPassword) == "" {
		fmt.Printf("%s must be set in the env!\n", constants.EnvPassword)
		os.Exit(constants.ExitNoPassword)
	}
}

func parseConfigFile(filepath string) (*model.Config, error) {
	_, err := os.Stat(filepath)

	// Stat will set an error if the file cannot be found.
	if err == nil {
		f, err := os.Open(filepath)
		if err != nil {
			log.Fatal("parseConfigFile: could not open configuration file. Exiting.")
		}
		defer f.Close()
		var cfg *model.Config
		decoder := yaml.NewDecoder(f)
		err = decoder.Decode(&cfg)
		if err != nil {
			log.Fatalf("parseConfigFile: could not decode YAML:\n%v\n", err)
		}
		return cfg, nil
	} else {
		log.Printf("parseConfigFile: could not find config: %v\n", filepath)
	}
	return nil, fmt.Errorf("config: could not find config file [%v]", filepath)
}
func devConfig() *model.Config {
	checkEnvVars()
	// FIXME consider turning this into an env var
	cfgPtr := flag.String("config", "config.yaml", "config file")
	flag.Parse()
	cfg, err := parseConfigFile(*cfgPtr)
	if err != nil {
		log.Println("config: could not load dev config. Exiting.")
		log.Fatalln(err)
	}
	return cfg
}

func readAuth() (*model.Auth, error) {
	_, err := os.Stat(constants.AuthPath)
	if err != nil {
		return &model.Auth{}, fmt.Errorf("readToken: cannot find default token file at [%v]", constants.AuthPath)
	}

	f, err := os.Open(constants.AuthPath)
	if err != nil {
		log.Fatal("readToken: could not open token file. Exiting.")
	}
	defer f.Close()
	var auth *model.Auth
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&auth)
	if err != nil {
		log.Fatalf("readToken: could not decode YAML:\n%v\n", err)
	}

	return auth, nil
}

func readConfig() *model.Config {
	// We expect config to be here:
	//   * /etc/session-counter/config.yaml
	// We expect there to be a token file at
	//   * /etc/session-counter/access-token
	//
	// If neither of those is true, we can check for a
	// the username and password to be in the ENV, and
	// for the config to be passed via command line.

	cfg, err := parseConfigFile(constants.ConfigPath)
	if err != nil {
		fmt.Printf("config: could not find config at default path [%v]\n", constants.ConfigPath)
		fmt.Println("config: loading dev config")
		return devConfig()
	}

	auth, err := readAuth()
	if err != nil {
		log.Fatal("readConfig: cannot find auth token")
	}

	os.Setenv(constants.AuthTokenKey, auth.Token)
	os.Setenv(constants.AuthEmailKey, auth.Email)
	return cfg
}

func run(ka *csp.Keepalive, cfg *model.Config) {
	log.Println("Starting run")
	// Create channels for process network
	ch_sec := make(chan bool)
	ch_nsec := make(chan bool)
	ch_macs := make(chan map[string]int)
	ch_macs_counted := make(chan map[string]int)
	ch_mfg := make(chan map[string]model.Entry)
	ch_uniq := make(chan map[model.UserMapping]int)

	// Run the process network.
	// Driven by a 1s `tick` process.
	// Thread the keepalive through the network
	go tick(ka, ch_sec)
	go tockEveryN(ka, 60, ch_sec, ch_nsec)
	go runWireshark(ka, cfg, ch_nsec, ch_macs)
	go ringBuffer(ka, cfg, ch_macs, ch_macs_counted)
	// ka *csp.Keepalive, cfg *model.Config, in <-chan map[string]int, remove <-chan int, out chan<- map[model.UserMapping]int
	go rawToUids(ka, cfg, ch_macs, ch_uniq, nil)
	go macToEntry(ka, cfg, ch_macs_counted, ch_mfg)
	go reportMap(ka, cfg, ch_mfg)
}

func keepalive(ka *csp.Keepalive, cfg *model.Config) {
	log.Println("Starting keepalive")
	var counter int64 = 0
	for {
		time.Sleep(time.Duration(cfg.Monitoring.PingInterval) * time.Second)
		ka.Publish(counter)
		counter = counter + 1
	}
}

func calcSessionId() string {
	h := sha256.New()
	email := os.Getenv(constants.AuthEmailKey)
	// FIXME: Use the email instead of the token.
	// Guaranteed to be unique. Current time along with our auth token, hashed.
	h.Write([]byte(fmt.Sprintf("%v%x", time.Now(), email)))
	sid := fmt.Sprintf("%x", h.Sum(nil))
	log.Println("Session id: ", sid)
	return sid
}

func main() {
	// Read in a config
	cfg := readConfig()
	// Add a "sessionId" to the mix.
	cfg.SessionId = calcSessionId()

	ka := csp.NewKeepalive()
	go ka.Start()
	go keepalive(ka, cfg)
	go run(ka, cfg)

	// Wait forever.
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
