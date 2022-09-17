package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Configure the logger as you need it
	logger := log.Default()

	// Catch and log the panic error
	defer func() {
		if p := recover(); p != nil {
			logger.Println(p)
		}
	}()

	// Simple command system
	switch os.Args[1] {
	case "command-a":
		logger.Println("Do stuff for command a")
	case "command-b":
		logger.Println("Do stuff for command b")
	case "command-c":
		logger.Println("Do stuff for command c")
	default:
		panic(fmt.Errorf("Command %s unknown", os.Args[1]))
	}
}
