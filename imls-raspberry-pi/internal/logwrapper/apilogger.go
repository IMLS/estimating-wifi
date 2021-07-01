package logwrapper

import (
	"fmt"
)

type ApiLogger struct {
	endpoint string
}

// type Writer interface {
//     Write(p []byte) (n int, err error)
// }

func NewApiLogger(endpoint string) (api *ApiLogger) {
	api = &ApiLogger{endpoint: endpoint}
	return api
}

func (api *ApiLogger) Write(p []byte) (n int, err error) {
	fmt.Printf("API: %v\n", string(p))
	return len(p), nil
}
