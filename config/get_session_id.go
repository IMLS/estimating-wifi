package config

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"time"

	"gsa.gov/18f/session-counter/constants"
)

func CreateSessionId() string {
	h := sha256.New()
	email := os.Getenv(constants.AuthEmailKey)
	// FIXME: Use the email instead of the token.
	// Guaranteed to be unique. Current time along with our auth token, hashed.
	h.Write([]byte(fmt.Sprintf("%v%x", time.Now(), email)))
	sid := fmt.Sprintf("%x", h.Sum(nil))[0:16]
	// Keep it short.
	log.Println("Session id: ", sid)
	return sid
}
