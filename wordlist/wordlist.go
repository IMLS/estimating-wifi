// Public domain wordlist used as source.
// https://github.com/MichaelWehar/Public-Domain-Word-Lists/blob/master/5000-more-common.txt

package wordlist

import (
	"bytes"
	"embed"
)

//go:embed wordlist.txt
var f embed.FS

var Wordlist = make([]string, 0)

func Init() {
	data, _ := f.ReadFile("wordlist.txt")
	words := bytes.Split(data, []byte("\n"))
	for _, pair := range words {
		Wordlist = append(Wordlist, string(pair))
	}
}
