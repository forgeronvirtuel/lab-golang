package main

import (
	"fmt"
	"sort"
	"time"
)

// Inefficient way to compute fibonacci
func fibonacci(n int) int {
	if n < 2 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func Compute(input []int, limit int) []int {
	res := make([]int, len(input))
	resultsch := make(chan int, limit)
	syncch := make(chan bool, limit)

	// Launch computations
	for _, n := range input {
		go func(n int, resultsch chan<- int) {
			token := <-syncch
			resultsch <- fibonacci(n)
			syncch <- token
		}(n, resultsch)
	}

	// Gives token to starts goroutine
	for i := 0; i < limit; i++ {
		syncch <- true
	}

	// Reading results
	for i := 0; i < len(input); i++ {
		res[i] = <-resultsch
	}

	return res
}

func main() {
	data := []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
	}

	for i := 1; i <= 1000; i++ {
		var times []time.Duration
		for j := 0; j < 10; j++ {
			start := time.Now()
			Compute(data, i)
			times = append(times, time.Since(start))
		}
		sort.Slice(times, func(i, j int) bool {
			return times[i].Microseconds() < times[j].Microseconds()
		})
		fmt.Printf("#%.02d, tmin = %s, tavg = %s, tmax = %s\n", i, times[0], times[len(times)/2], times[len(times)-1])
	}
}
