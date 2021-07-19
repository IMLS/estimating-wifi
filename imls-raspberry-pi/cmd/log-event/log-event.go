package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"gsa.gov/18f/internal/logwrapper"
	config "gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/version"
)

func main() {
	versionPtr := flag.Bool("version", false, "Get the software version and exit.")
	configPathPtr := flag.String("config", "", "Path to config.yaml. REQUIRED.")
	flag.Parse()
	rest := flag.Args()

	// If they just want the version, print and exit.
	if *versionPtr {
		fmt.Println(version.GetVersion())
		os.Exit(0)
	}

	if *configPathPtr == "" {
		log.Println("The flag -config MUST be provided.")
		os.Exit(1)
	}

	config.NewConfigFromPath(*configPathPtr)
	config.InitConfig()

	lw := logwrapper.NewLogger(nil)
	lw.Info(strings.Join(rest, " "))
}
