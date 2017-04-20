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

func BenchmarkLFQueueEnqueue(b *testing.B) {
	numItems := 1000

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

func BenchmarkLFQueueDequeue(b *testing.B) {
	numItems := 1000

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
