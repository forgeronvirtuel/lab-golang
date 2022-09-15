package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
)

type CopyError struct {
	Source string
	Err    error
}

func (c *CopyError) Error() string {
	return fmt.Sprintf("Cannot copy source file %s: %v", c.Source, c.Err)
}

type OpenError struct {
	File string
	Err  error
}

func (o *OpenError) Error() string {
	return fmt.Sprintf("Cannot open file %s: %v", o.File, o.Err)
}

func main() {
	cmd := os.Args[1]
	switch cmd {
	case "generation":
		limit, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		folder := os.Args[3]
		fmt.Printf("Writing %d files into %s\n", limit, folder)
		generateFilesCommand(limit, folder)
		break
	case "copy":
		copyCommand(os.Args[2:])
		break
	case "copy-folder":
		copyFolderCommand(os.Args[2], os.Args[3])
		break
	default:
		fmt.Printf("Unrecognized commands %s", cmd)
	}
}

func copyCommand(args []string) {
	destinationDir := args[0]
	var destinations []string
	sources := os.Args[1:]
	for _, src := range sources {
		bse := path.Base(src)
		destinations = append(destinations, path.Join(destinationDir, bse))
	}

	if err := copyFiles(sources, destinations); err != nil {
		fmt.Println(err)
	}
}

func copyFolderCommand(sourceFolder, destinationFolder string) {
	var src, dst []string
	dirsItem, err := os.ReadDir(sourceFolder)
	if err != nil {
		fmt.Println(&OpenError{sourceFolder, err})
		return
	}

	for _, item := range dirsItem {
		if item.IsDir() {
			continue
		}
		src = append(src, path.Join(sourceFolder, item.Name()))
		dst = append(dst, path.Join(destinationFolder, item.Name()))
	}

	if err := copyFiles(src, dst); err != nil {
		fmt.Println(err)
	}
}

func generateFilesCommand(limit int, folder string) {
	for i := 0; i < limit; i++ {
		filename := path.Join(folder, fmt.Sprintf("f%d.txt", i))
		f, err := os.Create(filename)

		fmt.Printf("\r[%f%%] %s", float32(i+1)/float32(limit)*100, filename)

		if err != nil {
			fmt.Println(err)
		}

		_, err = f.Write([]byte("Hello, World !"))
		if err != nil {
			fmt.Println(err)
		}

		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Print("\nGeneration completed\n")
}

func copyFiles(sources, destinations []string) error {
	if len(sources) != len(destinations) {
		return fmt.Errorf("Trying to copy %d sources into %d destinations.", len(sources), len(destinations))
	}

	// Bad idea: stacking file descriptor until the function is finished
	nb_file := len(sources)
	for i := 0; i < nb_file; i++ {

		fmt.Printf("\r[%f%%] %s -> %s", float32(i+1)/float32(nb_file)*100, sources[i], destinations[i])

		// Create the destination file
		dstFile, err := os.Create(destinations[i])
		defer dstFile.Close()
		if err != nil {
			return &OpenError{File: destinations[i], Err: err}
		}

		// Open the source file
		srcFile, err := os.Open(sources[i])
		defer srcFile.Close()
		if err != nil {
			return &OpenError{File: sources[i], Err: err}
		}

		// Copying content
		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return &CopyError{Source: sources[i], Err: err}
		}
	}
	return nil
}
