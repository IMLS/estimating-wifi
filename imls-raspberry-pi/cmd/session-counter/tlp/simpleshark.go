package tlp

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"gsa.gov/18f/cmd/session-counter/constants"
	"gsa.gov/18f/internal/interfaces"
	"gsa.gov/18f/internal/logwrapper"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"
	"gsa.gov/18f/internal/wifi-hardware-search/models"
)

func TShark2(adapter string) []string {
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

type SharkFn func(string) []string
type MonitorFn func(*models.Device)
type SearchFn func() *models.Device

// search.SetMonitorMode monitorFn
// searchFn search.SearchForMatchingDevice()
func SimpleShark(kb *KillBroker, in <-chan Ping,
	db interfaces.Database,
	setMonitorFn MonitorFn,
	searchFn SearchFn,
	sharkFn SharkFn) {
	cfg := state.GetConfig()
	cfg.Log().Info("STARTING THESHARK")
	var chKill chan interface{} = nil
	if kb != nil {
		chKill = kb.Subscribe()
	}

	adapter := ""
	var macMap map[string]int = make(map[string]int)
	var counter int = 0

	// Use the durations DB?
	// db := state.NewSqliteDB(cfg.GetDurationsDatabase().GetPath())

	for {
		select {

		case <-chKill:
			cfg.Log().Debug("EXITING")
			return

		case <-in:
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
				// Set monitor mode if the adapter changes.
				if dev.Logicalname != adapter {
					cfg.Log().Debug("setting monitor mode")
					setMonitorFn(dev)
					adapter = dev.Logicalname
				}

				// This will block for [cfg.Wireshark.Duration] seconds.
				macmap := sharkFn(dev.Logicalname)
				// Mark and remove too-short MAC addresses
				// for removal from the tshark findings.
				var keepers []int
				// for `k, _ :=` is the same as `for k :=`
				for _, k := range macmap {
					if len(k) >= constants.MACLENGTH {
						cfg.Log().Debug("keeping ", k, " is long enough")
						if _, ok := macMap[k]; !ok {
							macMap[k] = counter
							counter += 1
						}
						keepers = append(keepers, macMap[k])
					}
				}
				// out <- keepers
				StoreMacs(db, keepers)
			} else {
				cfg.Log().Info("no wifi devices found. no scanning carried out.")
				// out <- make([]string, 0)
			}
		}
	}
}

func macExists(db interfaces.Database, patronID int) bool {
	var d structs.Duration
	row := db.GetPtr().QueryRowx("SELECT patron_index FROM durations WHERE patron_index = ?", patronID)
	err := row.StructScan(&d)
	// Returns true if MAC found.
	return err == nil
}

func insertFirstSeen(db interfaces.Database, patronID int) {
	cfg := state.GetConfig()

	var d = structs.Duration{
		PiSerial:  cfg.GetSerial(),
		SessionID: fmt.Sprint(cfg.GetCurrentSessionID()),
		FCFSSeqID: cfg.GetFCFSSeqID(),
		DeviceTag: cfg.GetDeviceTag(),
		PatronID:  patronID,
		MfgID:     0,
		Start:     cfg.Clock.Now().Format(time.RFC3339),
		End:       cfg.Clock.Now().Format(time.RFC3339),
	}
	db.GetTableFromStruct(structs.Duration{}).InsertStruct(d)
}

func updateLastSeen(db interfaces.Database, patronID int) {
	cfg := state.GetConfig()
	_, err := db.GetPtr().Exec("UPDATE durations SET end = ? WHERE patron_index = ?",
		cfg.Clock.Now(),
		patronID)
	if err != nil {
		cfg.Log().Error("Could not update MAC end ", patronID)
	}
}

func StoreMacs(db interfaces.Database, keepers []int) {
	cfg := state.GetConfig()
	cfg.Log().Debug("keepers", keepers)
	for _, patronID := range keepers {
		if macExists(db, patronID) {
			cfg.Log().Debug(patronID, " exists")
			updateLastSeen(db, patronID)
		} else {
			cfg.Log().Debug(patronID, " being inserted")
			insertFirstSeen(db, patronID)
		}
	}
}
