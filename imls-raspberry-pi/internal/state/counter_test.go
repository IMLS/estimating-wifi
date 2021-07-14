package state

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	"gsa.gov/18f/config"
	"gsa.gov/18f/logwrapper"
)

type CounterSuite struct {
	suite.Suite
	cfg *config.Config
	lw  *logwrapper.StandardLogger
	c   *Counter
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *CounterSuite) SetupTest() {
	cfg := config.NewConfig()
	suite.cfg = cfg
	suite.cfg.LogLevel = "INFO"
	dir, _ := os.Getwd()
	cfg.Local.WebDirectory = dir
	suite.lw = logwrapper.NewLogger(cfg)
	suite.c = NewCounter(suite.cfg, "a")
}

func (suite *CounterSuite) TestCounter() {
	c := GetCounter(suite.cfg, "a")
	log.Println("TestCounter initial value", c.Value())
	if c.Value() != 0 {
		suite.Fail("counter was not at zero")
	}
}

func (suite *CounterSuite) TestCounter2() {
	// Using the "suite pointer"
	log.Println("TestCounter initial value", suite.c.Value())
	if suite.c.Value() != 0 {
		suite.Fail("counter was not at zero")
	}
}

func (suite *CounterSuite) TestIncCounter() {
	c := GetCounter(suite.cfg, "a")
	c.Increment()
	log.Println("TestIncCounter incremented value", c.Value())

	if c.Value() != 1 {
		suite.Fail("counter was not incremented")
	}
}

func (suite *CounterSuite) TestResetCounter() {
	c := GetCounter(suite.cfg, "a")
	c.Increment()
	log.Println("TestResetCounter incremented value", c.Value())

	if c.Value() != 2 {
		suite.Fail("counter was not incremented after reset")
	}
	c.Reset()

	if c.Value() != 0 {
		suite.Fail("counter was not reset")
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(CounterSuite))
}
