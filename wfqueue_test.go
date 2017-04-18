package wfqueue

import "testing"

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
	}
)

func TestEnqueue(t *testing.T) {
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

func TestDequeue(t *testing.T) {
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
