package state

func InitializeSessionID() int64 {
	return GetClock().Now().Unix()
}

func (dc *databaseConfig) GetCurrentSessionID() int64 {
	return dc.sessionID
}

func (dc *databaseConfig) IncrementSessionID() int64 {
	dc.sessionID = InitializeSessionID()
	return dc.sessionID
}
