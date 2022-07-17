package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	exitCodeSuccessful   = 0
	exitCodeUnsuccessful = 111
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// command
	childCommand := cmd[0]
	// arguments
	childCommandArgs := cmd[1:]

	command := exec.Command(childCommand, childCommandArgs...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	// update or create environment variables
	for key, value := range env {
		if !value.NeedRemove {
			if _, ok := os.LookupEnv(key); ok {
				_ = os.Unsetenv(key)
			}
			err := os.Setenv(key, value.Value)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return exitCodeUnsuccessful
			}
		} else if _, ok := os.LookupEnv(key); ok {
			_ = os.Unsetenv(key)
		}
	}

	command.Env = os.Environ()

	if err := command.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitCodeUnsuccessful
	}

	return exitCodeSuccessful
}
