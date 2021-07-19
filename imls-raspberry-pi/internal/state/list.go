package state

import "fmt"

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
	cfg := GetConfig()
	stmt := fmt.Sprintf("SELECT item FROM %v ORDER BY rowid", queue.name)
	rows, err := queue.db.GetPtr().Query(stmt)
	if err != nil {
		cfg.Log().Error("could not extract any items from queue/list ", queue.name)
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
	cfg := GetConfig()
	stmt, err := queue.db.GetPtr().Prepare(fmt.Sprintf("DELETE FROM %v WHERE item = ?", queue.name))
	if err != nil {
		cfg.Log().Error("could not prepare delete statement for ", queue.name)
		cfg.Log().Fatal(err.Error())
	}
	// lw.Debug("removing ", sessionid)
	res, err := stmt.Exec(sessionid)
	if err != nil {
		cfg.Log().Error("could not delete session from queue/list ", sessionid, " ", queue.name)
		cfg.Log().Error(res)
		cfg.Log().Error(err.Error())
	}
}
