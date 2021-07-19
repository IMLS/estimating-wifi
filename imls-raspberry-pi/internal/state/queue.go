package state

import (
	"errors"
	"fmt"

	"gsa.gov/18f/internal/interfaces"
)

type QueueRow struct {
	Rowid int
	Item  string
}
type Queue struct {
	name string
	db   interfaces.Database
}

func NewQueue(name string) (q *Queue) {
	cfg := GetConfig()
	tdb := NewSqliteDB(cfg.Databases.QueuesPath)
	t := tdb.InitTable(name)
	t.AddColumn("item", "TEXT UNIQUE")
	t.Create()
	q = &Queue{name: name, db: tdb}
	return q
}

func (queue *Queue) Enqueue(item string) {
	cfg := GetConfig()
	stmt := fmt.Sprintf("INSERT OR IGNORE INTO %v (item) VALUES (?)", queue.name)
	queue.db.Open()
	_, err := queue.db.GetPtr().Exec(stmt, item)
	if err != nil {
		cfg.Log().Error("error in enqueue insert")
		cfg.Log().Error(err.Error())
	}
}

func (queue *Queue) Peek() (string, error) {
	//lw := logwrapper.NewLogger(nil)
	qr := QueueRow{}

	queue.db.Open()
	err := queue.db.GetPtr().Get(&qr, fmt.Sprintf("SELECT rowid, item FROM %v ORDER BY rowid", queue.name))

	// The rowid value starts at 1. From the SQLite documentation.
	// if the Rowid is 0, we did not get anything back.
	if err != nil || qr.Item == "" || qr.Rowid == 0 {
		return "", errors.New("nothing on the queue")
	} else {
		//lw.Debug("PEEK found [ ", qr.Item, " ] on the queue [ ", queue.name, " ]")
		return qr.Item, nil
	}
}

func (queue *Queue) Dequeue() (string, error) {
	cfg := GetConfig()
	qr := QueueRow{}

	err := queue.db.GetPtr().Get(&qr, fmt.Sprintf("SELECT rowid, item FROM %v ORDER BY rowid", queue.name))
	if err != nil {
		return "", err
	} else {
		_, err := queue.db.GetPtr().Exec(fmt.Sprintf("DELETE FROM %v WHERE ROWID = ?", queue.name), qr.Rowid)
		if err != nil {
			cfg.Log().Error("failed to delete ", qr.Rowid, " in queue.Dequeue()")
			cfg.Log().Error(err.Error())
			return "", errors.New("failed to dequeue")
		}
		// cfg.Log().Debug("DEQUEUE result ", res)
		// cfg.Log().Debug("DEQUEUE removed [ ", qr.Item, " ] on the queue [ ", queue.name, " ]")
		return qr.Item, nil
	}
}
