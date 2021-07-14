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

func NewCounter(cfg *config.Config, name string) *Counter {
	fullpath := filepath.Join(cfg.Local.WebDirectory, DURATIONSDB)
	tdb := NewSqliteDB(DURATIONSDB, fullpath)
	tdb.AddStructAsTable(name, QueueRow{})
	q := &Queue{name: name, db: tdb}
	q.Enqueue("0")
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
	_, err := q.db.Ptr.Exec(fmt.Sprintf("UPDATE %v SET item = 0", q.name))
	if err != nil {
		lw.Error(err.Error())
	}
	q.db.Close()
}

func (q *Counter) Increment() int {
	lw := logwrapper.NewLogger(nil)
	i := q.Value()
	n, _ := strconv.Atoi(fmt.Sprintf("%v", i))
	n = n + 1
	q.db.Open()
	_, err := q.db.Ptr.Exec(fmt.Sprintf("UPDATE %v SET item = ?", q.name), n)
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

func (q *Counter) Prev() int {
	n := q.Value()
	n = n - 1
	return n
}
