package state

import (
	"fmt"
	"path/filepath"
	"strconv"

	"gsa.gov/18f/config"
)

type Counter struct {
	*Queue
}

func NewCounter(cfg *config.Config, name string) *Counter {
	fullpath := filepath.Join(cfg.Local.WebDirectory, DURATIONSDB)
	tdb := NewSqliteDB(DURATIONSDB, fullpath)
	tdb.AddStructAsTable(name, QueueRow{})
	q := &Queue{name: name, db: tdb}
	return &Counter{q}
}

func GetCounter(cfg *config.Config, name string) *Counter {
	fullpath := filepath.Join(cfg.Local.WebDirectory, DURATIONSDB)
	tdb := NewSqliteDB(DURATIONSDB, fullpath)
	q := &Queue{name: name, db: tdb}
	return &Counter{q}
}

func (q *Counter) Reset() {
	q.db.Open()
	q.db.Ptr.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %v", q.name))
	q.db.Close()
	q.Enqueue("0")
}

func (q *Counter) Increment() int {
	i := q.Dequeue()
	n, _ := strconv.Atoi(fmt.Sprintf("%v", i))
	n = n + 1
	q.Enqueue(fmt.Sprintf("%v", n))
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
