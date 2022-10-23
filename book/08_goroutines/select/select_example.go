package main

import (
	"fmt"
	"time"
)

// Inefficient way to compute fibonacci
func fibonacci(n int) int {
	if n < 2 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func main() {
	const inp = 10
	const timeout = 10 * time.Second

	ch := make(chan int)
	timeoutch := make(chan bool)

	go func() {
		ch <- fibonacci(inp)
	}()

	go func() {
		time.Sleep(timeout)
		timeoutch <- true
	}()

	select {
	case result := <-ch:
		fmt.Printf("fibonacci(%d) = %d", inp, result)
	case <-timeoutch:
		fmt.Println("timeout reach !")
	}
}
