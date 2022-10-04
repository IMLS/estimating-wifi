package models

// Saved searches are expressed as JSON.
// We can choose a field to search (or "ALL")
// and a regex as our query.
type Search struct {
	Field string `json:"field"`
	Query string `json:"query"`
}

type Device struct {
	Exists        bool
	Search        *Search
	Physicalid    int
	Description   string
	Businfo       string
	Logicalname   string
	Serial        string
	Mac           string
	Configuration string
	Vendor        string
}
