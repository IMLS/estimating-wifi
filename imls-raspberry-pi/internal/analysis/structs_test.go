package analysis

import (
	"testing"
	"time"
)

/*
type WifiEvent struct {
	ID                int    `json:"id" db:"id" sqlite:"INTEGER PRIMARY KEY AUTOINCREMENT"`
	FCFSSeqId         string `json:"fcfs_seq_id" db:"fcfs_seq_id" sqlite:"TEXT NOT NULL"`
	DeviceTag         string `json:"device_tag" db:"device_tag" sqlite:"TEXT NOT NULL"`
	Localtime         string `json:"localtimestamp" db:"localtimestamp" sqlite:"DATE NOT NULL"`
	SessionId         string `json:"session_id" db:"session_id" sqlite:"TEXT NOT NULL"`
	ManufacturerIndex int    `json:"manufacturer_index" db:"manufacturer_index" sqlite:"INTEGER NOT NULL"`
	PatronIndex       int    `json:"patron_index" db:"patron_index" sqlite:"INTEGER NOT NULL"`
}
*/

func TestAsMapWifi(t *testing.T) {
	e := WifiEvent{
		ID:                1,
		FCFSSeqId:         "asdf",
		DeviceTag:         "abd-dc",
		Localtime:         time.Now().Format(time.RFC3339),
		SessionId:         "hello",
		ManufacturerIndex: 0,
		PatronIndex:       0,
	}

	m := e.AsMap()
	if v, ok := m["id"]; ok {
		t.Log("map should not have `id` in it", v)
		t.Fail()
	}
}

/*

type Duration struct {
	Id        int    `db:"id" sqlite:"INTEGER PRIMARY KEY AUTOINCREMENT"`
	PiSerial  string `db:"pi_serial" sqlite:"TEXT"`
	SessionId string `db:"session_id" sqlite:"TEXT"`
	FCFSSeqId string `db:"fcfs_seq_id" sqlite:"TEXT"`
	DeviceTag string `db:"device_tag" sqlite:"TEXT"`
	PatronId  int    `db:"patron_index" sqlite:"INTEGER"`
	MfgId     int    `db:"manufacturer_index" sqlite:"INTEGER"`
	Start     string `db:"start" sqlite:"DATE"`
	End       string `db:"end" sqlite:"DATE"`
}
*/

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
