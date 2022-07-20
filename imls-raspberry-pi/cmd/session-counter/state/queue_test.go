package state

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type QueueSuite struct {
	suite.Suite
}

func (suite *QueueSuite) SetupTest() {
}

func (suite *QueueSuite) AfterTest(suiteName, testName string) {
}

func (suite *QueueSuite) TestQueueCreate() {
	NewQueue[string]("queue1")
}

func (suite *QueueSuite) TestEnqueue() {
	q := NewQueue[string]("queue1")
	q.Enqueue("123")
	q.Enqueue("abc")
}

func (suite *QueueSuite) TestMultiEnqueue() {
	q := NewQueue[string]("queue1")
	q.Enqueue("123")
	q.Enqueue("123")
}

func (suite *QueueSuite) TestQueueLength() {
	q := NewQueue[string]("queue1")
	q.Enqueue("123")
	q.Enqueue("123")
	if q.Length() != 2 {
		suite.Fail("queue is the wrong length")
	}
}

func (suite *QueueSuite) TestPeek() {
	q := NewQueue[string]("newqueue")
	_, err := q.Peek()
	if err == nil {
		suite.Fail("peek on an empty returned nil")
	}
}

func (suite *QueueSuite) TestDequeue() {
	q := NewQueue[string]("queue1")
	q.Enqueue("abc")
	shouldremove, _ := q.Peek()
	removed, err := q.Dequeue()
	if err != nil {
		suite.Fail("nothing on the queue")
	}
	if removed != shouldremove {
		suite.Fail("did not find appropriate next item.")
	}
}

func TestQueueSuite(t *testing.T) {
	suite.Run(t, new(QueueSuite))
}
