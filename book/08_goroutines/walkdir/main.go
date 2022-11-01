package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var verbose = flag.Bool("v", false, "show verbose progress messages")

func main() {

	flag.Parse()

	// Print the results periodically
	var tick <-chan time.Time
	if *verbose {
		tick = time.Tick(500 * time.Millisecond)
	}

	// Determines the initial directory
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	// Traverse the file tree
	filesizes := make(chan int64)
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go walkdir(root, &n, filesizes)
	}

	go func() {
		n.Wait()
		close(filesizes)
	}()

	var nfiles, nbytes int64
	var stop bool

	for !stop {
		select {
		case size, ok := <-filesizes:
			if !ok {
				stop = true
				break
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}
	printDiskUsage(nfiles, nbytes)
}

func printDiskUsage(nfiles int64, nbytes int64) {
	fmt.Printf("%d files %.1f GB\n", nfiles, float64(nbytes)/1e9)
}

func walkdir(dir string, n *sync.WaitGroup, filesizes chan<- int64) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkdir(subdir, n, filesizes)
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

var sema = make(chan struct{}, 20)

func dirents(dir string) []os.DirEntry {
	sema <- struct{}{}
	defer func() { <-sema }()
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
	}
	return entries
}
