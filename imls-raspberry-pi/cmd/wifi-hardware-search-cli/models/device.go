package models

import "gsa.gov/18f/find-ralink/config"

type Device struct {
	Exists        bool
	Search        *config.Search
	Physicalid    int
	Description   string
	Businfo       string
	Logicalname   string
	Serial        string
	Mac           string
	Configuration string
	Vendor        string
	Extract       string
}
