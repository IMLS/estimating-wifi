package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/jszwec/csvutil"
	"gsa.gov/18f/version"
)

// External: https://nicedoc.io/jszwec/csvutil#user-content-marshal-a-nameexamples_marshala

type config struct {
	Key         string
	Sqlite      bool
	CSV         bool
	Fcfs_seq_id string
	Device_tag  string
	GraphQL     string
	Events      string
	Wifi        string
	Stepsize    int
	TzOffset    int
}

/* {"id":1382983,
"event_id":22233,
"fcfs_seq_id":"KY0069-003",
"device_tag":"berea1",
"localtime":"2021-05-18T18:58:08Z",
"servertime":"2021-05-18T18:58:08Z",
"session_id":"1e666a9ebe6e3a95",
"manufacturer_index":24,
"patron_index":6467
*/

type WifiEvents struct {
	Data []WifiEvent `json:"data"`
}

type WifiEvent struct {
	ID                int       `json:"id"`
	EventId           int       `json:"event_id"`
	FCFSSeqId         string    `json:"fcfs_seq_id"`
	DeviceTag         string    `json:"device_tag"`
	Localtime         time.Time `json:"localtime"`
	Servertime        time.Time `json:"servertime"`
	SessionId         string    `json:"session_id"`
	ManufacturerIndex int       `json:"manufacturer_index"`
	PatronIndex       int       `json:"patron_index"`
}

func getWifiEvents(cfg *config) ([]WifiEvent, error) {
	fetching := true
	events := make([]WifiEvent, 0)

	offset := 0
	for fetching {
		client := &http.Client{}
		req, err := http.NewRequest("GET", cfg.Wifi, nil)
		if err != nil {
			log.Println(err)
			log.Fatal("Could not create HTTP request.")
		}
		// Add the API key to the header.
		req.Header.Add("X-Api-Key", cfg.Key)
		q := req.URL.Query()
		q.Add("limit", fmt.Sprint(cfg.Stepsize))
		q.Add("offset", fmt.Sprint(offset))
		q.Add("filter[fcfs_seq_id][_eq]", cfg.Fcfs_seq_id)
		q.Add("filter[device_tag][_eq]", cfg.Device_tag)
		req.URL.RawQuery = q.Encode()
		fmt.Printf("URL: %v\n", req.URL.String())

		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			log.Fatal("Failure in HTTP client execution.")
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		e := new(WifiEvents)
		json.Unmarshal(body, &e)
		events = append(events, e.Data...)
		if len(e.Data) < cfg.Stepsize {
			fetching = false
		} else {
			offset += cfg.Stepsize
		}
	}
	return events, nil
}

func getKeys(d map[string]interface{}) []string {
	keys := make([]string, 0)
	for k := range d {
		keys = append(keys, k)
	}
	return keys
}

func getSessions(events []WifiEvent) []string {
	eventset := make(map[string]interface{})
	for _, e := range events {
		eventset[e.SessionId] = true
	}
	return getKeys(eventset)
}

func remapEvents(events []WifiEvent) []WifiEvent {
	// Get the unique sessions in the dataset
	sessions := getSessions(events)
	log.Println(sessions)
	// We need things in order. This matters for remapping
	// the manufacturer and patron indicies.
	sort.Slice(events, func(i, j int) bool {
		// return events[i].Localtime.Before(events[j].Localtime)
		return events[i].ID < events[j].ID
	})

	for _, s := range sessions {
		// For each session, create a new patron/mfg mapping.
		manufacturerMap := make(map[int]int)
		manufacturerNdx := 0
		patronMap := make(map[int]int)
		patronNdx := 0

		// Now, remap every event.
		// This means putting it in a new session (based on days instead of event ids)
		// and renumbering the patron/
		for ndx, e := range events {
			if e.SessionId == s {
				// We will rewrite all session fields to the current day.
				// Need to modify the array, not the local variable `e`
				events[ndx].SessionId = fmt.Sprint(e.Localtime.Format("2006-01-02"))

				// If we have already mapped this manufacturer, then update
				// the event with the stored value.
				if val, ok := manufacturerMap[e.ManufacturerIndex]; ok {
					events[ndx].ManufacturerIndex = val
				} else {
					// Otherwise, update the map and the event.
					manufacturerMap[e.ManufacturerIndex] = manufacturerNdx
					events[ndx].ManufacturerIndex = manufacturerNdx
					manufacturerNdx += 1
				}

				// Now, check the patron index.
				if val, ok := patronMap[e.PatronIndex]; ok {
					events[ndx].PatronIndex = val
				} else {
					// fmt.Printf("[%v] Remapping %v to %v\n", s, e.PatronIndex, patronNdx)
					patronMap[e.PatronIndex] = patronNdx
					events[ndx].PatronIndex = patronNdx
					patronNdx += 1
				}
			}
		}
	}

	// Now, all of the events have been remapped.
	return events
}

