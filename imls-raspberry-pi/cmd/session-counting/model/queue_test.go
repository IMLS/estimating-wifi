package model

import (
	"testing"

	"gsa.gov/18f/config"
	"gsa.gov/18f/logwrapper"
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
	q := NewQueue(cfg, "queue1")
	lw.Debug("DUMPING")
	q.db.DebugDump("queue1")

}
func TestEnqueue(t *testing.T) {
	q := NewQueue(cfg, "queue1")
	q.Enqueue("123")
	q.Enqueue("abc")
}

func TestPeek(t *testing.T) {
	q := NewQueue(cfg, "queue1")
	lw.Debug("PEEK ", q.Peek())
}

func TestDequeue(t *testing.T) {
	q := NewQueue(cfg, "queue1")
	shouldremove := q.Peek()
	removed := q.Dequeue()
	if removed != shouldremove {
		lw.Error("DID NOT FIND APPROPRIATE NEXT ITEM.")
	}
	q.db.DebugDump("queue1")

}
