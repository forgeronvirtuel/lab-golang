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

func ComputeWithLimit(input []int, limit int) []int {
	res := make([]int, len(input))
	resultsch := make(chan int, limit)
	inputch := make(chan int)

	// Inject all inputs
	go func() {
		for _, v := range input {
			inputch <- v
		}
		close(inputch)
	}()

	// Launch computations
	for i := 0; i < limit; i++ {
		go func(inpch <-chan int, outch chan<- int) {
			for n := range inpch {
				outch <- fibonacci(n)
			}
		}(inputch, resultsch)
	}

	// Reading results
	for i := 0; i < len(input); i++ {
		res[i] = <-resultsch
	}

	return res
}

func main() {
	const startAt = 1
	const nbThread = 40
	const nbData = 10000
	var computeFunc = ComputeWithLimit

	data := make([]int, nbData)
	for i := 0; i < nbData; i++ {
		data[i] = 20
	}

	totTime := time.Duration(0)
	for i := startAt; i <= nbThread; i++ {
		var times []time.Duration
		for j := 0; j < 5; j++ {
			start := time.Now()
			computeFunc(data, i)
			totTime += time.Since(start)
			times = append(times, time.Since(start))
		}
		sort.Slice(times, func(i, j int) bool {
			return times[i].Microseconds() < times[j].Microseconds()
		})
		avg := (times[0] + times[1]) / 2
		fmt.Printf("#%.02d, avg = %s, tmin = %s, tavg = %s, tmax = %s\n", i, avg, times[0], times[len(times)/2], times[len(times)-1])
	}
}