func getEventsFromSession(events []WifiEvent, session string) []WifiEvent {
	filtered := make([]WifiEvent, 0)
	for _, e := range events {
		if e.SessionId == session {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

func fixLocaltime(cfg *config, events []WifiEvent) []WifiEvent {
	updated := make([]WifiEvent, 0)
	for _, e := range events {
		e.Localtime = e.Localtime.Add(time.Duration(cfg.TzOffset) * time.Hour)
		updated = append(updated, e)
	}
	return updated
}

func generateCSV(cfg *config) {
	allEvents, err := getWifiEvents(cfg)
	if err != nil {
		log.Println("no events retrieved.")
	}
	fixed := fixLocaltime(cfg, allEvents)
	remapped := remapEvents(fixed)

	for _, s := range getSessions(remapped) {
		events := getEventsFromSession(remapped, s)
		b, err := csvutil.Marshal(events)
		if err != nil {
			log.Println(err)
			log.Fatal("could not convert events to CSV")
		}

		_ = os.Mkdir("output", 0777)
		fcfs_tag := fmt.Sprintf("%v-%v", cfg.Fcfs_seq_id, cfg.Device_tag)
		outdir := filepath.Join("output", fcfs_tag)
		_ = os.Mkdir(outdir, 0777)
		f, err := os.Create(filepath.Join(outdir, fmt.Sprintf("%v-%v.csv", s, fcfs_tag)))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		f.Write(b)
	}

}

func main() {
	versionPtr := flag.Bool("version", false, "Get the software version and exit.")
	fcfsSeqIdPtr := flag.String("fcfs_seq_id", "", "Set the FCFS Seq Id to process.")
	deviceTagPtr := flag.String("device_tag", "", "Set the device tag to process.")
	sqlitePtr := flag.Bool("sqlite", false, "Generate an SQLite table of the data.")
	csvPtr := flag.Bool("csv", true, "Generate a CSV file of the data. Default.")
	graphQLPtr := flag.String("graphql", "https://api.data.gov/TEST/10x-imls/v1/graphql/", "GraphQL endpoint.")
	eventsPtr := flag.String("events", "https://api.data.gov/TEST/10x-imls/v1/search/events/", "Events REST endpoint.")
	wifiPtr := flag.String("wifi", "https://api.data.gov/TEST/10x-imls/v1/search/wifi/", "Wifi REST endpoint.")
	stepSizePtr := flag.Int("stepsize", 10000, "How many elements to retrieve per HTTPS GET.")
	// The server is recording time in Z, or Zulu time, which is GMT.
	tzOffsetPtr := flag.Int("tzoffset", -5, "localtime timezone offset from server")

	flag.Parse()

	// VERSION
	// If they just want the version, print it and exit.
	if *versionPtr {
		fmt.Println(version.GetVersion())
		os.Exit(0)
	}

	// Make sure we have an API key to work with.
	if os.Getenv("APIDATAGOV") == "" {
		fmt.Println("Please set APIDATAGOV in the environment before running.")
		os.Exit(-1)
	}

	if *fcfsSeqIdPtr == "" || *deviceTagPtr == "" {
		fmt.Println("Please set both fcfs_seq_id and device_tag.")
		os.Exit(-1)
	}

	cfg := config{Key: os.Getenv("APIDATAGOV"),
		Sqlite:      *sqlitePtr,
		CSV:         *csvPtr,
		Fcfs_seq_id: *fcfsSeqIdPtr,
		Device_tag:  *deviceTagPtr,
		GraphQL:     *graphQLPtr,
		Events:      *eventsPtr,
		Wifi:        *wifiPtr,
		Stepsize:    *stepSizePtr,
		TzOffset:    *tzOffsetPtr,
	}

	generateCSV(&cfg)
}
