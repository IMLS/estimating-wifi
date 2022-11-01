package mock_hw

import (
	"fmt"
	"math/rand"

	"github.com/rs/zerolog/log"
	"gsa.gov/18f/internal/config"
)

var consistentMACarray = make([]string, 0)

func generateFakeMac() string {
	var letterRunes = []rune("ABCDEF0123456789")
	b := make([]rune, 17)
	colons := [...]int{2, 5, 8, 11, 14}
	for i := 0; i < 17; i++ {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]

		for v := range colons {
			if i == colons[v] {
				b[i] = rune(':')
			}
		}
	}
	return string(b)
}

func RunFakeWireshark(device string) []string {
	minimum := config.GetFakesharkMinFound()
	maximum := config.GetFakesharkMaxFound()
	thisTime := rand.Intn(maximum-minimum) + minimum
	send := make([]string, thisTime)
	for i := 0; i < thisTime; i++ {
		send[i] = consistentMACarray[rand.Intn(len(consistentMACarray))]
	}

	log.Debug().
		Str("number of devices generated", fmt.Sprint(len(send))).
		Msg("runFakeWireshark")

	return send
}

func FakeWiresharkSetup() {

	num_macs := config.GetFakesharkNumMacs()
	// Create a pool of NUMMACS devices to draw from.
	// We will send NUMFOUNDPERMINUTE each minute
	consistentMACarray = make([]string, num_macs)
	for i := 0; i < num_macs; i++ {
		consistentMACarray[i] = generateFakeMac()
	}

}
