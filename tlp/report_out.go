package tlp

import (
	"log"
	"math/rand"
	"time"

	"gsa.gov/18f/session-counter/api"
	"gsa.gov/18f/session-counter/config"
	"gsa.gov/18f/session-counter/csp"
)

func report(service string, cfg *config.Config, session_id int, h map[string]int) (http_error_count int, err error) {
	tok, errGT := config.ReadAuth()
	http_error_count = 0

	if errGT != nil {
		log.Println("report:", service, "error in token fetch")
		log.Println(errGT)
		http_error_count = http_error_count + 1
	} else {
		for uid, count := range h {
			go func(id string, cnt int) {
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
				err := api.StoreDeviceCount(cfg, tok, session_id, id, cnt)
				if err != nil {
					log.Println("report:", service, "results POST failure")
					log.Println(err)
					http_error_count = http_error_count + 1
				}
			}(uid, count)
		}
	}

	return http_error_count, nil
}

func ReportOut(ka *csp.Keepalive, cfg *config.Config, ch_uidmap <-chan map[string]int) {
	log.Println("Starting ReportOut")
	ping, pong := ka.Subscribe("ReportOut", 30)
	http_error_count := 0

	// For event logging
	el := api.NewEventLogger(cfg)

	for {
		select {
		case <-ping:
			// Every [cfg.Monitoring.HTTPErrorIntervalMins] this value
			// is reset to zero. If we get too many errors in that number of
			// minutes, then we should fail the next pong request. This will
			// kill the program, and we'll restart.
			if http_error_count < cfg.Monitoring.MaxHTTPErrorCount {
				pong <- "ReportOut"
			} else {
				log.Printf("reportout: http_error_count threshold of %d reached\n", http_error_count)
			}
		// This is the [ uid -> ticks ] map (uid looks like "Next:0")
		case h := <-ch_uidmap:
			event_ndx := el.Log("logging_devices", nil)

			// This used to loop over "directus" and "reval"
			// We decided we will log only to reval, and it will handle validation and logging.
			service := "reval"
			errCount, err := report(service, cfg, event_ndx, h)
			if err != nil {
				log.Println("reportout: error in reporting to", service)
				log.Println(err)
				http_error_count += errCount
			}

		case <-time.After(time.Duration(cfg.Monitoring.HTTPErrorIntervalMins) * time.Minute):
			// If this much time has gone by, go ahead and reset the error count.
			if http_error_count != 0 {
				log.Printf("reportout: resetting http_error_count from %d to 0\n", http_error_count)
				http_error_count = 0
			}
		}
	}
}
