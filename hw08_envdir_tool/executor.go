package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	childCommand := cmd[3]
	childCommandParams := cmd[4:]

	command := exec.Command(childCommand, childCommandParams...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	//fmt.Println("empty: ", len(env["EMPTY"].Value), env["EMPTY"].Value+"test")

	for key, value := range env {
		if value.NeedRemove != true {
			if _, ok := os.LookupEnv(key); ok {
				_ = os.Unsetenv(key)
				err := os.Setenv(key, value.Value)
				if err != nil {
					fmt.Println("error", err, key, value.Value, len(value.Value))
				}
			}
			_ = os.Setenv(key, value.Value)
		} else {
			if _, ok := os.LookupEnv(key); ok {
				_ = os.Unsetenv(key)
			}
		}
	}
	//fmt.Println(envVariables, len(envVariables))
	command.Env = os.Environ()
	//fmt.Println("env: ", command.Env)

	err := command.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	return 0
}
