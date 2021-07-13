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
func writeImages(cfg *config.Config, durations []analysis.Duration) error {
	lw := logwrapper.NewLogger(nil)
	var reterr error

	yesterday := model.GetYesterday(cfg)
	if _, err := os.Stat(cfg.Local.WebDirectory); os.IsNotExist(err) {
		err := os.Mkdir(cfg.Local.WebDirectory, 0777)
		if err != nil {
			lw.Error("could not create web directory: ", cfg.Local.WebDirectory)
			reterr = err
		}
	}
	imagedir := filepath.Join(cfg.Local.WebDirectory, "images")
	if _, err := os.Stat(imagedir); os.IsNotExist(err) {
		err := os.Mkdir(imagedir, 0777)
		if err != nil {
			lw.Error("could not create image directory")
			reterr = err
		}
	}

	path := filepath.Join(imagedir, fmt.Sprintf("%04d-%02d-%02d-%v-%v-summary.png", yesterday.Year(), int(yesterday.Month()), int(yesterday.Day()), cfg.Auth.FCFSId, cfg.Auth.DeviceTag))
	// func DrawPatronSessions(cfg *config.Config, durations []Duration, outputPath string) {
	analysis.DrawPatronSessions(cfg, durations, path)
	return reterr
}

func WriteImages(ka *Keepalive, cfg *config.Config, kb *KillBroker,
	ch_durations_db chan *model.TempDB) {

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
			iq := model.NewQueue(cfg, "images")
			imagesToWrite := iq.AsList()
			lw.Debug("is there a session waiting to convert to images? [ ", imagesToWrite, "]")
			for _, nextImage := range imagesToWrite {
				durations := []analysis.Duration{}
				lw.Debug("looking for session ", nextImage, " in durations table to write images")
				db.Open()
				err := db.Ptr.Select(&durations, "SELECT * FROM durations WHERE session_id=?", nextImage)
				db.Close()
				lw.Debug("found ", len(durations), " durations in WriteImages")
				if err != nil {
					lw.Info("error in extracting durations for session", nextImage)
					lw.Error(err.Error())
				} else {
					err = writeImages(cfg, durations)
					if err != nil {
						lw.Error("error in writing images")
						lw.Error(err.Error())
					} else {
						iq.Remove(nextImage)
					}
				}
			}
		}
	}
}
