// Package structs holds internal state.
package structs

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"gsa.gov/18f/internal/interfaces"
)

type WifiEvents struct {
	Data []WifiEvent `json:"data"`
}

type WifiEvent struct {
	ID                int    `json:"rowid" db:"id" sqlite:"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL"`
	FCFSSeqID         string `json:"fcfs_seq_id" db:"fcfs_seq_id" type:"TEXT NOT NULL"`
	DeviceTag         string `json:"device_tag" db:"device_tag" type:"TEXT NOT NULL"`
	Localtime         string `json:"localtimestamp" db:"localtimestamp" type:"DATE NOT NULL"`
	SessionID         string `json:"session_id" db:"session_id" type:"TEXT NOT NULL"`
	ManufacturerIndex int    `json:"manufacturer_index" db:"manufacturer_index" type:"INTEGER NOT NULL"`
	PatronIndex       int    `json:"patron_index" db:"patron_index" type:"INTEGER NOT NULL"`
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
		if !strings.Contains(f.Tag.Get("type"), "AUTOINCREMENT") {
			col := strings.ReplaceAll(strings.Split(f.Tag.Get("db"), ",")[0], "\"", "")
			nom := strings.ReplaceAll(fmt.Sprintf("%v", reflect.Indirect(r).FieldByName(f.Name)), "\"", "")
			m[string(col)] = nom
		}
	}
	return m
}

func (wes WifiEvent) SelectAll(db interfaces.Database) []WifiEvent {
	we := []WifiEvent{}
	err := db.GetPtr().Select(&we, "SELECT * FROM WifiEvents")
	if err != nil {
		log.Println("Found no WifiEvents")
		log.Println(err.Error())
	}
	return we
}
