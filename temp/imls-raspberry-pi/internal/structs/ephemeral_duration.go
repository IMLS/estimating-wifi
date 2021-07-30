package structs

type EphemeralDuration struct {
	Start int64  `db:"start" type:"INTEGER"`
	End   int64  `db:"end" type:"INTEGER"`
	MAC   string `db:"mac" type:"TEXT"`
	// SessionID int    `db:"session_id" sqlite:"INTEGER"`
}
