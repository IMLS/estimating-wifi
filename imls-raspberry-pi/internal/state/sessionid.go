package state

import (
	"path/filepath"

	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/logwrapper"
)

func getCurrentSessionID(cfg *config.Config) int {
	lw := logwrapper.NewLogger(nil)
	fullpath := filepath.Join(cfg.Local.WebDirectory, DURATIONSDB)
	tdb := NewSqliteDB(DURATIONSDB, fullpath)
	if tdb.CheckTableExists("durations") {
		tdb.Open()
		defer tdb.Close()
		var sessionID int
		err := tdb.Ptr.Get(&sessionID, "SELECT MAX(session_id) FROM durations")
		if err != nil {
			lw.Error("error in finding max session id; returning 0")
			lw.Error(err.Error())
			return 0
		}
		return sessionID
	} else {
		lw.Error("durations table did not exist; returning session id 0")
		return 0
	}
}

func GetInitialSessionID(cfg *config.Config) int {
	return getCurrentSessionID(cfg)
}

func GetCurrentSessionID(cfg *config.Config) int {
	return cfg.SessionID
}

func GetNextSessionID(cfg *config.Config) int {
	return cfg.SessionID + 1
}

func GetPreviousSessionID(cfg *config.Config) int {
	return cfg.SessionID - 1
}
