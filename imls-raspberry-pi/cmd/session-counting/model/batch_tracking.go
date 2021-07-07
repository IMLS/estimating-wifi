package model

import "gsa.gov/18f/logwrapper"

type Batch struct {
	Session string `db:"session" sqlite:"TEXT"`
	Sent    int    `db:"sent" sqlite:"INTEGER"`
}

func (tdb *TempDB) GetUnsentBatches() []Batch {
	batches := []Batch{}
	tdb.Ptr.Select(&batches, "SELECT DISTINCT session, sent FROM batches WHERE sent = 0")
	return batches
}

/*
UPDATE table_name
SET column1 = value1, column2 = value2, ...
WHERE condition;
*/
func (tdb *TempDB) MarkAsSent(b Batch) {
	lw := logwrapper.NewLogger(nil)
	unsent_count_before := len(tdb.GetUnsentBatches())
	lw.Debug("Marking as sent ", b)
	_, err := tdb.Ptr.Exec("UPDATE batches SET sent = 1 WHERE session = ?", b.Session)
	if err != nil {
		lw.Error("could not update batch ", b.Session, " as sent.")
	}
	unsent_count_after := len(tdb.GetUnsentBatches())
	lw.Debug("unsent before [", unsent_count_before, "] unsent after [", unsent_count_after, "]")
}
