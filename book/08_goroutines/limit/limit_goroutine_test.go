package main

import (
	"fmt"
	"testing"
)

func BenchmarkCompute(b *testing.B) {
	data := []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
	}

	for i := 1; i < 10; i++ {
		b.Run(fmt.Sprintf("Goroutines: %d", i), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				Compute(data, i)
			}
		})
	}
}
