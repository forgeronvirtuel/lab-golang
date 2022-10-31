package main

import (
	"fmt"
	"time"
)

func main() {
	done := make(chan bool, 1)
	go spinner(100*time.Millisecond, done)
	const n = 45
	res := fibonacci(n)
	done <- true
	fmt.Printf("\rHeavy(%d) = %d\n", n, res)
}

func spinner(delay time.Duration, ch chan bool) {
	for {
		fmt.Println(len(ch))
		for _, r := range `|/-\` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func fibonacci(n int) int {
	if n < 2 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}
