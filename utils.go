package wfqueue

import (
	"sync/atomic"
)

func compareAndSwapNode(dest *atomic.Value, old *Node, new *Node) bool {
	if dest.Load() == nil || dest.Load().(*Node) == old {
		dest.Store(new)
		return true
	}
	return false
}

func compareAndSwapOpDesc(dest *atomic.Value, old *OpDesc, new *OpDesc) bool {
	if dest.Load() == nil || dest.Load().(*OpDesc) == old {
		dest.Store(new)
		return true
	}
	return false
}
