package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"gsa.gov/18f/config"
	"gsa.gov/18f/http"
	"gsa.gov/18f/version"
)

func isJsonOk(jsonString string) bool {
	var dat map[string]interface{}
	err := json.Unmarshal([]byte(jsonString), &dat)
	return err == nil
}

func readFileAsB64(path string) (string, error) {
	if _, err := os.Stat(path); err == nil {
		bs, err := os.ReadFile(path)
		if err != nil {
			return "NOREAD", err
		}
		sEnc := b64.StdEncoding.EncodeToString([]byte(bs))
		return sEnc, nil
	  } else if os.IsNotExist(err) {
		return "FILEDOESNOTEXIST", err
	  } else {
		// This is a rare condition. Leave the world.
		return "", err
	  }
}

func main() {
	versionPtr := flag.Bool("version", false, "Get the software version and exit.")
	tagPtr := flag.String("tag", "event", "Event tag.")
	infoPtr := flag.String("info", "{}", "Valid JSON for info. Exclusive w/ --file.")
	filePtr := flag.String("file", "", "Filepath to include as info. Exclusive w/ --info.")
	configPathPtr := flag.String("config", "", "Path to config.yaml. REQUIRED.")
	flag.Parse()

	// If they just want the version, print and exit.
	if *versionPtr {
		fmt.Println(version.GetVersion())
		os.Exit(0)
	}

	if *configPathPtr == "" {
		log.Println("The flag --config MUST be provided.")
		os.Exit(1)
	} else {
		config.SetConfigPath(*configPathPtr)
	}

	// Make sure we're exclusive between these two flags.
	if *filePtr != "" && *infoPtr != "{}" {
		fmt.Println("Only one of --info or --file permitted.")
		os.Exit(-1)
	}

	cfg := config.NewConfig()
	err := cfg.ReadConfig(*configPathPtr)
	if err != nil {
		log.Fatal("log-event: error loading config.")
	}

	// This reuses the *infoPtr for reporting below.
	// Saves us having more if-then statements. :shrug:
	// Not a nice reuse...
	if *filePtr != "" {
		b64, err := readFileAsB64(*filePtr)
		if err != nil {
			fmt.Println("Could not B64 convert file. Exiting.")
			os.Exit(-1)
		}
		*infoPtr = `{"file": "` + b64 + `"}`
	}


	if isJsonOk(*infoPtr) {
		logger := http.NewEventLogger(cfg)
		logger.LogJSON(*tagPtr, *infoPtr)
	} else {
		log.Fatal("BAD JSON: ", *infoPtr)
	}

}
