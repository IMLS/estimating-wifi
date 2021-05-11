package config

import (
	"crypto/sha256"
	"fmt"
	"time"
)

func CreateSessionId() string {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", time.Now())))
	sid := fmt.Sprintf("%x", h.Sum(nil))[0:16]
	return sid
}
