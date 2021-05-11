package http

import "gsa.gov/18f/config"

// We need some state.
type EventLogger struct {
	Cfg *config.Config
}

// {
// 	"tables": [
// 	  {
// 		"headers": [
// 		  "event_id",
// 		  "device_uuid",
// 		  "lib_user",
// 		  "localtime",
// 		  "servertime",
// 		  "session_id",
// 		  "device_id"
// 		], // End headers
// 		"whole_table_errors": [],
// 		"rows": [
// 		  {
// 			"row_number": 2,
// 			"errors": [],
// 			"data": {
// 			  "event_id": "-1",
// 			  "device_uuid": "1000000089bbf88b",
// 			  "lib_user": "matthew.jadud@gsa.gov",
// 			  "localtime": "2021-04-02T10:46:53-04:00",
// 			  "servertime": "2021-04-02T10:46:53-04:00",
// 			  "session_id": "9475068c05fea81f",
// 			  "device_id": "unknown:6"
// 			} //end data
// 		  } // end row
// 		], // end rows
// 		"valid_row_count": 1,
// 		"invalid_row_count": 0
// 	  } // end table
// 	],
// 	"valid": true
//   }
type RevalResponse struct {
	Tables []struct {
		Headers          []string      `json:"headers"`
		WholeTableErrors []interface{} `json:"whole_table_errors"`
		Rows             []struct {
			RowNumber int               `json:"row_number"`
			Errors    []interface{}     `json:"errors"`
			Data      map[string]string `json:"data"`
		} `json:"rows"`
		ValidRowCount   int `json:"valid_row_count"`
		InvalidRowCount int `json:"invalid_row_count"`
	} `json:"tables"`
	Valid bool `json:"valid"`
}
