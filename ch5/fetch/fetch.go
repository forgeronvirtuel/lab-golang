package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

func main() {
	filename, n, err := fetchAndWrite("http://google.fr")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Filename: ", filename)
	fmt.Println("Bytes written: ", n)
}

func fetchAndWrite(url string) (filename string, n int64, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	local := path.Base(resp.Request.URL.Path)
	if local == "/" {
		local = "index.html"
	}

	f, err := os.Create(local)
	if err != nil {
		return "", 0, err
	}

	n, err = io.Copy(f, resp.Body)
	if err != nil {
		return "", 0, err
	}

	if closeErr := f.Close(); closeErr != nil {
		return "", 0, closeErr
	}

	return local, n, err
}
