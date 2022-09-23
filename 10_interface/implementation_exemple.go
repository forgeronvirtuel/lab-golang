package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

//// No explicit implementation of io.Writer
//type IndexedMap struct {
//	// Very important properties
//}
//
//func (m *IndexedMap) Write(buf []byte) (n int, err error) {
//	// Write into the indexed map
//}
//
//func AddData(writer io.Writer) {
//	// Do very important stuff
//}
//
//func AddDataAndClose(stream io.WriteCloser) {
//	// Do very important stuff
//}

const debug = false
const buffer_version = false
const file_version = true

func main() {
	if file_version {
		var file *os.File
		if debug {
			file, _ = os.Open("debug.log")
		}
		fmt.Printf("%T %v\n", file)
		f(file)
	}

	if buffer_version {
		var buf *bytes.Buffer
		if debug {
			buf = new(bytes.Buffer)
		}
		fmt.Printf("%T %v\n", buf)
		f(buf)
		if debug {
			fmt.Println(buf)
		}
	}
}

func f(out io.Writer) {
	fmt.Printf("%T %v\n", out)
	if out != nil {
		out.Write([]byte("done !"))
	}
}
