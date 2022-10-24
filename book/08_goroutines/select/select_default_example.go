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
	const inp = 50
	const timeout = 1 * time.Second

	ch := make(chan int)

	go func() {
		ch <- fibonacci(inp)
	}()

	time.Sleep(timeout)
	select {
	case result := <-ch:
		fmt.Printf("fibonacci(%d) = %d", inp, result)
	default:
		fmt.Println("timeout reach !")
	}
}
