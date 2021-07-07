package tlp

import (
	"log"
	"encoding/json"

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
			//
			durations := make([]analysis.Duration, 0)
			err := db.Ptr.Select(&durations, "SELECT * FROM durations")
			if err != nil {
				lw.Info("error in extracting all durations")
				lw.Fatal(err.Error())
			}

			// convert []Duration to an array of map[string]string
			data := make([]map[string]string, 0)
			// this is awful. is there a better way?
			var inInterface map[string]string
			for _, duration := range durations {
				what, _ := json.Marshal(duration)
				json.Unmarshal(what, &inInterface)
				data = append(data, inInterface)
			}
			_, err = http.PostJSON(cfg, cfg.GetDataUri(), data)
			if err != nil {
				log.Println("could not log to API")
				log.Println(err.Error())
			}
		}
	}
}
