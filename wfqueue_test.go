package wfqueue

import (
	"testing"
)

type operation struct {
	value   int
	enqueue bool
}

var (
	operations = []operation{
		operation{1, true},
		operation{2, true},
		operation{3, true},
		operation{-1, false},
		operation{4, true},
		operation{-1, false},
		operation{-1, false},
		operation{5, true},
		operation{-1, false},
		operation{-1, false},
	}
)

func TestWFQueueEnqueue(t *testing.T) {
	tid := 1
	q := NewWFQueue(10)

	q.Enqueue(1, tid)
	if q.Len() != 1 {
		t.Errorf("queue size is expected to be 1, but returned %d instead", q.Len())
	}

	v, ok := q.Dequeue(tid)
	if !ok {
		t.Errorf("queue dequeue failed")
	}
	if v != 1 {
		t.Errorf("queue dequeue is expected to return 1, but returned %d instead", v)
	}

	q.Enqueue(2, tid)
	if q.Len() != 1 {
		t.Errorf("queue size is expected to be 1, but returned %d instead", q.Len())
	}

	_, ok = q.Dequeue(tid)
	if !ok {
		t.Errorf("queue dequeue failed")
	}
	if q.Len() != 0 {
		t.Errorf("queue size is expected to empty, but returned size %d instead", q.Len())
	}
}

func TestWFQueueDequeue(t *testing.T) {
	tid := 1
	q := NewWFQueue(10)

	q.Enqueue(1, tid)
	q.Enqueue(2, tid)

	if q.Len() != 2 {
		t.Errorf("queue size is expected to empty, but returned size %d instead", q.Len())
	}

	v, ok := q.Dequeue(tid)
	if !ok {
		t.Errorf("queue dequeue failed")
	}
	if v != 1 {
		t.Errorf("queue dequeue is expected to return 1, but returned %d instead", v)
	}

	v, ok = q.Dequeue(tid)
	if !ok {
		t.Errorf("queue dequeue failed")
	}
	if v != 2 {
		t.Errorf("queue dequeue is expected to return 1, but returned %d instead", v)
	}

	if q.Len() != 0 {
		t.Errorf("queue size is expected to empty, but returned size %d instead", q.Len())
	}
}

func BenchmarkWFQueueEnqueue10Items(b *testing.B) {
	benchmarkWFQueueEnqueue(10, 1, b)
}

func BenchmarkWFQueueEnqueue100Items(b *testing.B) {
	benchmarkWFQueueEnqueue(100, 1, b)
}

func BenchmarkWFQueueEnqueue1000Items(b *testing.B) {
	benchmarkWFQueueEnqueue(1000, 1, b)
}

func BenchmarkWFQueueEnqueue10000Items(b *testing.B) {
	benchmarkWFQueueEnqueue(10000, 1, b)
}

func BenchmarkWFQueueEnqueue100000Items(b *testing.B) {
	benchmarkWFQueueEnqueue(100000, 1, b)
}

func benchmarkWFQueueEnqueue(numItems, numThreads int, b *testing.B) {
	qs := make([]*WFQueue, b.N)
	for i := 0; i < b.N; i++ {
		qs[i] = NewWFQueue(numThreads)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queue := qs[i]
		for j := 0; j < numItems; j++ {
			queue.Enqueue(j, 0)
		}
	}
}

func BenchmarkWFQueueDequeue10Items(b *testing.B) {
	benchmarkWFQueueDequeue(10, 1, b)
}

func BenchmarkWFQueueDequeue100Items(b *testing.B) {
	benchmarkWFQueueDequeue(100, 1, b)
}

func BenchmarkWFQueueDequeue1000Items(b *testing.B) {
	benchmarkWFQueueDequeue(1000, 1, b)
}

func BenchmarkWFQueueDequeue10000Items(b *testing.B) {
	benchmarkWFQueueDequeue(10000, 1, b)
}

