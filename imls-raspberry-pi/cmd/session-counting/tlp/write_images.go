package tlp

import (
	"fmt"
	"os"
	"path/filepath"

	"gsa.gov/18f/analysis"
	"gsa.gov/18f/config"
	"gsa.gov/18f/logwrapper"
	"gsa.gov/18f/session-counter/model"
)

//This must happen after the data is updated for the day.
func writeImages(cfg *config.Config, durations []analysis.Duration) {
	lw := logwrapper.NewLogger(nil)
	yesterday := model.GetYesterday()
	if _, err := os.Stat(cfg.Local.WebDirectory); os.IsNotExist(err) {
		err := os.Mkdir(cfg.Local.WebDirectory, 0777)
		if err != nil {
			lw.Error("could not create web directory: %v", cfg.Local.WebDirectory)
		}
	}
	imagedir := filepath.Join(cfg.Local.WebDirectory, "images")
	if _, err := os.Stat(imagedir); os.IsNotExist(err) {
		err := os.Mkdir(imagedir, 0777)
		if err != nil {
			lw.Error("could not create image directory")
		}
	}

	path := filepath.Join(imagedir, fmt.Sprintf("%04d-%02d-%02d-%v-%v-summary.png", yesterday.Year(), int(yesterday.Month()), int(yesterday.Day()), cfg.Auth.FCFSId, cfg.Auth.DeviceTag))
	// func DrawPatronSessions(cfg *config.Config, durations []Duration, outputPath string) {
	analysis.DrawPatronSessions(cfg, durations, path)
}

func WriteImages(ka *Keepalive, cfg *config.Config, kb *KillBroker,
	ch_durations_db chan *model.TempDB,
	ch_proceed chan<- Ping) {

	lw := logwrapper.NewLogger(nil)
	lw.Debug("Starting WriteImages")
	var ping, pong chan interface{} = nil, nil
	var ch_kill chan interface{} = nil
	if kb != nil {
		ch_kill = kb.Subscribe()
	} else {
		ping, pong = ka.Subscribe("WriteImages", 30)
	}

	for {
		select {
		case <-ping:
			pong <- "WriteImages"
		case <-ch_kill:
			lw.Debug("exiting WriteImages")
			return
		case db := <-ch_durations_db:
			unsents := db.GetUnsentBatches()
			lw.Debug("write images found ", len(unsents), " unsent batches to write as images.")
			for _, unsent := range unsents {
				durations := []analysis.Duration{}
				lw.Debug("looking for session ", unsent.Session, " in durations table to write images")
				err := db.Ptr.Select(&durations, "SELECT * FROM durations WHERE session_id=?", unsent.Session)
				if err != nil {
					lw.Info("error in extracting durations for session", unsent.Session)
					lw.Error(err.Error())
				}
				writeImages(cfg, durations)
				ch_proceed <- Ping{}
			}
		}
	}
}
