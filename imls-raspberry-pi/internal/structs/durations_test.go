package structs

import (
	"testing"
	"time"
)

func TestAsMapDuration(t *testing.T) {
	e := Duration{
		Id:        1,
		PiSerial:  "asdf",
		DeviceTag: "abd-dc",
		Start:     time.Now().Format(time.RFC3339),
		End:       time.Now().Format(time.RFC3339),
		SessionId: "hello",
		MfgId:     0,
		PatronId:  0,
	}

	m := e.AsMap()
	if v, ok := m["id"]; ok {
		t.Log("map should not have `id` in it", v)
		t.Fail()
	}
}
