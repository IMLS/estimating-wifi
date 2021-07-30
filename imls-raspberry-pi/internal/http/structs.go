package http

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
