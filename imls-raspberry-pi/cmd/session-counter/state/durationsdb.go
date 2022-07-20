package state

import "gsa.gov/18f/cmd/session-counter/structs"

type DurationsDB struct {
	db map[int64]*structs.Duration
}

var pk int64 = 0

func NewDurationsDB() *DurationsDB {
	db := make(map[int64]*structs.Duration, 0)
	return &DurationsDB{db: db}
}

func (mkv *DurationsDB) Insert(d *structs.Duration) {
	pk += 1
	mkv.db[pk] = d
}

func (mkv *DurationsDB) InsertMany(ds []*structs.Duration) {
	for _, d := range ds {
		mkv.Insert(d)
	}
}
