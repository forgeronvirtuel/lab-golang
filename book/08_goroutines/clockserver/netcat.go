package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Println(err)
		}
	}()

	if _, err := io.Copy(os.Stdout, conn); err != nil {
		log.Println(err)
	}
}
