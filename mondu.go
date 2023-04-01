package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
)

// mondu doesn't stop on errors, so set this to false with the -q command line
// option if you don't want to see errors.
var showErrors bool = true

// For writing to stderr.
var errfmt = log.New(os.Stderr, "", 0)

// Calculate the total recursive size of the file or directory names given.
func Mondu(fnames []string) int64 {
	var total int64

	for size := range *getSizesConcurrently(fnames) {
		total += size
	}

	return total
}

func getSizesConcurrently(fnames []string) *chan int64 {
	wg := sync.WaitGroup{}
	results := make(chan int64)

	for _, fname := range fnames {
		go calculateSize(fname, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return &results
}

func calculateSize(pname string, wg *sync.WaitGroup, results chan int64) {
	wg.Add(1)
	defer wg.Done()

	info, err := os.Lstat(pname)

	if err != nil {
		if showErrors {
			errfmt.Println(err)
		}
		return
	}

	if info.Mode()&fs.ModeSymlink != 0 { // Symbolic Link
		// Do nothing.
	} else if info.Mode()&fs.ModeDevice != 0 { // Block Device
		// Do nothing--a block device doesn't really contribute to the size.
	} else if info.Mode()&fs.ModeType == 0 { // Regular File
		size := info.Size()

		if size != 0 {
			results <- size
		}
	} else if info.IsDir() {
		files, err := ioutil.ReadDir(pname)

		if err != nil {
			if showErrors {
				errfmt.Println(err)
			}
			return
		}

		for _, file := range files {
			go calculateSize(path.Join(pname, file.Name()), wg, results)
		}
	} else {
		if showErrors {
			errfmt.Println("WARNING: info.Mode() =", info.Mode(), "is unhandled.")
		}
	}
}

func main() {
	var fnames []string

	for _, arg := range os.Args {
		// For the sole command line argument.
		if arg == "-q" || arg == "--quiet" {
			showErrors = false
		} else {
			fnames = append(fnames, arg)
		}
	}

	size := Mondu(fnames)
	fmt.Println(size)
}
