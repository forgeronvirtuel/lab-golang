package main

import (
	"bytes"
	"crypto/rand"
)

const nbConsumer = 10
const nbMessage = 20

func getQuery() []byte {
	data := make([]byte, 1024)
	rand.Read(data)
	return data
}

func queryChecker(packets [][]byte, word []byte) bool {
	input := make(chan []byte, nbMessage)
	res := make(chan bool)

	// Create a pool of goroutine to treat data
	for i := 0; i < nbConsumer; i++ {
		go func(n int, input <-chan []byte, res chan<- bool) {
			data := <-input

			// Check the query
			res <- bytes.Compare(data, word) == 0
		}(i, input, res)
	}

	// insert all data for checking
	for _, p := range packets {
		input <- p
	}

	// stop at the first wrong result
	for ok := range res {
		if !ok {
			return false // all unfinished goroutine will be locked
		}
	}

	return true
}
