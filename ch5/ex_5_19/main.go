package main

import (
	"errors"
	"fmt"
)

func main() {
	v := panicAndRecover()
	fmt.Println(v)
}

func panicAndRecover() (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = p.(error)
		}
	}()
	genErr := errors.New("This is an error")
	panic(genErr)
}
