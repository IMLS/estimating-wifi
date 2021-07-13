package model

import (
	"os"
	"testing"
	"time"

	"gsa.gov/18f/analysis"
)

func TestDBCreate(t *testing.T) {
	tdb := NewSqliteDB("test1", "/tmp/test1.sqlite")
	tdb.Open()
	if tdb == nil {
		t.Log("failed to create tdb.")
		t.Fail()
	} else {
		tdb.Close()
		// No file is created by the underlying library unless
		// you execute something into the DB
		// tdb.Remove()
	}
}

func TestCreateTable(t *testing.T) {
	tdb := NewSqliteDB("test1", "/tmp/test1.sqlite")
	if tdb == nil {
		t.Log("failed to create tdb.")
		t.Fail()
	}
	tdb.AddTable("test", map[string]string{"a": "INTEGER", "b": "TEXT"})
	tdb.DropTable("test")
	tdb.Remove()
}

func TestInsert(t *testing.T) {
	tdb := NewSqliteDB("test1", "/tmp/test1.sqlite")
	if tdb == nil {
		t.Log("failed to create tdb.")
		t.Fail()
	}
	tdb.AddTable("test", map[string]string{"a": "INTEGER", "b": "TEXT"})
	tdb.Insert("test", map[string]interface{}{"a": "1", "b": "testing"})
	tdb.DebugDump("test")
}

func TestInsertAgain(t *testing.T) {
	tdb := NewSqliteDB("test1", "/tmp/test1.sqlite")
	if tdb == nil {
		t.Log("failed to create tdb.")
		t.Fail()
	}
	tdb.AddTable("test", map[string]string{"a": "INTEGER", "b": "TEXT"})
	tdb.Insert("test", map[string]interface{}{"a": "2", "b": cfg.Clock.Now().Format(time.RFC3339)})
	tdb.DebugDump("test")
}

func TestWifiTable(t *testing.T) {
	tdb := NewSqliteDB("wifi", "/tmp/wifi.sqlite")
	tdb.AddTable("wifi", map[string]string{
		"id":                 "INTEGER PRIMARY KEY AUTOINCREMENT",
		"event_id":           "INTEGER",
		"fcfs_seq_id":        "TEXT",
		"device_tag":         "TEXT",
		"localtimestamp":     "DATE",
		"session_id":         "TEXT",
		"manufacturer_index": "INTEGER",
		"patron_index":       "INTEGER"})
	tdb.Insert("wifi", map[string]interface{}{
		"event_id":           "42",
		"fcfs_seq_id":        "ME0000-000",
		"device_tag":         "somewhere",
		"localtimestamp":     cfg.Clock.Now().Format(time.RFC3339),
		"session_id":         "asdfasdhjfkjw3er4kjwefr",
		"manufacturer_index": "3",
		"patron_index":       "4",
	})
	tdb.Insert("wifi", map[string]interface{}{
		"event_id":           "43",
		"fcfs_seq_id":        "ME0000-000",
		"device_tag":         "somewhere",
		"localtimestamp":     cfg.Clock.Now().Format(time.RFC3339),
		"session_id":         "asdfasdhjfkjw3er4kjwefr",
		"manufacturer_index": "32",
		"patron_index":       "42",
	})
}

func TestWifiTable2(t *testing.T) {
	tdb := NewSqliteDB("wifi2", "/tmp/wifi2.db")
	tdb.AddStructAsTable("wifi", analysis.WifiEvent{})
	w := analysis.WifiEvent{
		FCFSSeqId:         "ME0000-000",
		DeviceTag:         "another-tag",
		Localtime:         cfg.Clock.Now().Format(time.RFC3339),
		SessionId:         "asdfasdfasdf",
		ManufacturerIndex: 0,
		PatronIndex:       1,
	}
	tdb.InsertStruct("wifi", w)

	tdb.Insert("wifi", map[string]interface{}{
		"event_id":           "42",
		"fcfs_seq_id":        "ME0000-000",
		"device_tag":         "somewhere",
		"localtimestamp":     cfg.Clock.Now().Format(time.RFC3339),
		"session_id":         "asdfasdhjfkjw3er4kjwefr",
		"manufacturer_index": "3",
		"patron_index":       "4",
	})
	tdb.Insert("wifi", map[string]interface{}{
		"event_id":           "43",
		"fcfs_seq_id":        "ME0000-000",
		"device_tag":         "somewhere",
		"localtimestamp":     cfg.Clock.Now().Format(time.RFC3339),
		"session_id":         "asdfasdhjfkjw3er4kjwefr",
		"manufacturer_index": "32",
		"patron_index":       "42",
	})

	tdb.DebugDump("wifi")
	tdb.Remove()
}

func TestCleanup(t *testing.T) {
	tdb1 := NewSqliteDB("test1", "/tmp/test1.sqlite")
	tdb1.Open()
	tdb2 := NewSqliteDB("wifi", "/tmp/wifi.sqlite")
	tdb2.Open()
	tdb1.Close()
	tdb2.Close()
	tdb1.Remove()
	tdb2.Remove()
}

func TestReflection(t *testing.T) {
	durations := make([]*analysis.Duration, 0)
	d := &analysis.Duration{
		PiSerial:  "12345",
		SessionId: "asdf",
		FCFSSeqId: "ME0000-000",
		DeviceTag: "a-device-tag",
		PatronId:  1,
		MfgId:     1,
		Start:     cfg.Clock.Now().Format(time.RFC3339),
		End:       cfg.Clock.Now().Format(time.RFC3339),
	}
	durations = append(durations, d)

	os.Remove("/tmp/durations.sqlite")
	tdb := NewSqliteDB("durations", "/tmp/durations.sqlite")
	tdb.AddStructAsTable("durations", analysis.Duration{})

	for _, d := range durations {
		tdb.Insert("durations", d.AsMap())
	}
	tdb.DebugDump("durations")
}
