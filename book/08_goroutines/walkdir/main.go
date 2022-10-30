package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Determines the initial directory
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// Traverse the file tree
	filesizes := make(chan int64)
	go func() {
		for _, root := range roots {
			walkdir(root, filesizes)
		}
		close(filesizes)
	}()

	// Prints the results
	var nfiles, nbytes int64
	for size := range filesizes {
		nfiles++
		nbytes += size
	}
	fmt.Printf("%d files %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

func walkdir(dir string, filesizes chan<- int64) {
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			walkdir(subdir, filesizes)
		} else {
			info, err := entry.Info()
			if err != nil {
				fmt.Fprintf(os.Stderr, "du: %v\n", err)
			} else {
				filesizes <- info.Size()
			}
		}
	}
}

func dirents(dir string) []os.DirEntry {
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
	}
	return entries
}
