package tlp

import (
	"log"
	"math/rand"
	"time"

	"gsa.gov/18f/session-counter/api"
	"gsa.gov/18f/session-counter/config"
	"gsa.gov/18f/session-counter/csp"
	"gsa.gov/18f/session-counter/model"
)

/* PROC reportMap
 * Takes a hashmap of [mfg id : count] and POSTs
 * each one to the server individually. We have no bulk insert.
 */
func ReportMap(ka *csp.Keepalive, cfg *config.Config, mfgs <-chan map[string]model.Entry) {
	log.Println("Starting reportMap")
	ping, pong := ka.Subscribe("reportMap", 5)

	var count int64 = 0
	http_error_count := 0
	directusServer := config.GetServer(cfg, "directus")

	for {
		select {
		case <-ping:
			// Every [cfg.Monitoring.HTTPErrorIntervalMins] this value
			// is reset to zero. If we get too many errors in that number of
			// minutes, then we should fail the next pong request. This will
			// kill the program, and we'll restart.
			if http_error_count < cfg.Monitoring.MaxHTTPErrorCount {
				pong <- "reportMap"
			} else {
				log.Printf("report: http_error_count threshold of %d reached\n", http_error_count)
			}

		case m := <-mfgs:
			count = count + 1
			log.Println("reporting: ", count)
			tok, errGT := api.GetToken(directusServer)
			if errGT != nil {
				log.Println("report: error in token fetch")
				log.Println(errGT)
				http_error_count = http_error_count + 1
			}

			for _, entry := range m {
				go func(entry model.Entry) {
					// Smear the requests out in time.
					time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
					err := api.Report_mfg(directusServer, tok, entry)
					if err != nil {
						log.Println("report: results POST failure")
						log.Println(err)
						http_error_count = http_error_count + 1
					}
				}(entry)
			}

			errRT := api.Report_telemetry(directusServer, tok)
			if errRT != nil {
				log.Println("report: error in telemetry POST")
				log.Println(errRT)
				http_error_count = http_error_count + 1
			}

		case <-time.After(time.Duration(cfg.Monitoring.HTTPErrorIntervalMins) * time.Minute):
			// If this much time has gone by, go ahead and reset the error count.
			if http_error_count != 0 {
				log.Printf("report: resetting http_error_count from %d to 0\n", http_error_count)
				http_error_count = 0
			}
		}
	}
}
