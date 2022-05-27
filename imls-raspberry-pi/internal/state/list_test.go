package state

import (
	"os"

	"github.com/stretchr/testify/suite"
)

type ListSuite struct {
	suite.Suite
}

func (suite *ListSuite) SetupTest() {
	temp, err := os.CreateTemp("", "list-test.ini")
	if err != nil {
		suite.Fail(err.Error())
	}
	SetConfigAtPath(temp.Name())
}

func (suite *ListSuite) AfterTest(suiteName, testName string) {
	// ensure a clean run.
	os.Remove(GetDurationsPath())
}

func (suite *ListSuite) TestList() {
	ls := NewList("ls1")
	ls.Push("hello")
	ls.Push("goodbye")
	asls := ls.AsList()
	shouldhave := []string{"hello", "goodbye"}
	allthere := true
	for _, s := range shouldhave {
		found := false
		for _, is := range asls {
			if is == s {
				found = true
			}
		}
		allthere = allthere || found
	}
	if !allthere {
		suite.Fail("missing value in list")
	}
}

func (suite *ListSuite) TestListRemove() {
	ls := NewList("ls1")
	ls.Push("hello")
	ls.Push("redshirt")
	ls.Push("goodbye")

	shouldhave := []string{"hello", "goodbye"}
	ls.Remove("redshirt")
	asls := ls.AsList()

	allthere := true
	redshirt := false
	for _, s := range shouldhave {
		found := false
		for _, is := range asls {
			if is == s {
				found = true
			}
			if is == "redshirt" {
				redshirt = true
			}
		}
		allthere = allthere || found
	}

	if !allthere {
		suite.Fail("missing value in list")
	}
	if redshirt {
		suite.Fail("failed to remove the redshirt")
	}
}
