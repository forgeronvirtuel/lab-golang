package main

import "fmt"

func main() {
	naturals := make(chan int)
	fibonaccis := make(chan int)

	go func() {
		for n := 0; n < 40; n++ {
			naturals <- n
		}
		close(naturals)
	}()

	go func() {
		for {
			n, ok := <-naturals
			if !ok {
				break
			}
			fibonaccis <- fibonacci(n)
		}
		close(fibonaccis)
	}()

	for f := range fibonaccis {
		fmt.Println(f)
	}
}

func fibonacci(n int) int {
	if n < 2 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}
