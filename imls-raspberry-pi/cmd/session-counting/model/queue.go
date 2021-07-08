package model

import (
	"fmt"
	"path/filepath"
	"sync"

	"gsa.gov/18f/config"
	"gsa.gov/18f/logwrapper"
)

type Item interface{}

type QueueRow struct {
	Rowid int
	Item  string `db:"item" sqlite:"TEXT"`
}

type Queue struct {
	name  string
	db    *TempDB
	mutex sync.Mutex
}

func NewQueue(cfg *config.Config, name string) (q *Queue) {
	durationsDB := "durations.sqlite"
	fullpath := filepath.Join(cfg.Local.WebDirectory, durationsDB)
	tdb := NewSqliteDB(durationsDB, fullpath)
	tdb.AddStructAsTable(name, QueueRow{})
	q = &Queue{name: name, db: tdb}
	return q
}

func (queue *Queue) Enqueue(sessionid string) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()
	qr := QueueRow{Item: sessionid}
	queue.db.InsertStruct(queue.name, qr)
}

func (queue *Queue) Peek() Item {
	lw := logwrapper.NewLogger(nil)
	queue.mutex.Lock()
	defer queue.mutex.Unlock()
	qr := QueueRow{}
	err := queue.db.Ptr.Get(&qr, fmt.Sprintf("SELECT rowid, item FROM %v ORDER BY rowid", queue.name))
	if err != nil {
		lw.Error(err.Error())
		return nil
	} else {
		return qr.Item
	}
}

func (queue *Queue) Dequeue() Item {
	lw := logwrapper.NewLogger(nil)
	queue.mutex.Lock()
	defer queue.mutex.Unlock()
	qr := QueueRow{}
	err := queue.db.Ptr.Get(&qr, fmt.Sprintf("SELECT rowid, item FROM %v ORDER BY rowid", queue.name))
	if err != nil {
		return nil
	} else {
		_, err := queue.db.Ptr.Exec(fmt.Sprintf("DELETE FROM %v WHERE ROWID = ?", queue.name), qr.Rowid)
		if err != nil {
			lw.Error("failed to delete ", qr.Rowid, " in queue.Dequeue()")
			lw.Error(err.Error())
		}

		return qr.Item
	}
}
