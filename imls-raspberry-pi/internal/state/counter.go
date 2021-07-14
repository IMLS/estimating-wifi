package state

import (
	"fmt"
	"path/filepath"
	"strconv"

	"gsa.gov/18f/config"
	"gsa.gov/18f/logwrapper"
)

type Counter struct {
	*Queue
}

// All of this needs to be refactored down into the TempDB space...
// Some convoluted things happen here. :/
func NewCounter(cfg *config.Config, name string) *Counter {
	lw := logwrapper.NewLogger(nil)
	fullpath := filepath.Join(cfg.Local.WebDirectory, DURATIONSDB)
	tdb := NewSqliteDB(DURATIONSDB, fullpath)
	tdb.Open()
	_, err := tdb.Ptr.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %v", name))
	if err != nil {
		lw.Error(err.Error())
	}
	tdb.Close()
	tdb.AddStructAsTable(name, QueueRow{})
	q := &Queue{name: name, db: tdb}
	q.db.Open()
	defer q.db.Close()
	_, err = q.db.Ptr.Exec(fmt.Sprintf("INSERT INTO %v (item) VALUES (0)", name))
	if err != nil {
		lw.Error(err.Error())
	}

	return &Counter{q}
}

func GetCounter(cfg *config.Config, name string) *Counter {
	fullpath := filepath.Join(cfg.Local.WebDirectory, DURATIONSDB)
	tdb := NewSqliteDB(DURATIONSDB, fullpath)
	q := &Queue{name: name, db: tdb}
	return &Counter{q}
}

func (q *Counter) Reset() {
	lw := logwrapper.NewLogger(nil)

	q.db.Open()
	_, err := q.db.Ptr.Exec(fmt.Sprintf("UPDATE %v SET item = 0 WHERE rowid = 1", q.name))
	if err != nil {
		lw.Error(err.Error())
	}
	q.db.Close()
}

func (q *Counter) Increment() int {
	lw := logwrapper.NewLogger(nil)
	i := q.Value()
	//lw.Debug("Before increment", i)
	n, _ := strconv.Atoi(fmt.Sprintf("%v", i))
	n = n + 1
	q.db.Open()
	_, err := q.db.Ptr.Exec(fmt.Sprintf("UPDATE %v SET item = ? WHERE rowid = 1", q.name), n)
	if err != nil {
		lw.Error(err.Error())
	}
	q.db.Close()
	return n
}

func (q *Counter) Value() int {
	i := q.Peek()
	n, _ := strconv.Atoi(fmt.Sprintf("%v", i))
	return n
}

func (q *Counter) PrevValue() int {
	n := q.Value()
	n = n - 1
	return n
}
