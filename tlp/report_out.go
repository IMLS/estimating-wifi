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
	http_error_count = 0

	svr := config.GetServer(cfg, service)
	tok, errGT := api.GetToken(svr)
	if errGT != nil {
		log.Println("report:", service, "error in token fetch")
		log.Println(errGT)
		http_error_count = http_error_count + 1
	} else {
		// If we had no problems getting a token, we can then report
		// the data to Directus.
		// First, grab an event ID.

		for uid, count := range h {
			go func(id string, cnt int) {
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

				err := api.StoreDeviceCount(cfg, svr, tok, session_id, id, cnt)
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

	session_id := 0

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
			for _, service := range []string{"directus", "reval"} {
				errCount, err := report(service, cfg, session_id, h)
				if err != nil {
					log.Println("reportout: error in reporting to", service)
					log.Println(err)
					http_error_count += errCount
				}
			}
			// FIXME Bump the session counter.
			// FIXME This should be the result of inserting an event.
			session_id = session_id + 1

		case <-time.After(time.Duration(cfg.Monitoring.HTTPErrorIntervalMins) * time.Minute):
			// If this much time has gone by, go ahead and reset the error count.
			if http_error_count != 0 {
				log.Printf("reportout: resetting http_error_count from %d to 0\n", http_error_count)
				http_error_count = 0
			}
		}
	}
}
