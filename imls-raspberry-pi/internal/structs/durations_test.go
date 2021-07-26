package structs

import (
	"testing"
	"time"
)

func TestAsMapDuration(t *testing.T) {
	e := Duration{
		// ID:        1,
		PiSerial:  "asdf",
		DeviceTag: "abd-dc",
		Start:     time.Now().Format(time.RFC3339),
		End:       time.Now().Format(time.RFC3339),
		SessionID: "hello",
		MfgID:     0,
		PatronID:  0,
	}

	m := e.AsMap()
	if v, ok := m["ID"]; ok {
		t.Log("map should not have `id` in it", v)
		t.Fail()
	}
}
