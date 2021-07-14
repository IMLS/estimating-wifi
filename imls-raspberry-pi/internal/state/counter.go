package state

import (
	"fmt"
	"strconv"

	"gsa.gov/18f/config"
)

type Counter struct {
	*Queue
}

func NewCounter(cfg *config.Config, name string) *Counter {
	q := NewQueue(cfg, name)
	// DROP THE TABLE.
	q.Enqueue("0")
	return &Counter{q}
}

func (q *Queue) Increment() int {
	i := q.Dequeue()
	n, _ := strconv.Atoi(fmt.Sprintf("%v", i))
	n = n + 1
	q.Enqueue(fmt.Sprintf("%v", n))
	return n
}

func (q *Queue) Value() int {
	i := q.Peek()
	n, _ := strconv.Atoi(fmt.Sprintf("%v", i))
	return n
}

func (q *Queue) Prev() int {
	n := q.Value()
	n = n - 1
	return n
}
