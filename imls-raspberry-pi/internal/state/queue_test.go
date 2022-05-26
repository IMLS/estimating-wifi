package state

import (
	"os"

	"github.com/stretchr/testify/suite"
)

type QueueSuite struct {
	suite.Suite
}

func (suite *QueueSuite) SetupTest() {
	temp, err := os.CreateTemp("", "queue-test.ini")
	if err != nil {
		suite.Fail(err.Error())
	}
	SetConfigAtPath(temp.Name())
}

func (suite *QueueSuite) AfterTest(suiteName, testName string) {
	os.Remove(GetDurationsPath())
}

func (suite *QueueSuite) TestQueueCreate() {
	NewQueue("queue1")
}

func (suite *QueueSuite) TestEnqueue() {
	q := NewQueue("queue1")
	q.Enqueue("123")
	q.Enqueue("abc")
}

func (suite *QueueSuite) TestMultiEnqueue() {
	q := NewQueue("queue1")
	q.Enqueue("123")
	q.Enqueue("123")
}

func (suite *QueueSuite) TestPeek() {
	q := NewQueue("newqueue")
	_, err := q.Peek()
	if err == nil {
		suite.Fail("peek on an empty returned nil")
	}
}

func (suite *QueueSuite) TestDequeue() {
	q := NewQueue("queue1")
	shouldremove, _ := q.Peek()
	removed, err := q.Dequeue()
	if err != nil {
		suite.Fail("nothing on the queue")
	}
	if removed != shouldremove {
		suite.Fail("did not find appropriate next item.")
	}
}
