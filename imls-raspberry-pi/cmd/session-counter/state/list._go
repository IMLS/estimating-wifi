package state

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type List struct {
	*Queue
}

func NewList(name string) *List {
	q := NewQueue(name)
	return &List{q}
}

// Push is a list abstraction layered over the same table. Pushing to the list
// is the same as enqueuing w.r.t. the DB.
func (queue *Queue) Push(sessionid string) {
	queue.Enqueue(sessionid)
}

func (queue *Queue) AsList() []string {
	stmt := fmt.Sprintf("SELECT item FROM %v ORDER BY rowid", queue.name)
	rows, err := queue.db.GetPtr().Query(stmt)
	if err != nil {
		log.Error().
			Str("queue", queue.name).
			Msg("could not extract any items")
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
	stmt, err := queue.db.GetPtr().Prepare(fmt.Sprintf("DELETE FROM %v WHERE item = ?", queue.name))
	if err != nil {
		log.Fatal().
			Str("queue", queue.name).
			Err(err).
			Msg("could not prepare delete statement")
	}
	res, err := stmt.Exec(sessionid)
	if err != nil {
		log.Error().
			Str("result", fmt.Sprintf("%v", res)).
			Str("sessionid", sessionid).
			Str("queue", queue.name).
			Err(err).
			Msg("could not delete session")
	}
}
