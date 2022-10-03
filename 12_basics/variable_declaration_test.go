package basics

import (
	"fmt"
	"testing"
)

func BenchmarkDeclarationsExperiment(b *testing.B) {
	for i := 0; i < b.N; i++ {
		j := uint8(i)
		if j < 0 {
			fmt.Sprintf("%d", j)
		}
	}
}

func BenchmarkShortDeclarationsExperiment(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var j = i
		if j < 0 {
			fmt.Sprintf("%d", j)
		}
	}
}

func TestPointers(t *testing.T) {
	var x int64

	var xPtr *int64
	var xPtrPtr **int64
	var xppp ***int64

	x = 30
	xPtr = &x       // memory address of x
	xPtrPtr = &xPtr // memory address of xPtr
	xppp = &xPtrPtr // memory address of xPtrPtr

	***xppp = 40

	fmt.Printf("%d\n", x)
}
