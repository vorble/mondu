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

// Should errors be shown.
var showErrors bool = true

// For writing to stderr.
var errfmt = log.New(os.Stderr, "", 0)

// Calculate the total recursive size of the file or directory names given.
func mondu(fnames []string) int64 {
	// Before each call to mondu_(), be sure to call wg.Add(1). mondu_() will call wg.Done() when it
	// is finished.
	var wg sync.WaitGroup
	// The results are a series of sizes that should be totaled.
	results := make(chan int64)

	for _, fname := range fnames {
		wg.Add(1)
		go mondu_(fname, &wg, results)
	}

	// The for loop that follows terminates after the wait group finishes.
	go func() {
		wg.Wait()
		close(results)
	}()

	var total int64

	for size := range results {
		total += size
	}

	return total
}

// path may be a file or directory path.
func mondu_(pname string, wg *sync.WaitGroup, results chan int64) {
	// Every return path must call wg.Done(), this is done via defer.
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
			wg.Add(1)
			go mondu_(path.Join(pname, file.Name()), wg, results)
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

	size := mondu(fnames)
	fmt.Println(size)
}
