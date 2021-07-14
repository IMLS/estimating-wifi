package tempdb

import (
	"log"
	"testing"

	"gsa.gov/18f/internal/config"
	"gsa.gov/18f/internal/logwrapper"
)

var cfg *config.Config = nil
var lw *logwrapper.StandardLogger = nil

func TestSetup(t *testing.T) {
	cfg = config.NewConfig()
	cfg.Local.WebDirectory = "/tmp"
	lw = logwrapper.NewLogger(cfg)
	lw.SetLogLevel("DEBUG")
}
func TestQueueCreate(t *testing.T) {
	log.Println(t.Name())
	q := NewQueue(cfg, "queue1")
	lw.Debug("DUMPING")
	q.db.DebugDump("queue1")

}
func TestEnqueue(t *testing.T) {
	log.Println(t.Name())
	q := NewQueue(cfg, "queue1")
	q.Enqueue("123")
	q.Enqueue("abc")
}

func TestMultiEnqueue(t *testing.T) {
	log.Println(t.Name())
	q := NewQueue(cfg, "queue1")
	q.Enqueue("123")
	q.Enqueue("123")
}

func TestPeek(t *testing.T) {
	log.Println(t.Name())
	q := NewQueue(cfg, "newqueue")
	v := q.Peek()
	if v != nil {
		t.Log("peek on an empty queue did not return nil")
		t.Fail()
	}
}

func TestDequeue(t *testing.T) {
	log.Println(t.Name())
	q := NewQueue(cfg, "queue1")
	shouldremove := q.Peek()
	removed := q.Dequeue()
	if removed != shouldremove {
		lw.Error("DID NOT FIND APPROPRIATE NEXT ITEM.")
	}
}

func TestList(t *testing.T) {
	log.Println(t.Name())
	ls := NewList(cfg, "ls1")
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
	ls := NewList(cfg, "ls1")
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

	lw.Debug("list after remove ", ls.AsList())

	if !allthere {
		t.Log("missing value in list")
		t.Fail()
	}
	if redshirt {
		t.Log("failed to remove the redshirt")
		t.Fail()
	}
}
