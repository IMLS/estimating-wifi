package main

import "testing"

func TestIsJsonOk(t *testing.T) {
	var jsonTests = []map[string]bool{
		{"{}": true},
		{"{": false},
		{`{"data": 3}`: true},
		{`{data: 3}`: false},
	}

	for ndx, pair := range jsonTests {
		for k, v := range pair {
			if isJsonOk(k) != v {
				t.Fatal("Test", ndx, ":", k, "should be", v)
			}
		}
	}
}
