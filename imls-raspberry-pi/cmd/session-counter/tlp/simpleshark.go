package tlp

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"syscall"

	"gsa.gov/18f/cmd/session-counter/constants"
	"gsa.gov/18f/internal/interfaces"
	"gsa.gov/18f/internal/logwrapper"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"
	"gsa.gov/18f/internal/wifi-hardware-search/models"
)

func TSharkRunner(adapter string) []string {
	cfg := state.GetConfig()
	lw := logwrapper.NewLogger(nil)
	tsharkCmd := exec.Command(cfg.Executables.Wireshark.Path,
		"-a", fmt.Sprintf("duration:%d", cfg.Executables.Wireshark.Duration),
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

func SimpleShark(db interfaces.Database,
	// setMonitorFn func(d *models.Device),
	// searchFn func() *models.Device,
	// sharkFn func(dev string) []string) bool {
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
		StoreMacs(db, keepers)
	} else {
		cfg.Log().Info("no wifi devices found. no scanning carried out.")
		return false
	}
	return true
}

func macExists(db interfaces.Database, mac string) bool {
	var ed structs.EphemeralDuration
	row := db.GetPtr().QueryRowx("SELECT mac FROM ephemeraldurations WHERE mac = ?", mac)
	err := row.StructScan(&ed)
	// Returns true if MAC found.
	return err == nil && ed.MAC == mac
}

func insertFirstSeen(db interfaces.Database, mac string) {
	cfg := state.GetConfig()

	_, err := db.GetPtr().Exec("INSERT INTO ephemeraldurations (mac, start, end) VALUES (?, ?, ?)",
		mac,
		cfg.Clock.Now().Unix(),
		cfg.Clock.Now().Unix())
	if err != nil {
		cfg.Log().Error("Could not do initial insert for ", mac)
		cfg.Log().Fatal(err.Error())
	}
}

func updateLastSeen(db interfaces.Database, mac string) {
	cfg := state.GetConfig()
	_, err := db.GetPtr().Exec(`UPDATE ephemeraldurations SET end=? WHERE mac=?`,
		cfg.Clock.Now().Unix(),
		mac)
	if err != nil {
		cfg.Log().Fatal("Could not update MAC end for ", mac)
	}
}

func StoreMacs(db interfaces.Database, keepers []string) {
	cfg := state.GetConfig()
	cfg.Log().Debug("keepers", keepers)
	for _, mac := range keepers {
		if macExists(db, mac) {
			cfg.Log().Debug(mac, " exists")
			updateLastSeen(db, mac)
		} else {
			cfg.Log().Debug(mac, " being inserted")
			insertFirstSeen(db, mac)
		}
	}
}
