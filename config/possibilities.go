package config

import (
	"bufio"
	"embed"
	"encoding/json"
	"log"
	"os"
)

// The text file is embedded at compile time.
// https://pkg.go.dev/embed#FS.ReadFile

//go:embed searches.json
var f embed.FS

type Search struct {
	Field string `json:"field"`
	Query string `json:"query"`
}

func GetSearches() []Search {

	data, _ := f.ReadFile("searches.json")
	searches := make([]Search, 0)
	err := json.Unmarshal(data, &searches)
	if err != nil {
		log.Fatal("could not unmarshal search strings from embedded data.")
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
