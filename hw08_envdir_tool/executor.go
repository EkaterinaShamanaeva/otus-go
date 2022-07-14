package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {

	envVariables := make([]string, len(env))
	index := 0
	for key, value := range env {
		if value.NeedRemove != true {
			envVariables[index] = fmt.Sprintf("%s=%s", key, value.Value)
			index++
		}
	}

	childCommand := cmd[2]
	fmt.Println(childCommand)
	childCommandParams := cmd[3:]
	command := exec.Command(childCommand, childCommandParams...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = envVariables
	if err := command.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus := exitError.Sys().(syscall.WaitStatus)
			os.Exit(waitStatus.ExitStatus())
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(111)
		}
	}
	os.Exit(0)
	return
}
