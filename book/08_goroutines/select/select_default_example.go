package main

import (
	"fmt"
	"runtime"
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

	var m runtime.MemStats

	timeoutch := make(chan bool)

	ch := make(chan int)

	go func() {
		ch <- fibonacci(inp)
	}()

	go func() {
		time.Sleep(timeout)
		timeoutch <- true
	}()

	stop := false
	for !stop {
		time.Sleep(300 * time.Millisecond)
		select {
		case result := <-ch:
			fmt.Printf("fibonacci(%d) = %d", inp, result)
			stop = true
		case <-timeoutch:
			fmt.Println("timeout reach !")
			stop = true
		default:
			runtime.ReadMemStats(&m)
			fmt.Printf("TotalAlloc: %d Bytes\n", m.TotalAlloc)
		}
	}
}
