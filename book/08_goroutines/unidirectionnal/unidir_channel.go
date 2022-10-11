package main

import (
	"crypto/rand"
	"fmt"
)

type ImageData []byte

func getImages(out chan<- ImageData, numbers int) {
	for i := 0; i < numbers; i++ {
		data := make(ImageData, 10)
		rand.Read(data) // replace with a way to get data
		out <- data
	}
	close(out)
}

func removeWhite(out chan<- ImageData, in <-chan ImageData) {
	for data := range in {
		for i := 0; i < len(data); i++ {
			data[i] -= 0x05
		}
		out <- data
	}
	close(out)
}

func isolate(out chan<- ImageData, in <-chan ImageData, min, max byte) {
	for data := range in {
		for i := 0; i < len(data); i++ {
			if data[i] < min || max < data[i] {
				data[i] = 0
			}
		}
		out <- data
	}
	close(out)
}

func main() {
	genDataCh := make(chan ImageData)
	rmWhiteCh := make(chan ImageData)
	isolateCh := make(chan ImageData)

	go getImages(genDataCh, 100)
	go removeWhite(rmWhiteCh, genDataCh)
	go isolate(isolateCh, rmWhiteCh, 0xAB, 0xAC)

	for res := range isolateCh {
		fmt.Println(res)
	}
}
