package structs

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type ByStart []Duration

func (a ByStart) Len() int { return len(a) }
func (a ByStart) Less(i, j int) bool {
	it, _ := time.Parse(time.RFC3339, a[i].Start)
	jt, _ := time.Parse(time.RFC3339, a[j].Start)
	return it.Before(jt)
}
func (a ByStart) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type Durations struct {
	Data []Duration `json:"data"`
}

type Duration struct {
	ID        int    `json:"id" db:"id" type:"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL"`
	PiSerial  string `json:"pi_serial" db:"pi_serial" type:"TEXT"`
	SessionID string `json:"session_id" db:"session_id" type:"TEXT"`
	FCFSSeqID string `json:"fcfs_seq_id" db:"fcfs_seq_id" type:"TEXT"`
	DeviceTag string `json:"device_tag" db:"device_tag" type:"TEXT"`
	PatronID  int    `json:"patron_index" db:"patron_index" type:"INTEGER"`
	Start     string `json:"start,string" db:"start" type:"DATE"`
	End       string `json:"end,string" db:"end" type:"DATE"`
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
		if !strings.Contains(f.Tag.Get("type"), "AUTOINCREMENT") {
			col := strings.ReplaceAll(strings.Split(f.Tag.Get("db"), ",")[0], "\"", "")
			nom := strings.ReplaceAll(fmt.Sprintf("%v", reflect.Indirect(r).FieldByName(f.Name)), "\"", "")
			m[string(col)] = nom
		}
	}
	return m
}
