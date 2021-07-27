package tlp

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"syscall"

	"gsa.gov/18f/cmd/session-counter/constants"
	"gsa.gov/18f/internal/logwrapper"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/wifi-hardware-search/models"
	"gsa.gov/18f/internal/wifi-hardware-search/search"
)

func tshark(adapter string) []string {
	cfg := state.GetConfig()
	lw := logwrapper.NewLogger(nil)
	tsharkCmd := exec.Command(cfg.GetWiresharkPath(),
		"-a", fmt.Sprintf("duration:%d", cfg.GetWiresharkDuration()),
		"-I", "-i", adapter,
		"-Tfields", "-e", "wlan.sa")

	tsharkOut, err := tsharkCmd.StdoutPipe()
	if err != nil {
		lw.Error("could not open wireshark pipe")
		lw.Error(err.Error())
	}
	// The closer is called on exe exit. Idomatic use does not
	// explicitly call the closer.
	// defer tsharkOut.Close()

	err = tsharkCmd.Start()
	if err != nil {
		lw.Error("could not exe wireshark")
		lw.Error(err.Error())
	}
	tsharkBytes, err := ioutil.ReadAll(tsharkOut)
	if err != nil {
		lw.Info("did not read wireshark bytes")
		lw.Error(err.Error())
	}

	//tsharkCmd.Wait()
	// From https://stackoverflow.com/questions/10385551/get-exit-code-go
	if err := tsharkCmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				lw.Fatal("tshark exit status ", status.ExitStatus(), " ", string(tsharkBytes))
			}
		} else {
			lw.Fatal("tsharkCmd.Wait()", err.Error())
		}
	}

	macs := strings.Split(string(tsharkBytes), "\n")

	return macs
}

// RunWireshark runs a subprocess for a duration of OBSERVE_SECONDS. Therefore,
// this process effectively blocks for that time. Gathers a hashmap of [MAC ->
// count] values. This hashmap is then communicated out. Empty MAC addresses are
// filtered out.
func RunWireshark(ka *Keepalive, kb *KillBroker, in <-chan bool, out chan []string) {
	cfg := state.GetConfig()
	cfg.Log().Info("starting RunWireshark")
	var ping, pong chan interface{} = nil, nil
	var chKill chan interface{} = nil
	if kb != nil {
		chKill = kb.Subscribe()
	} else {
		ping, pong = ka.Subscribe("RunWireshark", 15)
	}

	// Adapter count... every "ac" ticks, we look up the adapter.
	// (ac % 0) guarantees that we look it up the first time.
	ticker := 0
	adapter := ""

	for {
		select {

		case <-ping:
			// We ping faster than this process can reply. However, we have a long
			// enough timeout that we will *eventually* catch up with all of the pings.
			pong <- "RunWireshark"

		case <-chKill:
			cfg.Log().Debug("exiting RunWireshark")
			return

		case <-in:
			// Look up the adapter. Use the find-ralink library.
			// The % will trigger first time through, which we want.
			var dev *models.Device = nil
			// If the config doesn't have this in it, we get a divide by zero.
			dev = search.SearchForMatchingDevice()

			// Only do a reading and continue the pipeline
			// if we find an adapter.
			if dev != nil && dev.Exists {
				// Load the config for use.
				// cfg.Wireshark.Adapter = dev.Logicalname
				cfg.Log().Debug("found adapter: ", dev.Logicalname)
				// Set monitor mode if the adapter changes.
				if dev.Logicalname != adapter {
					cfg.Log().Debug("setting monitor mode")
					search.SetMonitorMode(dev)
					adapter = dev.Logicalname
				}

				// This will block for [cfg.Wireshark.Duration] seconds.
				macmap := tshark(dev.Logicalname)
				// Mark and remove too-short MAC addresses
				// for removal from the tshark findings.
				var keepers []string
				// for `k, _ :=` is the same as `for k :=`
				for _, k := range macmap {
					if len(k) >= constants.MACLENGTH {
						keepers = append(keepers, k)
					}
				}
				// How many devices did we find to keep?
				// lw.Length("wireshark keepers", keepers)
				out <- keepers
			} else {
				cfg.Log().Info("no wifi devices found. no scanning carried out.")
				// Report an empty array of keepers
				out <- make([]string, 0)
			}

			// Bump our ticker
			ticker += 1
		}
	}
}
