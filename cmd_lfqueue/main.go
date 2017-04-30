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
	doneChan := make(chan bool, numItems*numThreads)

	start := time.Now()
	for i := 0; i < numThreads; i++ {
		go func() {
			for j := 0; j < numItems; j++ {
				queue.Enqueue(j)
				doneChan <- true
			}
		}()
	}

	for n := 0; n < numItems*numThreads; n++ {
		<-doneChan
	}
	elapsed := time.Since(start)
	fmt.Printf("[LF][%d threads][%d items] enqueue takes %fs\n", numThreads, numItems, elapsed.Seconds())
}

func enqueueDequeuePair() {
	numTimes, _ := strconv.Atoi(os.Args[1])
	numThreads, _ := strconv.Atoi(os.Args[2])

	queue := wfqueue.NewLFQueue()
	doneChan := make(chan bool, numTimes*numThreads)
	start := time.Now()

	for i := 0; i < numThreads; i++ {
		go func() {
			for j := 0; j < numTimes; j++ {
				if j%2 == 0 {
					queue.Enqueue(j)
				} else {
					queue.Dequeue()
				}

				doneChan <- true
			}
		}()
	}

	for n := 0; n < numTimes*numThreads; n++ {
		<-doneChan
	}

	elapsed := time.Since(start)
	fmt.Printf("[LF][%d threads][%d times enqueue/dequeue pair] takes %fs\n", numThreads, numTimes, elapsed.Seconds())
}

func main() {
	//enqueueItems()
	enqueueDequeuePair()
}
