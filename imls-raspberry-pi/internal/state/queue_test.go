package state

import (
	"log"
	"testing"
)

var cfg *CFG

func TestSetup(t *testing.T) {
	NewConfig()
	cfg = GetConfig()
	cfg.Paths.WWW.Root = "/tmp"
	cfg.Logging.LogLevel = "DEBUG"
	cfg.Logging.Loggers = []string{"local:stderr"}
	cfg.Databases.QueuesPath = "/tmp/queues.sqlite"
	cfg.Databases.DurationsPath = "/tmp/durations.sqlite"
	InitConfig()
}
func TestQueueCreate(t *testing.T) {
	log.Println(t.Name())
	q := NewQueue("queue1")
	t.Log(q)
	// t.Log("current session id " + fmt.Sprint(cfg.GetCurrentSessionId()))
}
func TestEnqueue(t *testing.T) {
	log.Println(t.Name())
	q := NewQueue("queue1")
	q.Enqueue("123")
	q.Enqueue("abc")
}

func TestMultiEnqueue(t *testing.T) {
	log.Println(t.Name())
	q := NewQueue("queue1")
	q.Enqueue("123")
	q.Enqueue("123")
}

func TestPeek(t *testing.T) {
	log.Println(t.Name())
	q := NewQueue("newqueue")
	_, err := q.Peek()
	if err == nil {
		t.Log("peek on an empty returned nil")
		t.Log(err.Error())
		t.Fail()
	}
}

func TestDequeue(t *testing.T) {
	log.Println(t.Name())
	q := NewQueue("queue1")
	shouldremove, _ := q.Peek()
	removed, err := q.Dequeue()
	if err != nil {
		t.Fatal("nothing on the queue")
	}
	if removed != shouldremove {
		cfg.Log().Error("DID NOT FIND APPROPRIATE NEXT ITEM.")
	}
}
