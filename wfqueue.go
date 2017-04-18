package wfqueue

import "sync/atomic"
import "strconv"

// Node represents a Node in Wait Free Queue
type Node struct {
	value  int
	next   *atomic.Value
	enqTid *atomic.Value
	deqTid *atomic.Value
}

// NewNode initializes a new pointer to node struct
func NewNode(value, enqTid int) *Node {
	atomicEnqTid := new(atomic.Value)
	atomicEnqTid.Store(enqTid)
	atomicDeqTid := new(atomic.Value)
	atomicDeqTid.Store(-1)

	return &Node{
		value:  value,
		next:   new(atomic.Value),
		enqTid: atomicEnqTid,
		deqTid: atomicDeqTid,
	}
}

// OpDesc represents the state of a thread operation on the queue (either enqueue or dequeue)
type OpDesc struct {
	phase   int
	pending bool
	enqueue bool
	node    *Node
}

// NewOpDesc initializes a new pointer to OpDesc struct
func NewOpDesc(phase int, pending, enqueue bool, node *Node) *OpDesc {
	return &OpDesc{
		phase:   phase,
		pending: pending,
		enqueue: enqueue,
		node:    node,
	}
}

// WFQueue represents a wait free queue
type WFQueue struct {
	head       *atomic.Value
	tail       *atomic.Value
	state      []*atomic.Value
	numThreads int
}

// NewWFQueue takes in number of threads and initializes a new pointer to WFQueue struct
func NewWFQueue(numThreads int) *WFQueue {
	// make sentinal head and tail
	sentinal := NewNode(-1, -1)
	//var head, tail *atomic.Value
	head := new(atomic.Value)
	tail := new(atomic.Value)
	head.Store(sentinal)
	tail.Store(sentinal)

	stateSlice := make([]*atomic.Value, numThreads)
	// initialize the initial states
	for i := 0; i < numThreads; i++ {
		stateSlice[i] = new(atomic.Value)
		stateSlice[i].Store(NewOpDesc(-1, false, true, nil))
	}
	//var state *atomic.Value
	state := new(atomic.Value)
	state.Store(stateSlice)

	return &WFQueue{
		head:       head,
		tail:       tail,
		state:      stateSlice,
		numThreads: numThreads,
	}
}

// Enqueue enqueus new value in WFQueue
func (queue *WFQueue) Enqueue(value, tid int) {
	phase := queue.maxPhase()
	queue.state[tid].Store(NewOpDesc(phase, true, true, NewNode(value, tid)))
	queue.help(tid)

}

// Dequeue dequeues value from WFQueue
func (queue *WFQueue) Dequeue(tid int) (int, bool) {
	phase := queue.maxPhase() + 1
	queue.state[tid].Store(NewOpDesc(phase, true, false, nil))
	queue.help(phase)
	queue.helpDequeueFinish()

	node := queue.state[tid].Load().(*OpDesc).node
	if node == nil {
		return -1, false
	}

	return node.next.Load().(*Node).value, true
}

// Len returns the number of items in the WFQueue
func (queue *WFQueue) Len() int {
	head := queue.head.Load()
	num := 0
	for {
		if head == nil {
			break
		}
		node := head.(*Node)
		if node.value != -1 {
			num++
		}
		head = node.next.Load()
	}

	if num > 0 {
		return num - 1
	}
	return num
}

// String returns the string representation of WFQueue
func (queue *WFQueue) String() string {
	head := queue.head.Load()
	res := ""
	for {
		if head == nil {
			break
		}
		node := head.(*Node)
		res += strconv.Itoa(node.value) + " "
		head = node.next.Load()
	}

	return res
}

// -------------------------------------------------
//
// Helper functions
//
// -------------------------------------------------

func (queue *WFQueue) help(phase int) {
	for i := 0; i < len(queue.state); i++ {
		desc := queue.state[i].Load().(*OpDesc)
		if desc.pending && desc.phase <= phase {
			if desc.enqueue {
				queue.helpEnqueue(i, phase)
			} else {
				queue.helpDequeue(i, phase)
			}
		}
	}
}

