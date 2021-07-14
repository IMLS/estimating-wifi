package state

import (
	"fmt"
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

var s *CounterSuite

// // Make sure that VariableThatShouldStartAtFive is set to five
// // before each test
// func (suite *CounterSuite) SetupTest() {

// }

func (suite *CounterSuite) TestCounter() {
	log.Println("TestCounter initial value", suite.c.Value())
	if suite.c.Value() != 0 {
		suite.Fail("counter was not at zero")
	}
}

func (suite *CounterSuite) TestIncCounter() {
	suite.c.Reset()
	suite.c.Increment()
	log.Println("TestIncCounter incremented value", suite.c.Value())

	if suite.c.Value() != 1 {
		suite.Fail("counter was not incremented")
	}

	log.Println("Value after increment", suite.c.Value())

	log.Println("Check exists", suite.c.db.CheckTableExists(suite.c.name))
}

func (suite *CounterSuite) TestResetCounter() {
	log.Println("TestResetCounter initial value", suite.c.Value())
	for i := 0; i < 100; i++ {
		suite.c.Increment()
	}

	if suite.c.Value() != 101 {
		suite.Fail(fmt.Sprintf("counter is not 101; it is %d", suite.c.Value()))
	}

	suite.c.Reset()

	if suite.c.Value() != 0 {
		suite.Fail("counter was not reset")
	}

}

func TestSuite(t *testing.T) {
	cfg := config.NewConfig()
	s := &CounterSuite{}
	s.cfg = cfg
	s.cfg.LogLevel = "DEBUG"
	dir, _ := os.Getwd()
	cfg.Local.WebDirectory = dir
	s.lw = logwrapper.NewLogger(cfg)
	s.c = NewCounter(s.cfg, "a")
	suite.Run(t, s)
}
