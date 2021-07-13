package analysis

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type WifiEvents struct {
	Data []WifiEvent `json:"data"`
}

// https://stackoverflow.com/questions/18635671/how-to-define-multiple-name-tags-in-a-struct
//EventId           int       `json:"event_id" db:"event_id"`
type WifiEvent struct {
	ID                int    `json:"id" db:"id" sqlite:"INTEGER PRIMARY KEY AUTOINCREMENT"`
	FCFSSeqId         string `json:"fcfs_seq_id" db:"fcfs_seq_id" sqlite:"TEXT NOT NULL"`
	DeviceTag         string `json:"device_tag" db:"device_tag" sqlite:"TEXT NOT NULL"`
	Localtime         string `json:"localtimestamp" db:"localtimestamp" sqlite:"DATE NOT NULL"`
	SessionId         string `json:"session_id" db:"session_id" sqlite:"TEXT NOT NULL"`
	ManufacturerIndex int    `json:"manufacturer_index" db:"manufacturer_index" sqlite:"INTEGER NOT NULL"`
	PatronIndex       int    `json:"patron_index" db:"patron_index" sqlite:"INTEGER NOT NULL"`
}

func (w WifiEvent) AsMap() map[string]interface{} {
	m := make(map[string]interface{})
	rt := reflect.TypeOf(w)
	if rt.Kind() != reflect.Struct {
		panic("bad type")
	}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		r := reflect.ValueOf(w)
		// log.Println("tag db", f.Tag.Get("db"))
		if !strings.Contains(f.Tag.Get("sqlite"), "AUTOINCREMENT") {
			col := strings.ReplaceAll(strings.Split(f.Tag.Get("db"), ",")[0], "\"", "")
			nom := strings.ReplaceAll(fmt.Sprintf("%v", reflect.Indirect(r).FieldByName(f.Name)), "\"", "")
			m[string(col)] = nom
		}
	}
	return m
}

type Counter struct {
	Patrons          int
	Devices          int
	Transients       int
	PatronMinutes    int
	DeviceMinutes    int
	TransientMinutes int
}

func NewCounter(minMinutes int, maxMinutes int) *Counter {
	patron_min_mins = float64(minMinutes)
	patron_max_mins = float64(maxMinutes)
	return &Counter{0, 0, 0, 0, 0, 0}
}

func (c *Counter) Add(field int, minutes int) {
	switch field {
	case Patron:
		c.Patrons += 1
		c.PatronMinutes += minutes
	case Device:
		c.Devices += 1
		c.DeviceMinutes += minutes
	case Transient:
		c.Transients += 1
		c.TransientMinutes += minutes
	}
}

type ByStart []Duration

func (a ByStart) Len() int { return len(a) }
func (a ByStart) Less(i, j int) bool {
	it, _ := time.Parse(time.RFC3339, a[i].Start)
	jt, _ := time.Parse(time.RFC3339, a[j].Start)
	return it.Before(jt)
}
func (a ByStart) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

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

func (d Duration) AsMap() map[string]interface{} {
	m := make(map[string]interface{})
	rt := reflect.TypeOf(d)
	if rt.Kind() != reflect.Struct {
		panic("bad type")
	}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		r := reflect.ValueOf(d)
		// log.Println("tag db", f.Tag.Get("db"))
		if !strings.Contains(f.Tag.Get("sqlite"), "AUTOINCREMENT") {
			col := strings.ReplaceAll(strings.Split(f.Tag.Get("db"), ",")[0], "\"", "")
			nom := strings.ReplaceAll(fmt.Sprintf("%v", reflect.Indirect(r).FieldByName(f.Name)), "\"", "")
			m[string(col)] = nom
		}
	}
	return m
}
