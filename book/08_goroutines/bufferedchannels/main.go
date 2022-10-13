package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"time"
)

type Message []byte

const nbProducer = 5
const nbConsumer = 5
const nbMessage = 100
const msgSize = 256
const sizeBuffer = 3

func main() {
	start := time.Now()

	// Communication channel
	queue := make(chan Message, sizeBuffer)

	// Synchronization
	endConsumer := make(chan int)
	endProducer := make(chan int)

	// In the meanwhile, display a spinner
	go spinner(100 * time.Millisecond)

	// producers
	for i := 1; i <= nbProducer; i++ {
		go readMessage(queue, i, endProducer)
	}

	// consumers
	for i := 1; i <= nbConsumer; i++ {
		go filters(queue, i, endConsumer)
	}

	// Once all producer finished, close the queue
	for i := 0; i < nbProducer; i++ {
		<-endProducer
	}
	close(queue)

	// Sync the endConsumer of all tasks
	for i := 0; i < nbConsumer; i++ {
		<-endConsumer
	}

	fmt.Println("\rMain goroutine end.")
	log.Printf("Time to execute: %s", time.Since(start))
}

func readMessage(queue chan<- Message, producerNo int, end chan<- int) {
	for i := 0; i < nbMessage; i++ {
		data := make(Message, msgSize)
		rand.Read(data)
		queue <- data
	}
	fmt.Printf("\rProducer #%d finished\n", producerNo)
	end <- producerNo
}

func filters(queue <-chan Message, consumerNo int, end chan<- int) {
	var cnt, all int
	for packet := range queue {
		all++
		msb := packet[0] & 0x80
		if msb != 0 {
			cnt++
		}
	}
	fmt.Printf("\rConsumer #%d finished, filtered %d/%d packet\n", consumerNo, cnt, all)
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
