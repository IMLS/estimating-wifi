package state

import (
	"log"
	"testing"
)

var listcfg *CFG

func TestListSetup(t *testing.T) {
	listcfg = NewConfig()
	listcfg.Paths.WWW.Root = "/tmp"
	listcfg.Logging.LogLevel = "DEBUG"
	listcfg.Logging.Loggers = []string{"local:stderr"}
	listcfg.Databases.QueuesPath = "/tmp/queues.sqlite"
	listcfg.Databases.DurationsPath = "/tmp/durations.sqlite"
	InitConfig()
}

func TestList(t *testing.T) {
	log.Println(t.Name())
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
		t.Log("missing value in list")
		t.Fail()
	}
}

func TestListRemove(t *testing.T) {
	log.Println(t.Name())
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

	// listcfg.Log().Debug("list after remove ", ls.AsList())

	if !allthere {
		t.Log("missing value in list")
		t.Fail()
	}
	if redshirt {
		t.Log("failed to remove the redshirt")
		t.Fail()
	}
}
