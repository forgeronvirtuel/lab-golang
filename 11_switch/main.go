package main

import (
	"fmt"
	"os"
)

func main() {
	Freelance()
}

func Freelance() {
	switch os.Args[1] {
	case "Go":
		fmt.Println("Contact me !")
		break
	case "JavaScript":
		fmt.Println("Contact me !")
		break
	case "Python":
		fmt.Println("Contact me !")
	default:
		fmt.Println("Oups, need to learn !")
	}
}

func Freelance2() {
	switch os.Args[1] {
	case "Go", "Python", "JavaScript":
		fmt.Println("Contact me !")
	default:
		fmt.Println("Oups, need to learn !")
	}
}
