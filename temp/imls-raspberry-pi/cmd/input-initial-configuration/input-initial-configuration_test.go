package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"

	"gsa.gov/18f/cmd/input-initial-configuration/wordlist"
)

type test struct {
	have string
	want string
}

func TestYesOrNo(t *testing.T) {
	input := bytes.NewBufferString("no")
	if readYesOrNo(input) {
		t.Errorf("did not accept 'no'")
	}
	input = bytes.NewBufferString("Y")
	if !readYesOrNo(input) {
		t.Errorf("did not accept 'Y'")
	}
	input = bytes.NewBufferString("nah\nn")
	if readYesOrNo(input) {
		t.Errorf("accepted 'nah'")
	}
}

func TestDecode(t *testing.T) {

	tests := [...]test{
		{have: "compass passing", want: "lib"},
		{have: "period keep", want: "rar"},
		{have: "definite recover", want: "y"},
		{have: "state term", want: "2LV"},
		{have: "native harmony", want: "tzH"},
		{have: "forward metallic", want: "rVM"},
		{have: "water case", want: "C4u"},
		{have: "measure return", want: "0lR"},
		{have: "reason spiritual", want: "PDp"},
		{have: "external call", want: "Wg"},
		{have: "shoulder joint", want: "Yyl"},
		{have: "bearing uniform", want: "HLk"},
		{have: "country weather", want: "eoR"},
		{have: "form nature", want: "1HT"},
		{have: "power language", want: "3uc"},
		{have: "instrument northern", want: "tu4"},
		{have: "surface belief", want: "Jc"},
		{have: "present imperfect", want: "18f"},
		{have: "reduce from", want: "mat"},
		{have: "flat quiet", want: "tja"},
		{have: "particular phrase", want: "dud"},
		// {have: "", want: ""},
	}

	wordlist.Init()
	for _, test := range tests {
		ndx, _ := wordlist.GetPairIndex(strings.TrimSpace(test.have))
		result := decode(ndx)
		if result != test.want {
			t.Error(fmt.Sprintf("`%v`", test.have),
				"became", fmt.Sprintf("`%v`", result),
				"not", fmt.Sprintf("`%v`", test.want))
		} else {
			log.Println(test.have, "->", test.want)
		}
	}
}
