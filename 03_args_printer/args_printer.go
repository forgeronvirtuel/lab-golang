package main

import (
	"fmt"
	"os"
)

func main() {
	// v1
	printArgV1()

	// v2
	printArgV2()

	// v3
	printArgV3()

	// v4
	printArgV4()
}

func printArgV1() {
	for _, arg := range os.Args[0 : len(os.Args)-1] {
		fmt.Print(arg, " ")
	}
	fmt.Println(os.Args[len(os.Args)-1])
}

func printArgV2() {
	i := 0
	for {
		fmt.Print(os.Args[i], " ")
		i++
		if i == len(os.Args)-1 {
			break
		}
	}
	fmt.Println(os.Args[len(os.Args)-1])
}

func printArgV3() {
	i := 0
	for i < len(os.Args)-1 {
		fmt.Print(os.Args[i], " ")
		i++
	}
	fmt.Println(os.Args[len(os.Args)-1])
}

func printArgV4() {
	s := ""
	for i, v := range os.Args {
		if i == len(os.Args)-1 {
			s += v
			break
		}
		s += v + " "
	}
	fmt.Println(s)
}
