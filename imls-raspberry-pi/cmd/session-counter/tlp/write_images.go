package tlp

import (
	"fmt"
	"os"
	"path/filepath"

	"gsa.gov/18f/cmd/session-counter/model"
	"gsa.gov/18f/internal/analysis"
	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/logwrapper"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"
)

//This must happen after the data is updated for the day.
func writeImages(cfg *config.Config, durations []structs.Duration) error {
	lw := logwrapper.NewLogger(nil)
	var reterr error

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

	// FIXME: This filename kinda makes no sense if we're not running
	// a reset on a daily basis at midnight.
	yesterday := model.GetYesterday(cfg)
	imageFilename := fmt.Sprintf("%04d-%02d-%02d-%v-%v-%v-summary.png",
		yesterday.Year(),
		int(yesterday.Month()),
		int(yesterday.Day()),
		state.GetPreviousSessionID(cfg),
		cfg.Auth.FCFSId,
		cfg.Auth.DeviceTag)

	path := filepath.Join(imagedir, imageFilename)
	// func DrawPatronSessions(cfg *config.Config, durations []Duration, outputPath string) {
	analysis.DrawPatronSessions(cfg, durations, path)
	return reterr
}

func WriteImages(ka *Keepalive, cfg *config.Config, kb *KillBroker,
	chDurationsDB chan *state.TempDB) {

	lw := logwrapper.NewLogger(nil)
	lw.Debug("Starting WriteImages")
	var ping, pong chan interface{} = nil, nil
	var chKill chan interface{} = nil
	if kb != nil {
		chKill = kb.Subscribe()
	} else {
		ping, pong = ka.Subscribe("WriteImages", 30)
	}

	for {
		select {
		case <-ping:
			pong <- "WriteImages"
		case <-chKill:
			lw.Debug("exiting WriteImages")
			return
		case db := <-chDurationsDB:
			iq := state.NewQueue(cfg, "images")
			imagesToWrite := iq.AsList()
			lw.Debug("is there a session waiting to convert to images? [ ", imagesToWrite, "]")
			for _, nextImage := range imagesToWrite {
				durations := []structs.Duration{}
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
