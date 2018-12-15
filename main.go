// Commandline application to
// superfast copy file to another
// location

package main

import (
	"fmt"
	"os"
)

func main() {
	// Checking for input argumnets
	inputs := os.Args
	if len(inputs) < 3 {
		fmt.Fprintln(os.Stdout, "please, provide Input file and destination file path")
		os.Exit(1)
		return
	}
	// First argument is path to the input file
	inputPath := inputs[1]
	// Second argument is path to hte destination file
	destinationPath := inputs[2]

	fmt.Fprintln(os.Stdout, "Input provided by you are :")
	fmt.Fprintln(os.Stdout, inputPath)
	fmt.Fprintln(os.Stdout, destinationPath)

	copy(inputPath, destinationPath)

	fmt.Fprintln(os.Stdout, "Process completed succesfully")
	os.Exit(0)
}
