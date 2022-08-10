package config

import (
	"os"
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

func (suite *ConfigSuite) TestConfigDefaults() {
	var expected = []string{"local:stderr", "local:tmp", "api:directus"}
	result := GetLoggers()
	for i := 0; i < 3; i += 1 {
		if result[i] != expected[i] {
			suite.Fail("loggers were not equal")
		}
	}
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
