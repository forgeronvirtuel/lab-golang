package main

import "fmt"

func main() {
	n_lists := []int{23, 44, 13, 46, 29}
	var channels = make([]chan int, len(n_lists))

	// Launch one goroutine for each computation
	for i, n := range n_lists {
		ch := make(chan int)
		go func(n int, ch chan int) {
			ch <- fibonacci(n)
		}(n, ch)
		channels[i] = ch
	}

	//Wait in order for each goroutine to end
	for i, n := range n_lists {
		ch := channels[i]
		res := <-ch
		fmt.Printf("fibonacci(%d) = %d\n", n, res)
	}
}

// Inefficient way to compute fibonacci
func fibonacci(n int) int {
	if n < 2 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}
