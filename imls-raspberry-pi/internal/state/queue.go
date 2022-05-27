package state

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
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
	tdb := GetQueuesDatabase()
	t := tdb.InitTable(name)
	t.AddColumn("item", "TEXT UNIQUE")
	t.Create()
	q = &Queue{name: name, db: tdb}
	return q
}

func (queue *Queue) Enqueue(item string) {
	stmt := fmt.Sprintf("INSERT OR IGNORE INTO %v (item) VALUES (?)", queue.name)
	queue.db.Open()
	_, err := queue.db.GetPtr().Exec(stmt, item)
	if err != nil {
		log.Error().Err(err).Msg("could not insert")
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
	qr := QueueRow{}

	err := queue.db.GetPtr().Get(&qr, fmt.Sprintf("SELECT rowid, item FROM %v ORDER BY rowid", queue.name))
	if err != nil {
		return "", err
	} else {
		_, err := queue.db.GetPtr().Exec(fmt.Sprintf("DELETE FROM %v WHERE ROWID = ?", queue.name), qr.Rowid)
		if err != nil {
			log.Error().
				Int("rowid", qr.Rowid).
				Err(err).Msg("could not delete")
			return "", errors.New("failed to dequeue")
		}
		return qr.Item, nil
	}
}
