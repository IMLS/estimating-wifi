package tlp

import (
	"fmt"
	"os"
	"path/filepath"

	"gsa.gov/18f/cmd/session-counter/model"
	"gsa.gov/18f/internal/analysis"
	"gsa.gov/18f/internal/interfaces"
	"gsa.gov/18f/internal/logwrapper"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"
)

//This must happen after the data is updated for the day.
func writeImages(durations []structs.Duration, sessionid string) error {
	cfg := state.GetConfig()
	lw := logwrapper.NewLogger(nil)
	var reterr error

	if _, err := os.Stat(cfg.Paths.WWW.Root); os.IsNotExist(err) {
		err := os.Mkdir(cfg.Paths.WWW.Root, 0777)
		if err != nil {
			lw.Error("could not create web directory: ", cfg.Paths.WWW.Root)
			reterr = err
		}
	}
	if _, err := os.Stat(cfg.Paths.WWW.Images); os.IsNotExist(err) {
		err := os.Mkdir(cfg.Paths.WWW.Images, 0777)
		if err != nil {
			lw.Error("could not create image directory")
			reterr = err
		}
	}

	// FIXME: This filename kinda makes no sense if we're not running
	// a reset on a daily basis at midnight.
	yesterday := model.GetYesterday(cfg)
	imageFilename := fmt.Sprintf("%04d%02d%02d-%v-%v_%v.png",
		yesterday.Year(),
		int(yesterday.Month()),
		int(yesterday.Day()),
		sessionid,
		cfg.GetFCFSSeqID(),
		cfg.GetDeviceTag())

	path := filepath.Join(cfg.Paths.WWW.Images, imageFilename)
	// func DrawPatronSessions(cfg *config.Config, durations []Duration, outputPath string) {
	analysis.DrawPatronSessions(durations, path)
	return reterr
}

func WriteImages(ka *Keepalive, kb *KillBroker,
	chDurationsDB chan interfaces.Database) {
	cfg := state.GetConfig()
	cfg.Log().Debug("Starting WriteImages")
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
			cfg.Log().Debug("exiting WriteImages")
			return
		case db := <-chDurationsDB:
			iq := state.NewQueue("images")
			imagesToWrite := iq.AsList()
			cfg.Log().Debug("is there a session waiting to convert to images? [ ", imagesToWrite, "]")
			for _, nextImage := range imagesToWrite {
				durations := []structs.Duration{}
				cfg.Log().Debug("looking for session ", nextImage, " in durations table to write images")
				var count int
				db.GetPtr().QueryRow("SELECT COUNT(*) FROM durations WHERE session_id=?", nextImage).Scan(&count)
				cfg.Log().Debug("FOUND COUNT ", count)
				err := db.GetPtr().Select(&durations, "SELECT * FROM durations WHERE session_id=?", nextImage)
				cfg.Log().Debug("found ", len(durations), " durations in WriteImages")
				if err != nil {
					cfg.Log().Info("error in extracting durations for session", nextImage)
					cfg.Log().Error(err.Error())
				} else {
					err = writeImages(durations, nextImage)
					if err != nil {
						cfg.Log().Error("error in writing images")
						cfg.Log().Error(err.Error())
					} else {
						iq.Remove(nextImage)
					}
				}
			}
		}
	}
}
