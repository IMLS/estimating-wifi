package config

import (
	"regexp"
	"testing"
)

func Test_GetSerial(t *testing.T) {
	serial := GetSerial()
	found, err := regexp.MatchString(`DENT`, serial)
	if err != nil {
		t.Log("error in GetSerial regexp check")
		t.Fail()
	}
	if found {
		t.Log("Did not properly read /proc/cpuinfo")
		t.Fail()
	}

	t.Log("serial", serial)
}
