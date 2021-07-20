package tlp

import (
	"log"

	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/http"
	"gsa.gov/18f/internal/logwrapper"
	"gsa.gov/18f/internal/state"
	"gsa.gov/18f/internal/structs"
)

func BatchSend(ka *Keepalive, cfg *config.Config, kb *KillBroker,
	chDurationsDBIn <-chan *state.TempDB) {

	lw := logwrapper.NewLogger(nil)
	lw.Debug("Starting BatchSend")
	var ping, pong chan interface{} = nil, nil
	var chKill chan interface{} = nil
	if kb != nil {
		chKill = kb.Subscribe()
	} else {
		ping, pong = ka.Subscribe("GenerateDurations", 30)
	}

	for {
		select {
		case <-ping:
			pong <- "BatchSend"
		case <-chKill:
			lw.Debug("exiting BatchSend")
			return
		case db := <-chDurationsDBIn:
			// This only comes in on reset...
			sq := state.NewQueue(cfg, "sent")
			sessionsToSend := sq.AsList()

			for _, nextSessionIDToSend := range sessionsToSend {
				durations := []structs.Duration{}
				// lw.Debug("looking for session ", unsent.Session, " in durations table")
				db.Open()
				err := db.Ptr.Select(&durations, "SELECT * FROM durations WHERE session_id=?", nextSessionIDToSend)
				db.Close()

				if err != nil {
					lw.Info("error in extracting durations for session", nextSessionIDToSend)
					lw.Error(err.Error())
				}
				lw.Debug("found ", len(durations), " durations to send.")

				if len(durations) == 0 {
					lw.Info("found zero durations to send/draw. dequeing session [", nextSessionIDToSend, "]")
					sq.Remove(nextSessionIDToSend)
				} else if cfg.IsStoringToApi() {
					lw.Info("attempting to send batch [", nextSessionIDToSend, "][", len(durations), "] to the API server")
					// convert []Duration to an array of map[string]interface{}
					data := make([]map[string]interface{}, 0)
					for _, duration := range durations {
						data = append(data, duration.AsMap())
					}
					// After writing images, we come back and try and send the data remotely.
					lw.Debug("PostJSONing ", len(data), " duration datas")
					err = http.PostJSON(cfg, cfg.GetDataURI(), data)
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
	}
}
