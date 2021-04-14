package wordlist

import (
	"fmt"
)

// These are thin wrappers around the hash.
// They are legacy from when it was an array instead of a hash.
// They're nice names, though, so we'll keep the interface,
// since the other code relies on them.

// PURPOSE
// A helper to see if a wordpair is in the hash.
func contains(wp string) bool {
	_, found := Wordhash[wp]
	return found
}

// PURPOSE
// Checks to see if the value is in the map, and
// if so, returns the index.
func GetPairIndex(wp string) (int, error) {
	if contains(wp) {
		return Wordhash[wp], nil
	} else {
		return -1, fmt.Errorf("wordpair [%v] not found", wp)
	}
}
