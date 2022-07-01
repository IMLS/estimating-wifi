package model

import "errors"

type storable interface {
	int64 | string
}

type Queue[S storable] struct {
	name string
	fifo []S
}

func NewQueue[S storable](name string) (q *Queue[S]) {
	q = &Queue[S]{name: name, fifo: make([]S, 0)}
	return q
}

func (queue *Queue[S]) Enqueue(item S) {
	queue.fifo = append(queue.fifo, item)
}

func (queue *Queue[S]) Peek() (S, error) {
	if len(queue.fifo) > 0 {
		return queue.fifo[0], nil
	} else {
		// https://stackoverflow.com/questions/70585852/return-default-value-for-generic-type
		// By declaring an empty variable, we get the default "nil" return
		// type for a generic.
		var result S
		return result, errors.New("queue is empty in peek")
	}

}

func (queue *Queue[S]) Length() int {
	return len(queue.fifo)
}

func (queue *Queue[S]) Dequeue() (S, error) {
	if len(queue.fifo) > 0 {
		var s S
		s, queue.fifo = queue.fifo[0], queue.fifo[1:]
		return s, nil
	} else {
		var result S
		return result, errors.New("queue is empty in dequeue")
	}

}

func (q *Queue[S]) AsList() []S {
	return q.fifo
}

func (q *Queue[S]) Remove(s S) {
	n := make([]S, 0)
	for _, v := range q.fifo {
		if v != s {
			n = append(n, v)
		}
	}
	q.fifo = n
}
