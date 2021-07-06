package model

import (
	"testing"
	"time"
)

func TestDBCreate(t *testing.T) {
	tdb := NewSqliteDB("test1", "/tmp/test")
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
	tdb := NewSqliteDB("test1", "/tmp/test")
	if tdb == nil {
		t.Log("failed to create tdb.")
		t.Fail()
	}
	tdb.AddTable("test", map[string]string{"a": "INTEGER", "b": "TEXT"})
	tdb.DropTable("test")
	tdb.Remove()
}

func TestInsert(t *testing.T) {
	tdb := NewSqliteDB("test1", "/tmp/test")
	if tdb == nil {
		t.Log("failed to create tdb.")
		t.Fail()
	}
	tdb.AddTable("test", map[string]string{"a": "INTEGER", "b": "TEXT"})
	tdb.Insert("test", map[string]interface{}{"a": "1", "b": "testing"})
	tdb.DebugDump("test")
}

func TestInsertAgain(t *testing.T) {
	tdb := NewSqliteDB("test1", "/tmp/test")
	if tdb == nil {
		t.Log("failed to create tdb.")
		t.Fail()
	}
	tdb.AddTable("test", map[string]string{"a": "INTEGER", "b": "TEXT"})
	tdb.Insert("test", map[string]interface{}{"a": "2", "b": time.Now().Format(time.RFC3339)})
	tdb.DebugDump("test")
}

func TestWifiTable(t *testing.T) {
	tdb := NewSqliteDB("wifi", "/tmp")
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
		"localtimestamp":     time.Now().Format(time.RFC3339),
		"session_id":         "asdfasdhjfkjw3er4kjwefr",
		"manufacturer_index": "3",
		"patron_index":       "4",
	})
	tdb.Insert("wifi", map[string]interface{}{
		"event_id":           "43",
		"fcfs_seq_id":        "ME0000-000",
		"device_tag":         "somewhere",
		"localtimestamp":     time.Now().Format(time.RFC3339),
		"session_id":         "asdfasdhjfkjw3er4kjwefr",
		"manufacturer_index": "32",
		"patron_index":       "42",
	})
}

func TestCleanup(t *testing.T) {
	tdb1 := NewSqliteDB("test1", "/tmp/test")
	tdb2 := NewSqliteDB("wifi", "/tmp")
	tdb1.Close()
	tdb2.Close()
	tdb1.Remove()
	tdb2.Remove()
}
