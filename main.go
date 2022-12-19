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
	} else if info.Mode()&fs.ModeType == 0 { // Regular File
		partial.size += info.Size()
		results <- &partial
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
		log.Println("WARNING: info.Mode() =", info.Mode(), "is unhandled.")
	}
	waitGroup.Done()
}

func mondu(fnames []string) Progress {
	var waitGroup sync.WaitGroup
	results := make(chan *Progress)
	for _, fname := range fnames {
		waitGroup.Add(1)
		go recursive_size_find(fname, &waitGroup, results)
	}
	go func() {
		waitGroup.Wait()
		close(results)
	}()

	var progress Progress
	for partial := range results {
		progress.size += partial.size
	}
	return progress
}

func main() {
	progress := mondu(os.Args[1:])
	log.Println("Total:", progress.size)
}
