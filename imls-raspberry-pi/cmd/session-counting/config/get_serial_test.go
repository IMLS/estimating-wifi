package config

import (
	"regexp"
	"runtime"
	"testing"
)

func Test_GetSerial(t *testing.T) {
	// 20200505 MCJ
	// Only run this test on a Raspberry Pi
	if runtime.GOOS == "linux" && runtime.GOARCH == "arm" {
		serial := GetSerial()
		found, err := regexp.MatchString(FakeSerialCheck, serial)
		if err != nil {
			t.Log("error in GetSerial regexp check")
			t.Fail()
		}
		if found {
			t.Log("Did not properly read /proc/cpuinfo")
			t.Fail()
		}
		if len(serial) != 16 {
			t.Error("Serial is not length 16")
		}
		t.Log("serial", serial)
	}

}
