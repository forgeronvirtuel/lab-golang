package main

import (
	"github.com/forgeronvirtuel/indexedmap/indexedmap"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "", log.Lshortfile)
	imap := indexedmap.NewIndexedMap()
	imap.Add("key", "value")
	logger.Println("Content of `key`:", imap.Get("key"))
}
