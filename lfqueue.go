package wfqueue

import (
	"strconv"
	"sync/atomic"
)

// QueueNode represents a queue node in LFQueue
type QueueNode struct {
	value int
	next  *atomic.Value
}

// LFQueue represents a lock free queue for benchmark testing with WFQueue
type LFQueue struct {
	head *atomic.Value
	tail *atomic.Value
}

// NewLFQueue initializes a new LFQueue
func NewLFQueue() *LFQueue {
	sentinal := &QueueNode{
		value: -1,
		next:  new(atomic.Value),
	}

	head := new(atomic.Value)
	head.Store(sentinal)
	tail := new(atomic.Value)
	tail.Store(sentinal)
	return &LFQueue{
		head: head,
		tail: tail,
	}
}

// Enqueue put items into LFQueue
func (queue *LFQueue) Enqueue(val int) {
	newNode := &QueueNode{
		value: val,
		next:  new(atomic.Value),
	}

	for {
		//fmt.Printf("enqueueing: %d", val)
		last := queue.tail.Load().(*QueueNode)
		nextI := queue.tail.Load().(*QueueNode).next.Load()
		var next *QueueNode
		if nextI != nil {
			next = nextI.(*QueueNode)
		}
		if last == queue.tail.Load().(*QueueNode) {
			if next == nil {
				if compareAndSwapQueueNode(last.next, nil, newNode) {
					compareAndSwapQueueNode(queue.tail, last, newNode)
					return
				}
			}
		} else {
			compareAndSwapQueueNode(queue.tail, last, next)
			return
		}
	}
}

// Dequeue dequeues item for LFQueue
func (queue *LFQueue) Dequeue() (int, bool) {
	for {
		first := queue.head.Load().(*QueueNode)
		last := queue.tail.Load().(*QueueNode)
		nextI := first.next.Load()

		var next *QueueNode
		if nextI != nil {
			next = nextI.(*QueueNode)
		}

		if first == queue.head.Load().(*QueueNode) {
			if first != last && next != nil {
				if compareAndSwapQueueNode(queue.head, first, next) {
					return next.value, true
				}
			} else {
				return -1, false
			}
		}
	}
}

// String retruns string representation of LFQueue
func (queue *LFQueue) String() string {
	head := queue.head.Load()
	res := ""
	for {
		if head == nil {
			break
		}
		node := head.(*QueueNode)
		res += strconv.Itoa(node.value) + " "
		head = node.next.Load()
	}

	return res
}
