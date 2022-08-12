package state

type Duration struct {
	ID        int
	SessionID int64
	Start     int64
	End       int64
}

type DurationsDB struct {
	db map[int64]*Duration
}

var pk int64 = 0

func NewDurationsDB() *DurationsDB {
	db := make(map[int64]*Duration, 0)
	return &DurationsDB{db: db}
}

func (mkv *DurationsDB) ClearDurationsDB() {
	if mkv.db != nil {
		for k := range mkv.db {
			delete(mkv.db, k)
		}
	}
}

func (mkv *DurationsDB) Insert(d *Duration) {
	pk += 1
	mkv.db[pk] = d
}

func (mkv *DurationsDB) InsertMany(s int64, e EphemeralDB) {
	for _, ephemera := range e {
		mkv.Insert(&Duration{
			SessionID: s,
			Start:     ephemera.Start,
			End:       ephemera.End,
		})
	}
}

func (mkv *DurationsDB) GetSession(s int64) []*Duration {
	found := make([]*Duration, 0)
	for _, v := range mkv.db {
		if v.SessionID == s {
			found = append(found, v)
		}
	}

	return found
}

// Return the number of elements removed.
func (mkv *DurationsDB) DeleteSession(s int64) int {

	to_remove := make([]int64, 0)
	var size_before, size_after int

	size_before = len(mkv.db)

	// Find the hash keys to remove
	for k, v := range mkv.db {
		if v.SessionID == s {
			to_remove = append(to_remove, k)
		}
	}

	// Remove them from the hash "db"
	for _, v := range to_remove {
		delete(mkv.db, v)
	}

	size_after = len(mkv.db)

	return size_before - size_after
}
