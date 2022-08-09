package api

import (
	"encoding/json"
	"log"
	"testing"
)

func TestRevalResponseUnmarshall(t *testing.T) {
	testString := `{
		"tables": [
		  {
			"headers": [
			  "event_id",
			  "device_uuid",
			  "lib_user",
			  "localtime",
			  "servertime",
			  "session_id",
			  "device_id"
			],
			"whole_table_errors": [],
			"rows": [
			  {
				"row_number": 2,
				"errors": [],
				"data": {
				  "event_id": "-1",
				  "device_uuid": "1000000089bbf88b",
				  "lib_user": "matthew.jadud@gsa.gov",
				  "localtime": "2021-04-02T10:46:53-04:00",
				  "servertime": "2021-04-02T10:46:53-04:00",
				  "session_id": "9475068c05fea81f",
				  "device_id": "unknown:6"
				}
			  }
			],
			"valid_row_count": 1,
			"invalid_row_count": 0
		  }
		],
		"valid": true
	  }`

	var rev RevalResponse
	err := json.Unmarshal([]byte(testString), &rev)
	if err != nil {
		log.Println("unmarshalling error:", err)

	} else {
		log.Println(rev)

	}
}
