package main

import (
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

func main() {
	versionPtr := flag.Bool("version", false, "Get the software version and exit.")
	tagPtr := flag.String("tag", "event", "Event tag.")
	infoPtr := flag.String("info", "{}", "Valid JSON for info.")
	configPathPtr := flag.String("config", "", "Path to config.yaml (for testing).")
	authPathPtr := flag.String("auth", "", "Path to auth.yaml (for testing).")
	flag.Parse()

	if *configPathPtr != "" {
		config.SetConfigPath(*configPathPtr)
	}

	if *authPathPtr != "" {
		config.SetAuthPath(*authPathPtr)
	}

	// If they just want the version, print and exit.
	if *versionPtr {
		fmt.Println(version.GetVersion())
		os.Exit(0)
	}

	if isJsonOk(*infoPtr) {
		logger := http.NewEventLogger(config.ReadConfig())
		logger.LogJSON(*tagPtr, *infoPtr)
	} else {
		log.Fatal("BAD JSON: ", *infoPtr)
	}

}
