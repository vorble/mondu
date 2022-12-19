package main

import (
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
)

type Progress struct {
	size int64
}

// path may be a file or directory path.
func recursive_size_find(pname string, waitGroup *sync.WaitGroup, results chan *Progress) {
	info, err := os.Lstat(pname)
	if err != nil {
		log.Println("WARNING: stat():", err)
		waitGroup.Done()
		return
	}
	var partial Progress
	if info.Mode()&fs.ModeSymlink != 0 {
		// Don't follow symbolic links.
	} else if info.IsDir() {
		files, err := ioutil.ReadDir(pname)
		if err != nil {
			log.Println("WARNING: ReadDir():", err)
			waitGroup.Done()
			return
		}
		for _, file := range files {
			waitGroup.Add(1)
			go recursive_size_find(path.Join(pname, file.Name()), waitGroup, results)
		}
	} else {
		partial.size += info.Size()
		results <- &partial
	}
	waitGroup.Done()
}

func main() {
	var waitGroup sync.WaitGroup
	var waitGroup2 sync.WaitGroup
	results := make(chan *Progress)
	var progress Progress

	// User gives a list of files as arguments to the program.
	for _, arg := range os.Args[1:] {
		waitGroup.Add(1)
		go recursive_size_find(arg, &waitGroup, results)
	}

	waitGroup2.Add(1)
	go func() {
		go func() {
			waitGroup.Wait()
			close(results)
		}()
		for partial := range results {
			progress.size += partial.size
		}
		waitGroup2.Done()
	}()

	waitGroup2.Wait()
	log.Println("Total:", progress.size)
}
