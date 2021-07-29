package tlp

import (
	"log"

	"gsa.gov/18f/internal/http"
	"gsa.gov/18f/internal/interfaces"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"
)

func SimpleSend(db interfaces.Database) {
	cfg := state.GetConfig()
	cfg.Log().Debug("Starting BatchSend")
	// This only comes in on reset...
	sq := state.NewQueue("sent")
	sessionsToSend := sq.AsList()

	for _, nextSessionIDToSend := range sessionsToSend {
		durations := []structs.Duration{}
		// FIXME: Leaky Abstraction
		err := db.GetPtr().Select(&durations, "SELECT * FROM durations WHERE session_id=?", nextSessionIDToSend)

		if err != nil {
			cfg.Log().Info("error in extracting durations for session", nextSessionIDToSend)
			cfg.Log().Error(err.Error())
		}
		cfg.Log().Debug("found ", len(durations), " durations to send in session ", nextSessionIDToSend)

		if len(durations) == 0 {
			cfg.Log().Info("found zero durations to send/draw. dequeing session [", nextSessionIDToSend, "]")
			sq.Remove(nextSessionIDToSend)
		} else if cfg.IsStoringToAPI() {
			cfg.Log().Info("attempting to send batch [", nextSessionIDToSend, "][", len(durations), "] to the API server")
			// convert []Duration to an array of map[string]interface{}
			data := make([]map[string]interface{}, 0)
			for _, duration := range durations {
				data = append(data, duration.AsMap())
			}
			// After writing images, we come back and try and send the data remotely.
			cfg.Log().Debug("PostJSONing ", len(data), " duration datas")
			err = http.PostJSON(cfg, cfg.GetDurationsURI(), data)
			if err != nil {
				log.Println("could not log to API")
				log.Println(err.Error())
			} else {
				// If we successfully sent the data remotely, we can now mark it is as sent.
				sq.Remove(nextSessionIDToSend)
			}
		} else {
			// Always dequeue. We're storing locally "for free" into the
			// durations table before trying to do the send.
			sq.Remove(nextSessionIDToSend)
		}
	}

}
