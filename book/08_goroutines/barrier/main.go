package main

import (
	"fmt"
	"log"
	"time"
)

const nbProducer = 5
const nbConsumer = 5
const nbMessage = 43
const sizeBuffer = nbMessage * nbProducer

func main() {
	go spinner(100 * time.Millisecond)

	start := time.Now()

	// Communication channel
	params := make(chan int, sizeBuffer)
	results := make(chan int, sizeBuffer)

	// Synchronization
	endProducer := make(chan int, nbProducer)
	endConsumer := make(chan int, nbConsumer)

	// Launching all producers
	for i := 1; i <= nbProducer; i++ {
		go produceInt(params, nbMessage, i, endProducer)
	}

	// First barrier
	for i := 0; i < nbProducer; i++ {
		<-endProducer
	}
	close(params)

	// Launching all workers
	for i := 1; i <= nbConsumer; i++ {
		go worker(params, results, i, endConsumer)
	}

	// Second barrier
	for i := 0; i < nbConsumer; i++ {
		<-endConsumer
	}

	fmt.Println("\rMain goroutine end.")
	log.Printf("Time to execute: %s", time.Since(start))
}

func produceInt(queue chan<- int, nbMessage, producerNo int, end chan<- int) {
	for i := 0; i < nbMessage; i++ {
		queue <- i
	}
	fmt.Printf("\rProducer #%d finished\n", producerNo)
	end <- producerNo
}

func worker(queue <-chan int, results chan int, consumerNo int, end chan<- int) {
	for i := range queue {
		results <- fibonacci(i)
	}
	fmt.Printf("\rConsumer #%d finished\n", consumerNo)
	end <- consumerNo
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
