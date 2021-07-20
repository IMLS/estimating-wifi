// Package structs holds internal state.
package structs

import (
	"fmt"
	"reflect"
	"strings"
)

type WifiEvents struct {
	Data []WifiEvent `json:"data"`
}

type WifiEvent struct {
	ID                int    `json:"id" db:"id" sqlite:"INTEGER PRIMARY KEY AUTOINCREMENT"`
	FCFSSeqID         string `json:"fcfs_seq_id" db:"fcfs_seq_id" sqlite:"TEXT NOT NULL"`
	DeviceTag         string `json:"device_tag" db:"device_tag" sqlite:"TEXT NOT NULL"`
	Localtime         string `json:"localtimestamp" db:"localtimestamp" sqlite:"DATE NOT NULL"`
	SessionID         string `json:"session_id" db:"session_id" sqlite:"TEXT NOT NULL"`
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