func BenchmarkWFQueueDequeue100000Items(b *testing.B) {
	benchmarkWFQueueDequeue(100000, 1, b)
}

func benchmarkWFQueueDequeue(numItems, numThreads int, b *testing.B) {
	qs := make([]*WFQueue, b.N)
	for i := 0; i < b.N; i++ {
		qs[i] = NewWFQueue(numThreads)
	}
	for i := 0; i < b.N; i++ {
		queue := qs[i]
		for j := 0; j < numItems; j++ {
			queue.Enqueue(j, 0)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queue := qs[i]
		for j := 0; j < numItems; j++ {
			queue.Dequeue(0)
		}
	}
}

func BenchmarkWFQueueEnqueueDequeue(b *testing.B) {
	numItems := 1000

	qs := make([]*WFQueue, b.N)
	for i := 0; i < b.N; i++ {
		qs[i] = NewWFQueue(10)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queue := qs[i]
		for j := 0; j < numItems/len(operations); j++ {
			for _, op := range operations {
				if op.enqueue {
					queue.Enqueue(op.value, 0)
				} else {
					queue.Dequeue(0)
				}
			}
		}
	}
}

/*
func BenchmarkWFQueueConcurrentEnqueue10Items(b *testing.B) {
	benchmarkWFQueueConcurrentEnqueue(10, 2, b)
}

func BenchmarkWFQueueConcurrentEnqueue100Items(b *testing.B) {
	benchmarkWFQueueConcurrentEnqueue(100, 2, b)
}

func BenchmarkWFQueueConcurrentEnqueue1000Items(b *testing.B) {
	benchmarkWFQueueConcurrentEnqueue(1000, 2, b)
}

func BenchmarkWFQueueConcurrentEnqueue10000Items(b *testing.B) {
	benchmarkWFQueueConcurrentEnqueue(10000, 2, b)
}

func benchmarkWFQueueConcurrentEnqueue(numItems, numThreads int, b *testing.B) {
	qs := make([]*WFQueue, b.N)
	for i := 0; i < b.N; i++ {
		qs[i] = NewWFQueue(numThreads)
	}

	b.ResetTimer()
	queue := NewWFQueue(numThreads)
	total := b.N * numItems
	doneChan := make(chan bool, total)
	for i := 0; i < b.N; i++ {
		//queue := qs[i]
		go func(tid int, queue *WFQueue) {
			for j := 0; j < numItems; j++ {
				queue.Enqueue(j, tid%numThreads)
				doneChan <- true
			}
		}(i, queue)
	}

	for k := 0; k < total; k++ {
		<-doneChan
	}
}

func benchmarkWFQueueConcurrentEnqueue(numItems, numThreads int, b *testing.B) {
	qs := make([]*WFQueue, b.N)
	for i := 0; i < b.N; i++ {
		qs[i] = NewWFQueue(numThreads)
	}

	total := numItems * b.N
	doneChan := make(chan bool, total)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		//queue := NewWFQueue(numThreads)
		queue := qs[i]
		for j := 0; j < numItems; j++ {
			go func(v int, queue *WFQueue) {
				queue.Enqueue(v, v%numThreads)
				doneChan <- true
			}(j, queue)
		}
	}

	for k := 0; k < total; k++ {
		<-doneChan
	}
}

func BenchmarkWFQueueConcurrrentEnqueue10MaxThreads(b *testing.B) {
	benchmarkWFQueueConcurrentEnqueue(10, 10, b)
}

func BenchmarkWFQueueConcurrrentEnqueue100MaxThreads(b *testing.B) {
	benchmarkWFQueueConcurrentEnqueue(100, 100, b)
}

func BenchmarkWFQueueConcurrrentEnqueue1000MaxThreads(b *testing.B) {
	benchmarkWFQueueConcurrentEnqueue(100, 1000, b)
}

func BenchmarkWFQueueConcurrrentEnqueue10000MaxThreads(b *testing.B) {
	benchmarkWFQueueConcurrentEnqueue(100, 10000, b)
}
*/
