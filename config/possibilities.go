package config

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Saved searches are expressed as JSON.
// We can choose a field to search (or "ALL")
// and a regex as our query.
type Search struct {
	Field string `json:"field"`
	Query string `json:"query"`
}

// The text file is embedded at compile time.
// https://pkg.go.dev/embed#FS.ReadFile
//go:embed searches.json
var f embed.FS
var Verbose *bool = new(bool)

// PURPOSE
// GetSearches attempts to read in the JSON document from the filesystem
// and use that, or it attempts to use the embedded version. The embedded version
// is used as a fallback in the case that we cannot find a (presumably tweaked/custom)
// set of searches in /etc...
func GetSearches() []Search {
	searches := make([]Search, 0)

	// First, look for the file in /etc
	if _, err := os.Stat(SEARCHES_PATH); err == nil {
		if *Verbose {
			fmt.Println("using", SEARCHES_PATH)
		}
		// We found the version in /etc. This is probably a live installation.
		file, err := os.Open(SEARCHES_PATH)
		if err != nil {
			log.Fatalf("error opening [%v]", SEARCHES_PATH)
		}
		defer file.Close()

		data, _ := ioutil.ReadAll(file)
		// If we are successful here, we should end up with our array of
		// searches filled by the Unmarshal.
		err = json.Unmarshal(data, &searches)
		if err != nil {
			log.Fatal("could not unmarshal search strings from embedded data.")
		}
	} else if os.IsNotExist(err) {
		// If we're here, we're going to read the JSON document that is
		// embedded in the executable file.
		if *Verbose {
			fmt.Println("using embedded search data")
		}
		// Use the embedded file, which has a limited set of search terms.
		data, _ := f.ReadFile("searches.json")
		err := json.Unmarshal(data, &searches)
		if err != nil {
			log.Fatal("could not unmarshal search strings from embedded data.")
		}
	} else {
		// We really shouldn't end up here. Log it and die.
		log.Fatal("somehow could not find any search strings. exiting.")
	}

	return searches
}
