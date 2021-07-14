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
	"sync"
	"time"

	"github.com/briandowns/spinner"
	_ "github.com/mattn/go-sqlite3"
	"gsa.gov/18f/internal/analysis"

	"gsa.gov/18f/internal/version"
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
	DestDir     string
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

func getWifiEvents(cfg *config, events chan<- []analysis.WifiEvent, wg *sync.WaitGroup) {
	fetching := true
	// events := make([]analysis.WifiEvent, 0)
	offset := 0

	log.Println("Fetching wifi events...")
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

		// events = append(events, e.Data...)
		events <- e.Data

		if len(e.Data) < cfg.Stepsize {
			fetching = false
		} else {
			offset += cfg.Stepsize
			log.Println(fmt.Sprintf(" fetched %v events", offset))
		}
	}

	//return events, nil

	// Leave the gofuncs
	log.Println("Done reading events")
	events <- nil
	wg.Done()
}

// func getEvents(cfg *config) []analysis.WifiEvent {
// 	allEvents, err := getWifiEvents(cfg)
// 	if err != nil {
// 		log.Println("no events retrieved.")
// 	}
// 	return allEvents
// }

func generateSqlite(cfg *config, ch <-chan []analysis.WifiEvent, wg *sync.WaitGroup) {
	// _ = os.Mkdir("output", 0777)
	fcfs_tag := fmt.Sprintf("%v-%v", cfg.Fcfs_seq_id, cfg.Device_tag)
	db, err := sql.Open("sqlite3", string(filepath.Join(cfg.DestDir, fmt.Sprintf("%v.sqlite", fcfs_tag))))
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

	// stat, _ := db.Prepare(`INSERT INTO wifi
	// 		(event_id, fcfs_seq_id, device_tag, localtime, servertime, session_id, manufacturer_index, patron_index)
	// 		VALUES  (?, ?, ?, ?, ?, ?, ?, ?)`)
	// defer stat.Close()

	//. s := spinnerStart(" Writing data to SQLite table...")
	for {
		events := <-ch
		// Transactions are required to speed this up. Massively.
		// https://jmoiron.github.io/sqlx/
		// WIthout timing it, this runs around 2M events in a few miinutes.
		// But, it used to take *forever*.
		// Note this is storing the data in a DB on a RAMDISK. Speed will be slightly slower on an SSD.
		tx, _ := db.Begin()
		stat, _ := tx.Prepare(`INSERT INTO wifi 
			(event_id, fcfs_seq_id, device_tag, localtime, servertime, session_id, manufacturer_index, patron_index) 
			VALUES  (?, ?, ?, ?, ?, ?, ?, ?)`)
		defer stat.Close()

		if events == nil {
			stat.Close()
			db.Close()
			log.Println("Done writing events.")
			wg.Done()
			return
		} else {
			for _, e := range events {
				stat.Exec(e.EventId, e.FCFSSeqId, e.DeviceTag, e.Localtime.Format(time.RFC3339), e.Servertime.Format(time.RFC3339), e.SessionId, e.ManufacturerIndex, e.PatronIndex)
			}
			// This is the same as .Close()
			// Do all 10K inserts at once.
			err = tx.Commit()
			if err != nil {
				log.Fatal("could not execute transaction")
			}
		}
	}
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

func main() {
	versionPtr := flag.Bool("version", false, "Get the software version and exit.")
	getLibrariesPtr := flag.Bool("libraries", false, "Fetch a list of libraries in the dataset and exit.")
	fcfsSeqIdPtr := flag.String("fcfs_seq_id", "", "Set the FCFS Seq Id to process.")
	deviceTagPtr := flag.String("device_tag", "", "Set the device tag to process.")
	sqlitePtr := flag.Bool("sqlite", false, "Generate an SQLite table of the data.")
	graphQLPtr := flag.String("graphql", "https://api.data.gov/TEST/10x-imls/v1/graphql/", "GraphQL endpoint.")
	eventsPtr := flag.String("events", "https://api.data.gov/TEST/10x-imls/v1/search/events/", "Events REST endpoint.")
	wifiPtr := flag.String("wifi", "https://api.data.gov/TEST/10x-imls/v1/search/wifi/", "Wifi REST endpoint.")
	stepSizePtr := flag.Int("stepsize", 10000, "How many elements to retrieve per HTTPS GET. Default is 10K.")
	// The server is recording time in Z, or Zulu time, which is GMT.
	tzOffsetPtr := flag.Int("tzoffset", -5, "localtime timezone offset from server")
	destPtr := flag.String("dest", "", "Destination directory for sqlite db.")

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
		CSV:         false,
		Fcfs_seq_id: *fcfsSeqIdPtr,
		Device_tag:  *deviceTagPtr,
		GraphQL:     *graphQLPtr,
		Events:      *eventsPtr,
		Wifi:        *wifiPtr,
		Stepsize:    *stepSizePtr,
		TzOffset:    *tzOffsetPtr,
		DestDir:     *destPtr,
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

	var wg sync.WaitGroup
	wg.Add(2)
	ch := make(chan []analysis.WifiEvent)
	go getWifiEvents(&cfg, ch, &wg)
	go generateSqlite(&cfg, ch, &wg)
	wg.Wait()

}

// https://gist.github.com/htr3n/344f06ba2bb20b1056d7d5570fe7f596
// diskutil erasevolume HFS+ 'RAMDISK' `hdiutil attach -nobrowse -nomount ram://4194304
// go build ; ./cache-to-sqlite --fcfs_seq_id GA0027-004 --device_tag in-ops --dest /Volumes/RAMDISK
