package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/wang502/wfqueue"
)

func enqueueItems() {
	numItems, _ := strconv.Atoi(os.Args[1])
	numThreads, _ := strconv.Atoi(os.Args[2])

	queue := wfqueue.NewWFQueue(numThreads)
	doneChan := make(chan bool, numItems*numThreads)

	start := time.Now()
	for i := 0; i < numThreads; i++ {
		go func(tid int) {
			for j := 0; j < numItems; j++ {
				fmt.Printf("thread %d tries to enqueue\n", tid)
				queue.Enqueue(j, tid)
				doneChan <- true
			}
		}(i)
	}

	for n := 0; n < numItems*numThreads; n++ {
		<-doneChan
		fmt.Printf("finished an operation, n:%d\n", n)
	}
	elapsed := time.Since(start)
	fmt.Printf("[WF][%d threads] each enqueue [%d items] takes total %fs\n", numThreads, numItems, elapsed.Seconds())
}

func enqueueDequeuePair() {
	numTimes, _ := strconv.Atoi(os.Args[1])
	numThreads, _ := strconv.Atoi(os.Args[2])

	queue := wfqueue.NewWFQueue(numThreads)
	doneChan := make(chan bool, numTimes)
	start := time.Now()
	ch := make(chan bool, numThreads)

	for i := 0; i < numTimes; i++ {
		ch <- true
		go func(v int, ch chan bool) {
			if v%2 == 0 {
				queue.Enqueue(v, v%numThreads)
			} else {
				queue.Dequeue(v % numThreads)
			}

			doneChan <- true
			<-ch
		}(i, ch)
	}

	for j := 0; j < numTimes; j++ {
		<-doneChan
	}

	elapsed := time.Since(start)
	fmt.Printf("[WF][%d threads][%d times enqueue/dequeue pair] takes %fs\n", numThreads, numTimes, elapsed.Seconds())
}

func enqueueDequeuePairII() {
	numTimes, _ := strconv.Atoi(os.Args[1])
	numThreads, _ := strconv.Atoi(os.Args[2])

	queue := wfqueue.NewWFQueue(numThreads)
	doneChan := make(chan bool, numTimes*numThreads)
	start := time.Now()

	for i := 0; i < numThreads; i++ {
		go func(tid int) {
			for j := 0; j < numTimes; j++ {
				if j%2 == 0 {
					fmt.Printf("thread %d tries to enqueue\n", tid)
					queue.Enqueue(j, tid)
				} else {
					fmt.Printf("thread %d tries to dequeue\n", tid)
					queue.Dequeue(tid)
				}

				doneChan <- true
			}
		}(i)
	}

	for n := 0; n < numTimes*numThreads; n++ {
		<-doneChan
		fmt.Printf("finished an operation, n:%d\n", n)
	}

	elapsed := time.Since(start)
	fmt.Printf("[WF][%d threads][%d times enqueue/dequeue pair] takes %fs\n", numThreads, numTimes, elapsed.Seconds())
}

func main() {
	enqueueItems()
	//enqueueDequeuePairII()
}
