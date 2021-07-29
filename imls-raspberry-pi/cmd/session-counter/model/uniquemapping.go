// Package model manipulates internal state during execution
package model

import (
	"gsa.gov/18f/internal/interfaces"
)

// This probably should be a proper database.
type uniqueMappingDB struct {
	lastid  *int
	uid     map[string]int
	mfg     map[string]string
	anonmfg map[string]int
	tick    map[string]int
	cfg     interfaces.Config
}

func NewUMDB(cfg interfaces.Config) *uniqueMappingDB {
	umdb := &uniqueMappingDB{
		lastid:  new(int),
		uid:     make(map[string]int),
		mfg:     make(map[string]string),
		anonmfg: make(map[string]int),
		tick:    make(map[string]int),
		cfg:     cfg}
	return umdb
}

func (umdb uniqueMappingDB) WipeDB() {
	*umdb.lastid = 0
	umdb.uid = make(map[string]int)
	umdb.mfg = make(map[string]string)
	umdb.anonmfg = make(map[string]int)
	umdb.tick = make(map[string]int)
}

func (umdb uniqueMappingDB) AdvanceTime() {
	// Bump all the ticks by one.
	for mac := range umdb.mfg {
		umdb.tick[mac] = umdb.tick[mac] + 1
	}
}

func (umdb uniqueMappingDB) RemoveOldMappings(window int) {
	remove := make([]string, 0)
	// Find everything we need to remove.
	for mac := range umdb.mfg {
		if umdb.tick[mac] >= window {
			remove = append(remove, mac)
		}
	}
	// Remove everything that's old.
	// But not the mfg anonymization.
	// That is kept until reset.
	for _, mac := range remove {
		delete(umdb.uid, mac)
		delete(umdb.mfg, mac)
		delete(umdb.tick, mac)
	}
}
