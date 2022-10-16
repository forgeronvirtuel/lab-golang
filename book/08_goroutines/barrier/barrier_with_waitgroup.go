package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

const nbProducer = 5
const nbConsumer = 5
const nbMessage = 10
const sizeBuffer = nbMessage * nbProducer

func main() {
	go spinner(100 * time.Millisecond)

	start := time.Now()

	// Communication channel
	params := make(chan int, sizeBuffer)
	results := make(chan int, sizeBuffer)

	// Synchronization
	var wg sync.WaitGroup

	// Launching all producers
	wg.Add(nbProducer)
	for i := 1; i <= nbProducer; i++ {
		go produceInt(params, nbMessage, i, &wg)
	}

	wg.Wait()
	close(params)

	// Launching all workers
	wg.Add(nbConsumer)
	for i := 1; i <= nbConsumer; i++ {
		go worker(params, results, i, &wg)
	}

	// Second barrier
	wg.Wait()

	fmt.Println("\rMain goroutine end.")
	log.Printf("Time to execute: %s", time.Since(start))
}

func produceInt(queue chan<- int, nbMessage, producerNo int, end *sync.WaitGroup) {
	defer end.Done()
	for i := 0; i < nbMessage; i++ {
		queue <- i
	}
	fmt.Printf("\rProducer #%d finished\n", producerNo)
}

func worker(queue <-chan int, results chan int, consumerNo int, end *sync.WaitGroup) {
	defer end.Done()
	for i := range queue {
		results <- fibonacci(i)
	}
	fmt.Printf("\rConsumer #%d finished\n", consumerNo)
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `|/-\` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

// Inefficient way to compute fibonacci
func fibonacci(n int) int {
	if n < 2 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}
