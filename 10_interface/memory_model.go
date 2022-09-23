package main

import (
	"io"
	"os"
)

func main() {
	var w io.Writer
	w = os.Stdout
	w.Write([]byte("Hello, World!"))
}
