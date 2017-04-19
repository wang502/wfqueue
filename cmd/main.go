package main

import (
	"fmt"

	"github.com/wang502/wfqueue"
)

func main() {
	numThreads := 10
	queue := wfqueue.NewWFQueue(numThreads)

	fmt.Printf("queue size: %d\n", queue.Len())
	done := make(chan bool, numThreads)
	for i := 0; i < 10; i++ {
		go func(val, tid int) {
			queue.Enqueue(val, tid)
			done <- true
		}(i, i)
	}
	for i := 0; i < numThreads; i++ {
		fmt.Println(<-done)
	}
	close(done)
	fmt.Printf("queue: %s\n", queue)
	fmt.Printf("queue size: %d\n", queue.Len())

	ch := make(chan int, numThreads)
	for i := 0; i < 10; i++ {
		go func(tid int) {
			v, ok := queue.Dequeue(tid)
			if ok {
				ch <- v
			}
		}(i)
	}
	for i := 0; i < numThreads; i++ {
		fmt.Printf("dequeued: %d\n", <-ch)
	}
	fmt.Printf("queue: %s\n", queue)
	fmt.Printf("queue size: %d\n", queue.Len())
	close(ch)

	fmt.Printf("enqueue: 10\n")
	queue.Enqueue(10, 0)
	fmt.Printf("queue: %s\n", queue)
	fmt.Printf("queue size: %d\n", queue.Len())
	fmt.Printf("head: %d\n", queue.Head())
	fmt.Printf("tail: %d\n", queue.Tail())

	v, ok := queue.Dequeue(1)
	if ok {
		fmt.Printf("dequeued: %d\n", v)
	}
	fmt.Printf("queue: %s\n", queue)
}
