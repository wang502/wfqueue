package wfqueue

import "testing"

func TestLFQueueEnqueue(t *testing.T) {
	q := NewLFQueue()

	q.Enqueue(1)

	v, ok := q.Dequeue()
	if !ok {
		t.Errorf("queue dequeue failed")
	}
	if v != 1 {
		t.Errorf("queue dequeue is expected to return 1, but returned %d instead", v)
	}

	q.Enqueue(2)

	_, ok = q.Dequeue()
	if !ok {
		t.Errorf("queue dequeue failed")
	}
}

func TestLFQueueDequeue(t *testing.T) {
	q := NewLFQueue()

	q.Enqueue(1)
	q.Enqueue(2)

	v, ok := q.Dequeue()
	if !ok {
		t.Errorf("queue dequeue failed")
	}
	if v != 1 {
		t.Errorf("queue dequeue is expected to return 1, but returned %d instead", v)
	}

	v, ok = q.Dequeue()
	if !ok {
		t.Errorf("queue dequeue failed")
	}
	if v != 2 {
		t.Errorf("queue dequeue is expected to return 1, but returned %d instead", v)
	}

}

func BenchmarkLFQueueEnqueue10Items(b *testing.B) {
	benchmarkLFQueueEnqueue(10, 1, b)
}

func BenchmarkLFQueueEnqueue100Items(b *testing.B) {
	benchmarkLFQueueEnqueue(100, 1, b)
}

func BenchmarkLFQueueEnqueue1000Items(b *testing.B) {
	benchmarkLFQueueEnqueue(1000, 1, b)
}

func BenchmarkLFQueueEnqueue10000Items(b *testing.B) {
	benchmarkLFQueueEnqueue(10000, 1, b)
}

func BenchmarkLFQueueEnqueue100000Items(b *testing.B) {
	benchmarkLFQueueEnqueue(100000, 1, b)
}

func benchmarkLFQueueEnqueue(numItems, numThreads int, b *testing.B) {
	qs := make([]*LFQueue, b.N)
	for i := 0; i < b.N; i++ {
		qs[i] = NewLFQueue()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queue := qs[i]
		for j := 0; j < numItems; j++ {
			queue.Enqueue(j)
		}
	}
}

func BenchmarkLFQueueDequeue10Items(b *testing.B) {
	benchmarkLFQueueDequeue(10, 1, b)
}

func BenchmarkLFQueueDequeue100Items(b *testing.B) {
	benchmarkLFQueueDequeue(100, 1, b)
}

func BenchmarkLFQueueDequeue1000Items(b *testing.B) {
	benchmarkLFQueueDequeue(1000, 1, b)
}

func BenchmarkLFQueueDequeue10000Items(b *testing.B) {
	benchmarkLFQueueDequeue(10000, 1, b)
}

func BenchmarkLFQueueDequeue100000Items(b *testing.B) {
	benchmarkLFQueueDequeue(100000, 1, b)
}

func benchmarkLFQueueDequeue(numItems, numThreads int, b *testing.B) {
	qs := make([]*LFQueue, b.N)
	for i := 0; i < b.N; i++ {
		qs[i] = NewLFQueue()
	}
	for i := 0; i < b.N; i++ {
		queue := qs[i]
		for j := 0; j < numItems; j++ {
			queue.Enqueue(j)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queue := qs[i]
		for j := 0; j < numItems; j++ {
			queue.Dequeue()
		}
	}
}

func BenchmarkLFQueueEnqueueDequeue(b *testing.B) {
	numItems := 1000

	qs := make([]*LFQueue, b.N)
	for i := 0; i < b.N; i++ {
		qs[i] = NewLFQueue()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		queue := qs[i]
		for j := 0; j < numItems/len(operations); j++ {
			for _, op := range operations {
				if op.enqueue {
					queue.Enqueue(op.value)
				} else {
					queue.Dequeue()
				}
			}
		}
	}
}

func BenchmarkLFQueueConcurrentEnqueue10Items(b *testing.B) {
	benchmarkLFQueueConcurrentEnqueue(10, b)
}

func BenchmarkLFQueueConcurrentEnqueue100Items(b *testing.B) {
	benchmarkLFQueueConcurrentEnqueue(100, b)
}

func BenchmarkLFQueueConcurrentEnqueue1000Items(b *testing.B) {
	benchmarkLFQueueConcurrentEnqueue(1000, b)
}

func BenchmarkLFQueueConcurrentEnqueue10000Items(b *testing.B) {
	benchmarkLFQueueConcurrentEnqueue(10000, b)
}

func BenchmarkLFQueueConcurrentEnqueue100000Items(b *testing.B) {
	benchmarkLFQueueConcurrentEnqueue(100000, b)
}

func benchmarkLFQueueConcurrentEnqueue(numItems int, b *testing.B) {
	qs := make([]*LFQueue, b.N)
	for i := 0; i < b.N; i++ {
		qs[i] = NewLFQueue()
	}

	b.ResetTimer()

	total := b.N * numItems
	doneChan := make(chan bool, total)

	for i := 0; i < b.N; i++ {
		//queue := NewLFQueue()
		queue := qs[i]
		for j := 0; j < numItems; j++ {
			go func(j int, queue *LFQueue) {
				queue.Enqueue(j)
				doneChan <- true
			}(j, queue)
		}
	}

	for j := 0; j < total; j++ {
		<-doneChan
	}
}
