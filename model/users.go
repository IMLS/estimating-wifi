package model

// We want to take raw MAC addresses and map them to user ids.
// The user id will be a counter. It starts at 0 every time we restart.
// So, it is not a "tracking" id. Or, it provides an ephemeral notion of
// tracking a single user within our "uniqueness window."
//
// This means we must map mac -> id.
// This is one step in a pipeline.
type UserMapping struct {
	Mfg string
	Id  int
}
