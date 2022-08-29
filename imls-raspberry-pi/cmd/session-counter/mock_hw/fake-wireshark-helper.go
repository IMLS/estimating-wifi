package mock_hw

import (
	"math/rand"

	"gsa.gov/18f/cmd/session-counter/tlp"
	"gsa.gov/18f/internal/wifi-hardware-search/models"
)

var NUMMACSHELPER int
var NUMFOUNDPERMINUTEHELPER int
var consistentMACsHelper = make([]string, 0)

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

func runFakeWireshark(device string) []string {

	thisTime := rand.Intn(NUMFOUNDPERMINUTEHELPER)
	send := make([]string, thisTime)
	for i := 0; i < thisTime; i++ {
		send[i] = consistentMACsHelper[rand.Intn(len(consistentMACsHelper))]
	}
	return send
}

func FakeWiresharkHelper(numfoundperminute int, nummacs int) {
	// Create a pool of NUMMACS devices to draw from.
	// We will send NUMFOUNDPERMINUTE each minute
	NUMFOUNDPERMINUTEHELPER = numfoundperminute
	NUMMACSHELPER = nummacs
	consistentMACsHelper = make([]string, NUMMACSHELPER)
	for i := 0; i < NUMMACSHELPER; i++ {
		consistentMACsHelper[i] = generateFakeMac()
	}

	tlp.SimpleShark(
		// search.SetMonitorMode,
		func(d *models.Device) {},
		// search.SearchForMatchingDevice,
		func() *models.Device { return &models.Device{Exists: true, Logicalname: "fakewan0"} },
		// tlp.TSharkRunner
		runFakeWireshark)
}
