package structs

import (
	"testing"
	"time"
)

func TestAsMapWifi(t *testing.T) {
	e := WifiEvent{
		ID:                1,
		FCFSSeqID:         "asdf",
		DeviceTag:         "abd-dc",
		Localtime:         time.Now().Format(time.RFC3339),
		SessionID:         "hello",
		ManufacturerIndex: 0,
		PatronIndex:       0,
	}

	m := e.AsMap()
	if v, ok := m["id"]; ok {
		t.Log("map should not have `id` in it", v)
		t.Fail()
	}
}
