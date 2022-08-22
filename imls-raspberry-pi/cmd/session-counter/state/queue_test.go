//nolint:typecheck
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

func (suite *QueueSuite) TestDequeueEmpty() {
	q := NewQueue[string]("queue1")
	removed, err := q.Dequeue()
	if removed != "" {
		suite.Fail("dequeued a value from empty queue")
	}
	if err == nil {
		suite.Fail("dequeued an empty queue")
	}
}

func (suite *QueueSuite) TestDequeueAsList() {
	q := NewQueue[int64]("queue1")
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(100)
	l := q.AsList()
	if l[0] != 1 {
		suite.Fail("incorrect first value")
	}
	if l[1] != 2 {
		suite.Fail("incorrect second value")
	}
	if l[2] != 100 {
		suite.Fail("incorrect third value")
	}
}

func (suite *QueueSuite) TestDequeueRemove() {
	q := NewQueue[int64]("queue1")
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(100)
	q.Remove(2)
	l := q.AsList()
	if l[0] != 1 {
		suite.Fail("incorrect first value")
	}
	if l[1] != 100 {
		suite.Fail("incorrect second value")
	}
}

func TestQueueSuite(t *testing.T) {
	suite.Run(t, new(QueueSuite))
}
