package state

import (
	"path/filepath"

	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/logwrapper"
)

func getCurrentSessionId(cfg *config.Config) int {
	lw := logwrapper.NewLogger(nil)
	fullpath := filepath.Join(cfg.Local.WebDirectory, DURATIONSDB)
	tdb := NewSqliteDB(DURATIONSDB, fullpath)
	if tdb.CheckTableExists("durations") {
		tdb.Open()
		defer tdb.Close()
		var sessionId int
		err := tdb.Ptr.Get(&sessionId, "SELECT IFNULL(MAX(session_id), 0) FROM durations")
		if err != nil {
			lw.Error("error in finding max session id; returning 0")
			lw.Error(err.Error())
			return 0
		}
		return sessionId
	} else {
		lw.Error("durations table did not exist; returning session id 0")
		return 0
	}
}

func GetCurrentSessionId(cfg *config.Config) int {
	return getCurrentSessionId(cfg)
}

func GetNextSessionId(cfg *config.Config) int {
	return cfg.SessionId + 1
}

func GetPreviousSessionId(cfg *config.Config) int {
	return cfg.SessionId - 1
}
