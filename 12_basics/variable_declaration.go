package basics

import "fmt"

func f() {
	var i uint
	i = 20
	fmt.Printf("i = (%T) %v \n", i, i)

	j := 30
	fmt.Printf("j = (%T) %v \n", j, j)

	v := 400
	k := uint8(v) // Oops, capacity overflow
	fmt.Printf("k = (%T) %v \n", k, k)
}

func DeclarationsExperiment(param int) string {
	var i = param
	return fmt.Sprintf("%d", i)
}

func ShortDeclarationsExperiment(param int) string {
	i := uint8(param)
	return fmt.Sprintf("%d", i)
}

func pointers() {
	var x int64
	x = 30

	*(&x) = 40

	fmt.Printf("%d\n", x)
}
