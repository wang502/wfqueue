package main

import (
	"fmt"

	"github.com/wang502/wfqueue"
)

func main() {
	numThreads := 10
	queue := wfqueue.NewLFQueue()

	done := make(chan bool, numThreads)
	for i := 0; i < 10; i++ {
		go func(val, tid int) {
			queue.Enqueue(val)
			done <- true
		}(i, i)
	}

	for i := 0; i < numThreads; i++ {
		fmt.Println(<-done)
	}
	close(done)
	fmt.Printf("queue: %s\n", queue)

	ch := make(chan int, numThreads)
	for i := 0; i < 10; i++ {
		go func(tid int) {
			v, ok := queue.Dequeue()
			if ok {
				ch <- v
			}
		}(i)
	}
	for i := 0; i < numThreads; i++ {
		fmt.Printf("dequeued: %d\n", <-ch)
	}
}
