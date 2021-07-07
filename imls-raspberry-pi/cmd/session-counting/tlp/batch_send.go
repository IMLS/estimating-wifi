package tlp

import (
	"log"

	"gsa.gov/18f/analysis"
	"gsa.gov/18f/config"
	"gsa.gov/18f/http"
	"gsa.gov/18f/logwrapper"
	"gsa.gov/18f/session-counter/model"
)

func BatchSend(ka *Keepalive, cfg *config.Config, kb *Broker,
	ch_durations_db chan *model.TempDB) {

	lw := logwrapper.NewLogger(nil)
	lw.Debug("Starting BatchSend")
	var ping, pong chan interface{} = nil, nil
	var ch_kill chan interface{} = nil
	if kb != nil {
		ch_kill = kb.Subscribe()
	} else {
		ping, pong = ka.Subscribe("GenerateDurations", 30)
	}

	for {
		select {
		case <-ping:
			pong <- "BatchSend"
		case <-ch_kill:
			lw.Debug("exiting BatchSend")
			return
		case db := <-ch_durations_db:
			unsents := db.GetUnsentBatches()
			lw.Debug("found ", len(unsents), " batches that are unsent")
			for _, unsent := range unsents {
				durations := []analysis.Duration{}
				lw.Debug("looking for session ", unsent.Session, " in durations table")
				err := db.Ptr.Select(&durations, "SELECT * FROM durations WHERE session_id=?", unsent.Session)
				if err != nil {
					lw.Info("error in extracting durations for session", unsent.Session)
					lw.Error(err.Error())
				}
				lw.Debug("found ", len(durations), " durations to send.")

				// convert []Duration to an array of map[string]interface{}
				data := make([]map[string]interface{}, 0)
				for _, duration := range durations {
					data = append(data, duration.AsMap())
				}
				lw.Debug("PostJSONing ", len(data), " duration datas")
				_, err = http.PostJSON(cfg, cfg.GetDurationsUri(), data)
				if err != nil {
					log.Println("could not log to API")
					log.Println(err.Error())
				} else {
					db.MarkAsSent(unsent)
				}
			}

		}
	}
}
