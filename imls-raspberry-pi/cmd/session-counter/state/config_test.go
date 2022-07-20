package state

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct {
	suite.Suite
}

func (suite *ConfigSuite) SetupTest() {
	tempDB, err := os.CreateTemp("", "config-test.ini")
	if err != nil {
		suite.Fail(err.Error())
	}
	SetConfigAtPath(tempDB.Name())
}

func (suite *ConfigSuite) AfterTest(suiteName, testName string) {
	// ensure a clean run.
	os.Remove(GetDurationsPath())
}

func (suite *ConfigSuite) TestConfigDefaults() {
	_, filename, _, _ := runtime.Caller(0)
	path := filepath.Dir(filename)
	durationsPath := filepath.Join(path, "test", "durations.sqlite")
	SetDurationsPath(durationsPath)
	var expected = []string{"local:stderr", "local:tmp", "api:directus"}
	result := GetLoggers()
	for i := 0; i < 3; i += 1 {
		if result[i] != expected[i] {
			suite.Fail("loggers were not equal")
		}
	}
	if GetDurationsDatabase().GetPath() != durationsPath {
		suite.Fail("duration path was not equal")
	}
	os.Remove(durationsPath)
}

func (suite *ConfigSuite) TestConfigWrite() {
	SetDeviceTag("a random string")
	result := GetDeviceTag()
	if result != "a random string" {
		suite.Fail("write was not reflected")
	}
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}
