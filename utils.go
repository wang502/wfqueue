package wfqueue

import (
	"sync/atomic"
)

func compareAndSwap(dest *atomic.Value, old interface{}, new interface{}) bool {
	if dest.Load() == nil || dest.Load() == old {
		dest.Store(new)
		return true
	}
	return false
}
