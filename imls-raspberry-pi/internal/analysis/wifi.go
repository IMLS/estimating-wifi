package analysis

import "time"

type WifiEvents struct {
	Data []WifiEvent `json:"data"`
}

// https://stackoverflow.com/questions/18635671/how-to-define-multiple-name-tags-in-a-struct
type WifiEvent struct {
	ID                int       `json:"id" db:"id"`
	EventId           int       `json:"event_id" db:"event_id"`
	FCFSSeqId         string    `json:"fcfs_seq_id" db:"fcfs_seq_id"`
	DeviceTag         string    `json:"device_tag" db:"device_tag"`
	Localtime         time.Time `json:"localtime" db:"localtime"`
	Servertime        time.Time `json:"servertime" db:"servertime"`
	SessionId         string    `json:"session_id" db:"session_id"`
	ManufacturerIndex int       `json:"manufacturer_index" db:"manufacturer_index"`
	PatronIndex       int       `json:"patron_index" db:"patron_index"`
}
