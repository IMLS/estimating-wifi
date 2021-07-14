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
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *CounterSuite) SetupTest() {
	cfg := config.NewConfig()
	suite.cfg = cfg
	suite.cfg.LogLevel = "DEBUG"
	dir, _ := os.Getwd()
	cfg.Local.WebDirectory = dir
	suite.lw = logwrapper.NewLogger(cfg)
}

func (suite *CounterSuite) TestCounter() {
	c := NewCounter(suite.cfg, "a")
	log.Println(c.Value())
	if c.Peek() != 0 {
		suite.Fail("counter was not at zero")
	}
}

func (suite *CounterSuite) TestIncCounter() {
	c := NewCounter(suite.cfg, "a")
	c.Increment()
	log.Println(c.Value())

	if c.Peek() != 1 {
		suite.Fail("counter was not incremented")
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(CounterSuite))
}
