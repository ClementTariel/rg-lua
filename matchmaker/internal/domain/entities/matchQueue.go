package entities

import (
	"errors"
	"sync"
)

const QUEUE_SIZE = 10

type MatchQueue struct {
	pendingMatchs [QUEUE_SIZE]PendingMatch
	index         int
	length        int
	mu            sync.Mutex
}

func (q *MatchQueue) IsFull() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.length == QUEUE_SIZE
}

func (q *MatchQueue) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.length == 0
}

func (q *MatchQueue) Pop() (PendingMatch, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.length == 0 {
		return PendingMatch{}, errors.New("Empty queue")
	}
	pendingMatch := q.pendingMatchs[(q.index+QUEUE_SIZE-q.length)%QUEUE_SIZE]
	q.length -= 1
	return pendingMatch, nil
}

func (q *MatchQueue) Push(pendingMatch PendingMatch) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.length == QUEUE_SIZE {
		return false
	}
	q.length += 1
	q.pendingMatchs[q.index] = pendingMatch
	q.index = (q.index + 1) % QUEUE_SIZE
	return true
}

func NewMatchQueue() MatchQueue {
	return MatchQueue{
		pendingMatchs: [QUEUE_SIZE]PendingMatch{},
		index:         0,
		length:        0,
	}
}
