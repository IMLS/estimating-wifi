package tempdb

import (
	"fmt"
	"path/filepath"
	"sync"

	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/logwrapper"
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

type List struct {
	*Queue
}

func NewList(cfg *config.Config, name string) *List {
	q := NewQueue(cfg, name)
	return &List{q}
}

func NewQueue(cfg *config.Config, name string) (q *Queue) {
	fullpath := filepath.Join(cfg.Local.WebDirectory, DURATIONSDB)
	tdb := NewSqliteDB(DURATIONSDB, fullpath)
	tdb.AddStructAsTable(name, QueueRow{})
	q = &Queue{name: name, db: tdb}
	return q
}

// The list abstraction is layered over the same table.
// Pushing to the list is the same as enqueuing w.r.t. the DB.
func (queue *Queue) Push(sessionid string) {
	queue.Enqueue(sessionid)
}

func (queue *Queue) AsList() []string {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()
	lw := logwrapper.NewLogger(nil)

	stmt := fmt.Sprintf("SELECT item FROM %v", queue.name)
	queue.db.Open()
	defer queue.db.Close()

	rows, err := queue.db.Ptr.Query(stmt)
	if err != nil {
		lw.Error("could not extract any items from queue/list ", queue.name)
	}

	sessions := make([]string, 0)
	for rows.Next() {
		var sid string
		rows.Scan(&sid)
		sessions = append(sessions, sid)
	}

	return sessions

}

func (queue *Queue) Remove(sessionid string) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()
	lw := logwrapper.NewLogger(nil)
	queue.db.Open()
	defer queue.db.Close()
	stmt, err := queue.db.Ptr.Prepare(fmt.Sprintf("DELETE FROM %v WHERE item = ?", queue.name))
	if err != nil {
		lw.Error("could not prepare delete statement for ", queue.name)
		lw.Fatal(err.Error())
	}
	// lw.Debug("removing ", sessionid)
	res, err := stmt.Exec(sessionid)
	if err != nil {
		lw.Error("could not delete session from queue/list ", sessionid, " ", queue.name)
		lw.Error(res)
		lw.Error(err.Error())
	}
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
