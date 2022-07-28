package main

import (
	"fmt"
	"os"
)

func main() {
	// parse dir, command and arguments
	args := os.Args
	pathToDir := args[1]
	commandAndParams := args[2:]
	// return map of environment variables
	envMap, errReadDir := ReadDir(pathToDir)
	if errReadDir == nil {
		// run command
		code := RunCmd(commandAndParams, envMap)
		os.Exit(code)
	} else {
		fmt.Fprintln(os.Stderr, errReadDir)
		os.Exit(exitCodeUnsuccessful)
	}
}
