package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	_ "github.com/mattn/go-sqlite3"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"

	"gsa.gov/18f/internal/version"
)

// External: https://nicedoc.io/jszwec/csvutil#user-content-marshal-a-nameexamples_marshala

type config struct {
	Key       string
	Sqlite    bool
	CSV       bool
	FcfsSeqID string
	DeviceTag string
	GraphQL   string
	Durations string
	Stepsize  int
	TzOffset  int
	DestDir   string
}

func spinnerStart(msg string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Suffix = msg
	s.Start()
	return s
}

func getDurations(cfg *config, events chan<- []structs.Duration, wg *sync.WaitGroup) {
	fetching := true
	offset := 0

	log.Println("Fetching durations...")
	for fetching {
		client := &http.Client{}
		req, err := http.NewRequest("GET", cfg.Durations, nil)
		if err != nil {
			log.Println(err)
			log.Fatal("Could not create HTTP request.")
		}
		// Add the API key to the header.
		req.Header.Add("X-Api-Key", cfg.Key)
		q := req.URL.Query()
		q.Add("limit", fmt.Sprint(cfg.Stepsize))
		q.Add("offset", fmt.Sprint(offset))
		q.Add("filter[fcfs_seq_id][_eq]", cfg.FcfsSeqID)
		q.Add("filter[device_tag][_eq]", cfg.DeviceTag)
		req.URL.RawQuery = q.Encode()
		// fmt.Printf("URL: %v\n", req.URL.String())

		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			log.Fatal("Failure in HTTP client execution.")
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		e := new(structs.Durations)
		json.Unmarshal(body, &e)
		events <- e.Data

		if len(e.Data) < cfg.Stepsize {
			fetching = false
		} else {
			offset += cfg.Stepsize
			log.Println(fmt.Sprintf(" fetched %v durations", offset))
		}
	}

	events <- nil
	wg.Done()
}

func generateSqlite(cfg *config, ch <-chan []structs.Duration, wg *sync.WaitGroup) {
	fcfsTag := fmt.Sprintf("%v-%v", cfg.FcfsSeqID, cfg.DeviceTag)
	db := state.NewSqliteDB(fcfsTag + ".sqlite")
	db.CreateTableFromStruct(structs.Duration{})

	total := 0
	for {
		events := <-ch
		if events == nil {
			db.Close()
			log.Println("Done; wrote", total, "durations.")
			wg.Done()
			return
		} else {
			durations := make([]interface{}, 0)
			total += len(events)
			for _, e := range events {
				durations = append(durations, e)
			}
			s := spinnerStart(fmt.Sprint(" inserting ", len(durations), " durations"))
			s.Start()
			db.GetTableFromStruct(structs.Duration{}).InsertMany(durations)
			s.Stop()
		}
	}
}

func getLibraries(cfg *config) map[string][]string {
	s := spinnerStart(" Fetching event events...")
	set := make(map[string][]string)
	// Fetch the last 50K events;
	for count := 0; count < 10; count++ {
		client := &http.Client{}
		req, err := http.NewRequest("GET", cfg.Durations, nil)
		if err != nil {
			log.Println(err)
			log.Fatal("Could not create HTTP request.")
		}
		// Add the API key to the header.
		req.Header.Add("X-Api-Key", cfg.Key)
		q := req.URL.Query()
		q.Add("limit", fmt.Sprint(cfg.Stepsize))
		q.Add("offset", fmt.Sprint(cfg.Stepsize*count))
		q.Add("fields", "fcfs_seq_id,device_tag")
		q.Add("sort", "-session_id")
		req.URL.RawQuery = q.Encode()
		log.Printf("URL: %v\n", req.URL.String())

		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			log.Fatal("Failure in HTTP client execution.")
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		evts := new(structs.Durations)
		json.Unmarshal(body, &evts)

		for _, e := range evts.Data {
			set[fmt.Sprint(e.FCFSSeqID, e.DeviceTag)] = []string{e.FCFSSeqID, e.DeviceTag}
		}
	}
	s.Stop()

	return set
}

func main() {
	versionPtr := flag.Bool("version", false, "Get the software version and exit.")
	getLibrariesPtr := flag.Bool("libraries", false, "Fetch a list of libraries in the dataset and exit.")
	fcfsSeqIDPtr := flag.String("fcfs_seq_id", "", "Set the FCFS Seq Id to process.")
	deviceTagPtr := flag.String("device_tag", "", "Set the device tag to process.")
	sqlitePtr := flag.Bool("sqlite", false, "Generate an SQLite table of the data.")
	graphQLPtr := flag.String("graphql", "https://api.data.gov/TEST/10x-imls/v1/graphql/", "GraphQL endpoint.")
	// logsPtr := flag.String("events", "https://api.data.gov/TEST/10x-imls/v2/search/events/", "Events REST endpoint.")
	durationsPtr := flag.String("durations", "https://api.data.gov/TEST/10x-imls/v2/search/durations/", "Durations REST endpoint.")
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

	if !*getLibrariesPtr && (*fcfsSeqIDPtr == "" || *deviceTagPtr == "") {
		fmt.Println("Please set both fcfs_seq_id and device_tag.")
		os.Exit(-1)
	}

	cfg := config{Key: os.Getenv("APIDATAGOV"),
		Sqlite:    *sqlitePtr,
		CSV:       false,
		FcfsSeqID: *fcfsSeqIDPtr,
		DeviceTag: *deviceTagPtr,
		GraphQL:   *graphQLPtr,
		Durations: *durationsPtr,
		Stepsize:  *stepSizePtr,
		TzOffset:  *tzOffsetPtr,
		DestDir:   *destPtr,
	}

	if *getLibrariesPtr {
		libs := getLibraries(&cfg)
		fmt.Println("fcfs_seq_id,device_tag")
		re := regexp.MustCompile(`[A-Z]{2}[0-9]{4}-[0-9]{3}`)
		for k, v := range libs {
			if re.MatchString(k) {
				fmt.Printf("--fcfs_seq_id %v --device_tag %v\n", v[0], v[1])
			}
		}
		os.Exit(0)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	ch := make(chan []structs.Duration)
	go getDurations(&cfg, ch, &wg)
	go generateSqlite(&cfg, ch, &wg)
	wg.Wait()
}