func (queue *WFQueue) helpEnqueue(tid, phase int) {
	for queue.isStillPending(tid, phase) {
		last := queue.tail.Load().(*Node)
		nextI := last.next.Load()

		var next *Node
		if nextI != nil {
			next = nextI.(*Node)
		}

		if last == queue.tail.Load().(*Node) {
			if next == nil {
				if queue.isStillPending(tid, phase) {
					if compareAndSwapNode(last.next, next, queue.state[tid].Load().(*OpDesc).node) {
						queue.helpEnqueueFinish()
						return
					}
				}
			}
		} else {
			queue.helpEnqueueFinish()
		}
	}
}

func (queue *WFQueue) helpEnqueueFinish() {
	last := queue.tail.Load().(*Node)
	nextI := last.next.Load()

	var next *Node
	if nextI != nil {
		next = nextI.(*Node)
	}
	if next != nil {
		tid := next.enqTid.Load().(int)
		curDesc := queue.state[tid].Load().(*OpDesc)
		if last == queue.tail.Load().(*Node) && queue.state[tid].Load().(*OpDesc).node == next {
			newDesc := NewOpDesc(queue.state[tid].Load().(*OpDesc).phase, false, true, next)
			//queue.state[tid].Store(newDesc)
			compareAndSwapOpDesc(queue.state[tid], curDesc, newDesc)
			compareAndSwapNode(queue.tail, last, next)
		}
	}
}

func (queue *WFQueue) helpDequeue(tid, phase int) {
	for queue.isStillPending(tid, phase) {
		first := queue.head.Load().(*Node)
		last := queue.tail.Load().(*Node)
		nextI := first.next.Load()

		var next *Node
		if nextI != nil {
			next = nextI.(*Node)
		}

		if first == queue.head.Load().(*Node) {
			if first == last {
				if next == nil {
					curDesc := queue.state[tid].Load().(*OpDesc)
					if last == queue.tail.Load().(*Node) && queue.isStillPending(tid, phase) {
						newDesc := NewOpDesc(queue.state[tid].Load().(*OpDesc).phase, false, false, nil)
						compareAndSwapOpDesc(queue.state[tid], curDesc, newDesc)
					}
				} else {
					queue.helpEnqueueFinish()
				}
			} else {
				curDesc := queue.state[tid].Load().(*OpDesc)
				node := curDesc.node
				if !queue.isStillPending(tid, phase) {
					break
				}
				if first == queue.head.Load().(*Node) && node != first {
					newDesc := NewOpDesc(queue.state[tid].Load().(*OpDesc).phase, true, false, first)
					if !compareAndSwapOpDesc(queue.state[tid], curDesc, newDesc) {
						continue
					}
				}

				first.deqTid.Store(tid)
				queue.helpDequeueFinish()
			}
		}
	}
}

func (queue *WFQueue) helpDequeueFinish() {
	first := queue.head.Load().(*Node)
	nextI := first.next.Load()

	var next *Node
	if nextI != nil {
		next = nextI.(*Node)
	}

	tid := first.deqTid.Load().(int)
	if tid != -1 {
		curDesc := queue.state[tid].Load().(*OpDesc)
		if first == queue.head.Load().(*Node) && next != nil {
			newDesc := NewOpDesc(queue.state[tid].Load().(*OpDesc).phase, false, false, queue.state[tid].Load().(*OpDesc).node)
			compareAndSwapOpDesc(queue.state[tid], curDesc, newDesc)
			compareAndSwapNode(queue.head, first, next)
		}
	}
}

func (queue *WFQueue) maxPhase() int {
	maxPhase := -1
	for i := 0; i < queue.numThreads; i++ {
		phase := queue.state[i].Load().(*OpDesc).phase
		if phase > maxPhase {
			maxPhase = phase
		}
	}
	return maxPhase
}

func (queue *WFQueue) isStillPending(tid, phase int) bool {
	return queue.state[tid].Load().(*OpDesc).pending && queue.state[tid].Load().(*OpDesc).phase <= phase
}