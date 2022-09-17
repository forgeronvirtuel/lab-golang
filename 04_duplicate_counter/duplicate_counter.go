package main

import (
	"bufio"
	"fmt"
	"os"
)

func duplicateCount(lines []string) map[string]int {
	counts := make(map[string]int)
	for _, l := range lines {
		counts[l]++
	}
	return counts
}

func readFile(filepath string) ([]string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines, nil
}

func main() {
	var text []string
	for _, filepath := range os.Args[1:] {
		lines, err := readFile(filepath)
		if err != nil {
			panic(err)
		}
		text = append(text, lines...)
	}

	counts := duplicateCount(text)
	for k, v := range counts {
		fmt.Printf("%s: %d\n", k, v)
	}
}
