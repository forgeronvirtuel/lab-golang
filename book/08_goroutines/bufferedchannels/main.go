package main

import (
	"crypto/rand"
	"fmt"
)

type Message []byte

const nbProducer = 1
const nbConsumer = 1
const nbMessage = 200
const msgSize = 1024

func main() {
	// Communication channel
	queue := make(chan Message, 1)

	// Synchronization
	endConsumer := make(chan int)
	endProducer := make(chan int)

	// producers
	for i := 0; i < nbProducer; i++ {
		go readMessage(queue, i, endProducer)
	}

	// consumers
	for i := 0; i < nbConsumer; i++ {
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

	fmt.Println("Main goroutine end.")
}

func readMessage(queue chan<- Message, producerNo int, end chan<- int) {
	for i := 0; i < nbMessage; i++ {
		data := make(Message, msgSize)
		rand.Read(data)
		queue <- data
	}
	fmt.Printf("Producer #%d finished\n", producerNo)
	end <- producerNo
}

func filters(queue <-chan Message, consumerNo int, end chan<- int) {
	var cnt int
	for packet := range queue {
		msb := packet[0] & 0x80
		if msb != 0 {
			cnt++
		}
	}
	fmt.Printf("Consumer #%d finished, filtered %d packet\n", consumerNo, cnt)
	end <- consumerNo
}
