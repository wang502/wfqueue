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

	queue := wfqueue.NewLFQueue()
	doneChan := make(chan bool, numItems)
	start := time.Now()
	ch := make(chan bool, numThreads)

	for i := 0; i < numItems; i++ {
		ch <- true
		go func(v int, ch chan bool) {
			queue.Enqueue(v)
			doneChan <- true
			<-ch
		}(i, ch)
	}

	for j := 0; j < numItems; j++ {
		<-doneChan
	}

	elapsed := time.Since(start)
	fmt.Printf("[LF][%d threads][%d items] enqueue takes %fs\n", numThreads, numItems, elapsed.Seconds())
}

func enqueueDequeuePair() {
	numTimes, _ := strconv.Atoi(os.Args[1])
	numThreads, _ := strconv.Atoi(os.Args[2])

	queue := wfqueue.NewLFQueue()
	doneChan := make(chan bool, numTimes)
	start := time.Now()
	ch := make(chan bool, numThreads)

	for i := 0; i < numTimes; i++ {
		ch <- true
		go func(v int, ch chan bool) {
			if v%2 == 0 {
				queue.Enqueue(v)
			} else {
				queue.Dequeue()
			}

			doneChan <- true
			<-ch
		}(i, ch)
	}

	for j := 0; j < numTimes; j++ {
		<-doneChan
	}

	elapsed := time.Since(start)
	fmt.Printf("[LF][%d threads][%d times enqueue/dequeue pair] takes %fs\n", numThreads, numTimes, elapsed.Seconds())
}

func main() {
	//enqueueItems()
	enqueueDequeuePair()
}
