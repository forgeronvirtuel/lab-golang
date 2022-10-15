package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"runtime"
	"time"
)

const nbConsumer = 10
const nbMessage = 20

func queryChecker(packets [][]byte, word []byte) bool {
	input := make(chan []byte, nbMessage)
	res := make(chan bool)

	// Create a pool of goroutine to treat data
	for i := 0; i < nbConsumer; i++ {
		go func(n int, input <-chan []byte, res chan<- bool) {
			data := <-input
			// Check the query
			res <- bytes.Contains(data, word)
		}(i, input, res)
	}

	// insert all data for checking
	for _, p := range packets {
		input <- p
	}

	// stop at the first wrong result
	for i := 0; i < nbConsumer; i++ {
		if !(<-res) {
			return false // all unfinished goroutine will be locked
		}
	}
	return true
}

const MBdiv = 1024.0 * 1024.0
const nbIter = 1_000

func main() {
	m := runtime.MemStats{}

	for i := 0; i < nbIter; i++ {
		packets := make([][]byte, nbMessage)
		for i, _ := range packets {
			data := make([]byte, 1024)
			rand.Read(data)
			word := []byte("HTTP")
			copy(data[0:len(word)], word)
			packets[i] = data
		}
		if ok := queryChecker(packets, []byte("HTTP")); !ok {
			fmt.Printf("[nok] ")
		} else {
			fmt.Printf("[ ok] ")
		}

		runtime.ReadMemStats(&m)
		fmt.Printf("HeapAlloc: %0.2f MB\r", float64(m.HeapAlloc)/MBdiv)
		time.Sleep(5 * time.Microsecond)
	}
}
