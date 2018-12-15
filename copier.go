package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

func copy(inputFile, destinationFile string) {

	var wg sync.WaitGroup
	chans := make(chan int, 4) // Controlling go routines count to 4 or less at any point of time

	// Need to verify that if inputFile
	// actully exists
	if inputFileStat, err := os.Stat(inputFile); err != nil {
		fmt.Fprintln(os.Stderr, "inputFileError: ", err.Error())
		os.Exit(1)
		return

	} else {
		// If input file exist then
		// we can copy it to destiation
		fmt.Fprintln(os.Stdout, inputFileStat.Name(), inputFileStat.Size())
		var offset, partsize int64

		offset = 0                                 // offset from need to read
		partsize = inputFileStat.Size() / int64(8) // size of part we are going to write concurrently

		inputFileHandle, inputFileError := os.OpenFile(inputFile, os.O_RDONLY|os.O_RDWR, os.ModePerm)
		defer inputFileHandle.Close()
		if inputFileError != nil {
			fmt.Fprintln(os.Stderr, "Could not open input file ")
			os.Exit(1)
			return
		}

		destinationFileHandle, destinationFileError := os.OpenFile(destinationFile, os.O_CREATE|os.O_RDWR, os.ModePerm)
		defer destinationFileHandle.Close()
		if destinationFileError != nil {
			fmt.Fprintln(os.Stderr, "Could not create destination file")
			os.Exit(1)
			return
		}

		chans <- 1
		wg.Add(1)
		go watcher(destinationFile, chans, &wg)
		// Creating 8 concurrent routines
		// which will write the same file
		// from different offsets
		for i := 1; i < 9; i++ {
			wg.Add(1)
			chans <- 1
			go partProcess(i, offset, partsize, inputFileHandle, destinationFileHandle, &wg, chans)
			if offset > inputFileStat.Size() {
				break
			}
			offset += partsize + 1
			if i == 8 {
				<-chans
				partsize = inputFileStat.Size() - offset
			}
		}

		wg.Wait()
	}

}

func partProcess(id int, offset, size int64, sourcefileHanel, destfileHandel *os.File, wg *sync.WaitGroup, c chan int) {

	defer func() {
		// fmt.Printf("copy done for  %d \n", id)
		wg.Done()
		<-c
	}()

	buff := make([]byte, size+1)
	// fmt.Printf("id %d Writing from offset %d till %d \n", id, offset, offset+size)

	_, readErr := sourcefileHanel.ReadAt(buff, offset)
	if readErr != nil && readErr != io.EOF {
		fmt.Fprintln(os.Stderr, "failed to read chunk:", readErr.Error())
		return
	}

	_, writeErr := destfileHandel.WriteAt(buff, offset)
	if writeErr != nil {
		fmt.Fprintln(os.Stderr, "failed to write chunk:", readErr.Error())
		return
	}

	if readErr == io.EOF {
		return
	}

	return
}

func watcher(fileToWatch string, c chan int, wg *sync.WaitGroup) {

	var prev int64
	var curr int64

	defer wg.Done()

	// At starting file size is 0
	prev = 0
	for {
		// If channel is empty
		// then every routine has finished
		// work, so break watch loop ...
		if len(c) == 0 {
			return
			break
		}
		// Sleep for one second ...
		time.Sleep(time.Second * 1)
		// Read file stat again ...
		fileInfo, statErr := os.Stat(fileToWatch)
		if statErr != nil {
			fmt.Fprintln(os.Stderr, "failed to read  state of destination file ")
			return
		}
		curr = fileInfo.Size()
		speed := (curr - prev) / 1024 / 1024 // Bytes/kb/mb  Speed ... .
		prev = curr
		color.Green("copying at speed %d mbps", speed)

	}
}
