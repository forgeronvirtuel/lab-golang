package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)
	go waitForData(ch)
	fmt.Println("main goroutine: waiting for two seconds")
	time.Sleep(2 * time.Second)
	fmt.Println("main goroutine: sending 1")
	ch <- 1
	fmt.Println("main goroutine: receiving data")
	x := <-ch
	fmt.Printf("main goroutine: received %d\n", x)
	fmt.Println("main goroutine: waiting one second")
	time.Sleep(1 * time.Second)
	fmt.Println("main goroutine: end")
}

func waitForData(ch chan int) {
	fmt.Println("waitForData: receiving data")
	x := <-ch
	fmt.Printf("waitForData goroutine: received %d\n", x)
	fmt.Println("waitForData goroutine: waiting 2 seconds")
	time.Sleep(2 * time.Second)
	fmt.Println("waitForData goroutine: sending 2")
	ch <- 2
	fmt.Println("waitForData goroutine: end")
}
