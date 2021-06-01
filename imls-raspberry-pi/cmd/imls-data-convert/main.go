package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/briandowns/spinner"
	_ "github.com/mattn/go-sqlite3"
	"gsa.gov/18f/analysis"

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

func spinnerStart(msg string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = msg
	s.Start()
	return s
}

func getWifiEvents(cfg *config) ([]analysis.WifiEvent, error) {
	fetching := true
	events := make([]analysis.WifiEvent, 0)
	offset := 0

	s := spinnerStart(" Fetching wifi events...")
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
		// fmt.Printf("URL: %v\n", req.URL.String())

		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			log.Fatal("Failure in HTTP client execution.")
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		e := new(analysis.WifiEvents)
		json.Unmarshal(body, &e)
		events = append(events, e.Data...)
		if len(e.Data) < cfg.Stepsize {
			fetching = false
		} else {
			offset += cfg.Stepsize
		}
	}
	s.Stop()
	return events, nil
}

func fixLocaltime(cfg *config, events []analysis.WifiEvent) []analysis.WifiEvent {
	updated := make([]analysis.WifiEvent, 0)
	for _, e := range events {
		e.Localtime = e.Localtime.Add(time.Duration(cfg.TzOffset) * time.Hour)
		updated = append(updated, e)
	}
	return updated
}

func fixEvents(cfg *config) []analysis.WifiEvent {
	allEvents, err := getWifiEvents(cfg)
	if err != nil {
		log.Println("no events retrieved.")
	}
	fixed := fixLocaltime(cfg, allEvents)
	remapped := analysis.RemapEvents(fixed)
	return remapped
}

func generateCSV(cfg *config, remapped []analysis.WifiEvent) {
	for _, s := range analysis.GetSessions(remapped) {
		events := analysis.GetEventsFromSession(remapped, s)
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

func generateSqlite(cfg *config, remapped []analysis.WifiEvent) {
	_ = os.Mkdir("output", 0777)
	fcfs_tag := fmt.Sprintf("%v-%v", cfg.Fcfs_seq_id, cfg.Device_tag)
	outdir := filepath.Join("output", fcfs_tag)
	_ = os.Mkdir(outdir, 0777)
	db, err := sql.Open("sqlite3", string(filepath.Join(outdir, fmt.Sprintf("%v.sqlite", fcfs_tag))))
	if err != nil {
		log.Fatal("could not open SQLite3 DB for writing.")
	}
	defer db.Close()
	//ID,EventId,FCFSSeqId,DeviceTag,Localtime,Servertime,SessionId,ManufacturerIndex,PatronIndex
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS wifi 
			(id INTEGER PRIMARY KEY, 
			event_id INTEGER, 
			fcfs_seq_id TEXT,
			device_tag TEXT,
			localtime TEXT,
			servertime TEXT,
			session_id TEXT,
			manufacturer_index INTEGER,
			patron_index INTEGER)`)
	if err != nil {
		log.Println("error creating table")
		log.Fatal(err)
	}

	stat, _ := db.Prepare(`INSERT INTO wifi 
			(event_id, fcfs_seq_id, device_tag, localtime, servertime, session_id, manufacturer_index, patron_index) 
			VALUES  (?, ?, ?, ?, ?, ?, ?, ?)`)
	s := spinnerStart(" Writing data to SQLite table...")
	for ndx, e := range remapped {
		if ndx%1000 == 0 {
			time.Sleep(10 * time.Millisecond)
		}
		stat.Exec(e.EventId, e.FCFSSeqId, e.DeviceTag, e.Localtime.Format(time.RFC3339), e.Servertime.Format(time.RFC3339), e.SessionId, e.ManufacturerIndex, e.PatronIndex)
	}
	defer stat.Close()
	s.Stop()
}

func getLibraries(cfg *config) map[string][]string {
	s := spinnerStart(" Fetching event events...")
	set := make(map[string][]string)
	// Fetch the last 50K events;
	for count := 0; count < 10; count++ {
		client := &http.Client{}
		req, err := http.NewRequest("GET", cfg.Events, nil)
		if err != nil {
			log.Println(err)
			log.Fatal("Could not create HTTP request.")
		}
		// Add the API key to the header.
		req.Header.Add("X-Api-Key", cfg.Key)
		q := req.URL.Query()
		q.Add("limit", fmt.Sprint(cfg.Stepsize))
		q.Add("offset", fmt.Sprint(cfg.Stepsize*count))
		q.Add("fields", "tag,fcfs_seq_id,device_tag")
		req.URL.RawQuery = q.Encode()
		//log.Printf("URL: %v\n", req.URL.String())

		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			log.Fatal("Failure in HTTP client execution.")
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		evts := new(analysis.EventEvents)
		json.Unmarshal(body, &evts)

		for _, e := range evts.Data {
			set[fmt.Sprint(e.FCFSSeqId, e.DeviceTag)] = []string{e.FCFSSeqId, e.DeviceTag}
		}
	}
	s.Stop()

	return set
}

func dedupe(events []analysis.WifiEvent) []analysis.WifiEvent {
	clean := make([]analysis.WifiEvent, 0)
	for _, e := range events {
		found := false
		for _, c := range clean {
			if e.ID == c.ID {
				found = true
			}
		}
		if !found {
			clean = append(clean, e)
		}
	}
	return clean
}

func main() {
	versionPtr := flag.Bool("version", false, "Get the software version and exit.")
	getLibrariesPtr := flag.Bool("get-libraries", false, "Fetch a list of libraries in the dataset and exit.")
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

	if !*getLibrariesPtr && (*fcfsSeqIdPtr == "" || *deviceTagPtr == "") {
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

	if *getLibrariesPtr {
		libs := getLibraries(&cfg)
		fmt.Println("fcfs_seq_id,device_tag")

		for k, v := range libs {
			ismatch, err := regexp.Match(`[A-Z]{2}[0-9]{4}-[0-9]{3}`, []byte(k))
			if err == nil && ismatch {
				fmt.Printf("%v,%v\n", v[0], v[1])
			}
		}
		os.Exit(0)
	}

	remapped := fixEvents(&cfg)
	deduped := dedupe(remapped)
	generateCSV(&cfg, deduped)
	if *sqlitePtr {
		generateSqlite(&cfg, remapped)
	}
}
