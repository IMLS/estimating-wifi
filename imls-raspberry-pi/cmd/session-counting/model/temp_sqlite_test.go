package model

import (
	"log"
	"reflect"
	"strings"
	"testing"
	"time"

	"gsa.gov/18f/analysis"
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

func getFieldName(tag, key string, s interface{}) (fieldname string) {
	rt := reflect.TypeOf(s)
	if rt.Kind() != reflect.Struct {
		panic("bad type")
	}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		v := strings.Split(f.Tag.Get(key), ",")[0]
		if v == tag {
			return f.Name
		}
	}
	return ""
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
		Start:     time.Now().Format(time.RFC3339),
		End:       time.Now().Format(time.RFC3339),
	}
	log.Println("asmap", d.AsMap())
	durations = append(durations, d)
	log.Println("durations", durations)
	tdb := NewSqliteDB("durations-test", "/tmp")
	tdb.AddTable("durations", map[string]string{
		"id":                 "INTEGER PRIMARY KEY AUTOINCREMENT",
		"pi_serial":          "TEXT",
		"fcfs_seq_id":        "TEXT",
		"device_tag":         "TEXT",
		"session_id":         "TEXT",
		"manufacturer_index": "INTEGER",
		"patron_index":       "INTEGER",
		"start":              "DATE",
		"end":                "DATE",
	})
	fields := tdb.GetFields("durations")
	log.Println("fields", fields)
	for _, d := range durations {
		values := make(map[string]interface{})
		// For each field name in the DB
		for _, field := range fields {
			// Get the struct field name.
			fieldname := getFieldName(field, "db", analysis.Duration{})
			if len(fieldname) > 0 {
				// Reflect on the duration
				r := reflect.ValueOf(d)
				v := reflect.Indirect(r).FieldByName(fieldname)
				log.Println("fieldname", fieldname, "value:", v)
				values[field] = v.String()
			}
		}
		tdb.Insert("durations", values)
	}

	tdb.DebugDump("durations")
}
