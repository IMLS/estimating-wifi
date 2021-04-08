package config

import (
	"bufio"
	"embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gsa.gov/18f/find-ralink/constants"
)

// The text file is embedded at compile time.
// https://pkg.go.dev/embed#FS.ReadFile

//go:embed searches.json
var f embed.FS
var Verbose bool = false

type Search struct {
	Field string `json:"field"`
	Query string `json:"query"`
}

func GetSearches() []Search {
	searches := make([]Search, 0)

	// First, look for the file in /etc
	if _, err := os.Stat(constants.SEARCHES_PATH); err == nil {
		if Verbose {
			fmt.Println("using", constants.SEARCHES_PATH)
		}
		// We found the version in /etc. This is probably a live installation.
		file, err := os.Open(constants.SEARCHES_PATH)
		if err != nil {
			log.Fatal("error opening /etc/session-counter/searches.json")
		}
		defer file.Close()
		data, _ := ioutil.ReadAll(file)
		err = json.Unmarshal(data, &searches)
		if err != nil {
			log.Fatal("could not unmarshal search strings from embedded data.")
		}

	} else if os.IsNotExist(err) {
		if Verbose {
			fmt.Println("using embedded search data")
		}
		// Use the embedded file, which has a limited set of search terms.
		data, _ := f.ReadFile("searches.json")
		err := json.Unmarshal(data, &searches)
		if err != nil {
			log.Fatal("could not unmarshal search strings from embedded data.")
		}
	} else {
		log.Fatal("somehow could not find any search strings. exiting.")
	}

	return searches
}

// https://stackoverflow.com/questions/5884154/read-text-file-into-string-array-and-write
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)
		lines = append(lines, line)
	}
	return lines, scanner.Err()
}
