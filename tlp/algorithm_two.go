package tlp

import (
	"gsa.gov/18f/session-counter/config"
	"gsa.gov/18f/session-counter/csp"
	"gsa.gov/18f/session-counter/model"
)

func AlgorithmTwo(ka *csp.Keepalive, cfg *config.Config, in <-chan map[string]int, out chan<- map[model.UserMapping]int, kill <-chan bool) {

}
