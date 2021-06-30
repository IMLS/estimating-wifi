package model

import (
	"fmt"

	"gsa.gov/18f/config"
	"gsa.gov/18f/session-counter/api"
)

// This probably should be a proper database.
type uniqueMappingDB struct {
	lastid  *int
	uid     map[string]int
	mfg     map[string]string
	anonmfg map[string]int
	tick    map[string]int
	cfg     *config.Config
}

func NewUMDB(cfg *config.Config) *uniqueMappingDB {
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

func (umdb uniqueMappingDB) UpdateMapping(mac string) {

	_, found := umdb.mfg[mac]
	// If we didn't find the mac we're supposed to update, then we need to add it.
	if !found {
		// Assign the next id.
		umdb.uid[mac] = *umdb.lastid
		// Increment for the next found address.
		*umdb.lastid = *umdb.lastid + 1
		// 20210412 MCJ
		// Now manufactuerers are being numbered as they come in.
		// This makes sure that we don't leak info. If the first device
		// we see after powerup is an "Apple" device, it will become
		// mfg "0". If the third device we see is an "Apple" device, then
		// Apple devices will be mfg 3. Effectively random, and does not
		// leak any info.

		// Get the actual manufactuerer. This pares down the MAC appropriately.
		// Grab a manufacturer for this MAC
		mfg := api.MacToMfg(umdb.cfg, mac)
		// Do we have a mfg mapping?
		// If we do, use it. If not, create a new mapping.
		mfgid, found := umdb.anonmfg[mfg]
		if !found {
			mfgid = len(umdb.anonmfg)
		}
		umdb.anonmfg[mfg] = mfgid
		umdb.mfg[mac] = mfg
		umdb.tick[mac] = 0
	} else {
		// If this address is already known, update
		// when we last saw it.
		umdb.tick[mac] = 0
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

func (umdb uniqueMappingDB) AsUserMappings() map[string]int {
	h := make(map[string]int)
	for mac := range umdb.mfg {
		userm := fmt.Sprintf("%v:%d", umdb.anonmfg[api.MacToMfg(umdb.cfg, mac)], umdb.uid[mac])
		h[userm] = umdb.tick[mac]
	}

	return h
}
