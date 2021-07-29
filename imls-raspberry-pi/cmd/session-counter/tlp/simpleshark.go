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
)

func TSharkRunner(adapter string) []string {
	cfg := state.GetConfig()
	lw := logwrapper.NewLogger(nil)
	tsharkCmd := exec.Command(
		cfg.GetWiresharkPath(),
		"-a", fmt.Sprintf("duration:%d", cfg.GetWiresharkDuration()),
		"-I", "-i", adapter,
		"-Tfields", "-e", "wlan.sa")

	tsharkOut, err := tsharkCmd.StdoutPipe()
	if err != nil {
		lw.Error("could not open wireshark pipe")
		lw.Error(err.Error())
	}
	// The closer is called on exe exit. Idomatic use does not
	// explicitly call the closer. BUT DO WE HAVE LEAKS?
	defer tsharkOut.Close()

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

type SharkFn func(string) []string
type MonitorFn func(*models.Device)
type SearchFn func() *models.Device

func SimpleShark(
	setMonitorFn MonitorFn,
	searchFn SearchFn,
	sharkFn SharkFn) bool {

	cfg := state.GetConfig()

	// Look up the adapter. Use the find-ralink library.
	// The % will trigger first time through, which we want.
	var dev *models.Device = nil
	// If the config doesn't have this in it, we get a divide by zero.
	dev = searchFn()

	// Only do a reading and continue the pipeline
	// if we find an adapter.
	if dev != nil && dev.Exists {
		// Load the config for use.
		// cfg.Wireshark.Adapter = dev.Logicalname
		cfg.Log().Debug("found adapter: ", dev.Logicalname)
		setMonitorFn(dev)
		// This blocks for monitoring...
		macmap := sharkFn(dev.Logicalname)
		// Mark and remove too-short MAC addresses
		// for removal from the tshark findings.
		var keepers []string
		for _, mac := range macmap {
			if len(mac) >= constants.MACLENGTH {
				keepers = append(keepers, mac)
			}
		}
		StoreMacs(keepers)
	} else {
		cfg.Log().Info("no wifi devices found. no scanning carried out.")
		return false
	}
	return true
}

func StoreMacs(keepers []string) {
	cfg := state.GetConfig()
	// Do not log MAC addresses...
	cfg.Log().Debug("found ", len(keepers), " keepers")
	for _, mac := range keepers {
		state.RecordMAC(mac)
	}
}
