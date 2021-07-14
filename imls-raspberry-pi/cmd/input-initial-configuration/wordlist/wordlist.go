// Public domain wordlist used as source.
// https://github.com/MichaelWehar/Public-Domain-Word-Lists/blob/master/5000-more-common.txt

package wordlist

import (
	"bytes"
	"embed"
)

// We embed the entire wordlist into the executable.
// This needs to stay in sync with the web version, or the decoding
// will not work.
//go:embed wordlist.txt
var f embed.FS

// Map wordpairs to line numbers
// We use the line number as a binary value to
// decode into three, six-bit values, and those become
// the decoded key piece.
var Wordhash = make(map[string]int)

// Read the wordlist into the hash.
func Init() {
	data, _ := f.ReadFile("wordlist.txt")
	words := bytes.Split(data, []byte("\n"))
	for lineno, pair := range words {
		// Wordlist = append(Wordlist, string(pair))
		Wordhash[string(pair)] = lineno
	}
}
