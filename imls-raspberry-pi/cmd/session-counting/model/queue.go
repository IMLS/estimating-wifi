package model

import (
	"fmt"
	"path/filepath"
	"sync"

	"gsa.gov/18f/config"
	"gsa.gov/18f/logwrapper"
	"gsa.gov/18f/session-counter/constants"
)

type Item interface{}

type QueueRow struct {
	Rowid int
	Item  string `db:"item" sqlite:"TEXT UNIQUE"`
}

type Queue struct {
	name  string
	db    *TempDB
	mutex sync.Mutex
}

func NewQueue(cfg *config.Config, name string) (q *Queue) {
	fullpath := filepath.Join(cfg.Local.WebDirectory, constants.DURATIONSDB)
	tdb := NewSqliteDB(constants.DURATIONSDB, fullpath)
	tdb.AddStructAsTable(name, QueueRow{})
	q = &Queue{name: name, db: tdb}
	return q
}

func (queue *Queue) Enqueue(sessionid string) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()
	lw := logwrapper.NewLogger(nil)

	//qr := QueueRow{Item: sessionid}
	//queue.db.InsertStruct(queue.name, qr)
	stmt := fmt.Sprintf("INSERT OR IGNORE INTO %v (item) VALUES (?)", queue.name)
	queue.db.Open()
	_, err := queue.db.Ptr.Exec(stmt, sessionid)
	queue.db.Close()
	if err != nil {
		lw.Error("error in enqueue insert")
		lw.Error(err.Error())
	}
}

func (queue *Queue) Peek() Item {
	lw := logwrapper.NewLogger(nil)
	queue.mutex.Lock()
	defer queue.mutex.Unlock()
	qr := QueueRow{}

	queue.db.Open()
	err := queue.db.Ptr.Get(&qr, fmt.Sprintf("SELECT rowid, item FROM %v ORDER BY rowid", queue.name))
	queue.db.Close()

	// The rowid value starts at 1. From the SQLite documentation.
	// if the Rowid is 0, we did not get anything back.
	if err != nil || qr.Item == "" || qr.Rowid == 0 {
		lw.Debug("nothing to peek at on the queue [", queue.name, "]")
		lw.Debug(err.Error())
		return nil
	} else {
		lw.Debug("PEEK found [ ", qr.Item, " ] on the queue [ ", queue.name, " ]")
		return qr.Item
	}
}

func (queue *Queue) Dequeue() Item {
	lw := logwrapper.NewLogger(nil)
	queue.mutex.Lock()
	defer queue.mutex.Unlock()
	qr := QueueRow{}

	queue.db.Open()
	err := queue.db.Ptr.Get(&qr, fmt.Sprintf("SELECT rowid, item FROM %v ORDER BY rowid", queue.name))
	if err != nil {
		queue.db.Close()
		return nil
	} else {
		res, err := queue.db.Ptr.Exec(fmt.Sprintf("DELETE FROM %v WHERE ROWID = ?", queue.name), qr.Rowid)
		if err != nil {
			queue.db.Close()
			lw.Error("failed to delete ", qr.Rowid, " in queue.Dequeue()")
			lw.Error(err.Error())
		}
		lw.Debug("DEQUEUE result ", res)
		lw.Debug("DEQUEUE removed [ ", qr.Item, " ] on the queue [ ", queue.name, " ]")
		queue.db.Close()
		return qr.Item
	}
}
