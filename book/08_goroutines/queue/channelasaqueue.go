package main

import "fmt"

func main() {
	ch := make(chan int, 3)

	for i := 0; i < 3; i++ {
		ch <- i * i
	}

	for i := 0; i < 2; i++ {
		fmt.Println(<-ch)
	}

	// Deadlock here
	for i := 3; i < 6; i++ {
		ch <- i * i
	}
}
