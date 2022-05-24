package tlp

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog/log"
	"gsa.gov/18f/internal/analysis"
	"gsa.gov/18f/internal/interfaces"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"
)

//This must happen after the data is updated for the day.
func writeImages(durations []structs.Duration, sessionid string) error {
	cfg := state.GetConfig()
	var reterr error

	if _, err := os.Stat(cfg.GetWWWRoot()); os.IsNotExist(err) {
		err := os.Mkdir(cfg.GetWWWRoot(), 0777)
		if err != nil {
			log.Error().
				Str("web directory", cfg.GetWWWRoot()).
				Msg("could not create root directory")
			reterr = err
		}
	}
	if _, err := os.Stat(cfg.GetWWWImages()); os.IsNotExist(err) {
		err := os.Mkdir(cfg.GetWWWImages(), 0777)
		if err != nil {
			log.Error().
				Str("web directory", cfg.GetWWWImages()).
				Msg("could not create image directory")
			reterr = err
		}
	}

	// FIXME: This filename kinda makes no sense if we're not running
	// a reset on a daily basis at midnight.
	// yesterday := model.GetYesterday(cfg)
	yesterday := state.GetClock().Now().In(time.Local)
	imageFilename := fmt.Sprintf("%04d%02d%02d-%v-%v_%v.png",
		yesterday.Year(),
		int(yesterday.Month()),
		int(yesterday.Day()),
		sessionid,
		cfg.GetFCFSSeqID(),
		cfg.GetDeviceTag())

	path := filepath.Join(cfg.GetWWWImages(), imageFilename)
	// func DrawPatronSessions(cfg *config.Config, durations []Duration, outputPath string) {
	analysis.DrawPatronSessions(durations, path)
	return reterr
}

func WriteImages(db interfaces.Database) {
	iq := state.NewQueue("images")
	imagesToWrite := iq.AsList()

	log.Info().
		Int("sessions", len(imagesToWrite)).
		Msg("about to write to image queue")

	for _, nextImage := range imagesToWrite {
		durations := []structs.Duration{}
		var count int
		db.GetPtr().QueryRow("SELECT COUNT(*) FROM durations WHERE session_id=?", nextImage).Scan(&count)
		err := db.GetPtr().Select(&durations, "SELECT * FROM durations WHERE session_id=?", nextImage)

		log.Debug().
			Int("durations", len(durations)).
			Msg("found durations")
		if err != nil {
			log.Error().
				Err(err).
				Str("session", nextImage).
				Msg("could not extract durations")
		} else {
			err = writeImages(durations, nextImage)
			if err != nil {
				log.Error().
					Err(err).
					Msg("could not write images")
			} else {
				iq.Remove(nextImage)
			}
		}
	}
}
