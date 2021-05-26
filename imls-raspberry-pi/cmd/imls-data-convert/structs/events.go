package structs

import "time"

type EventEvents struct {
	Data []EventEvent `json:"data"`
}

type EventEvent struct {
	Info       string    `json:"info"`
	ID         int       `json:"id"`
	PiSerial   string    `json:"pi_serial"`
	FCFSSeqId  string    `json:"fcfs_seq_id"`
	DeviceTag  string    `json:"device_tag"`
	SessionId  string    `json:"session_id"`
	Localtime  time.Time `json:"localtime"`
	Servertime time.Time `json:"servertime"`
	// Tag is our event tag; eg. "logging_devices"
	Tag string `json:"tag"`
}
