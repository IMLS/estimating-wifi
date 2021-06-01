package analysis

import "time"

type WifiEvents struct {
	Data []WifiEvent `json:"data"`
}

type WifiEvent struct {
	ID                int       `json:"id"`
	EventId           int       `json:"event_id"`
	FCFSSeqId         string    `json:"fcfs_seq_id"`
	DeviceTag         string    `json:"device_tag"`
	Localtime         time.Time `json:"localtime"`
	Servertime        time.Time `json:"servertime"`
	SessionId         string    `json:"session_id"`
	ManufacturerIndex int       `json:"manufacturer_index"`
	PatronIndex       int       `json:"patron_index"`
}
