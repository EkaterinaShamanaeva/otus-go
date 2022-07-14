package main

import (
	"fmt"
	"log"
	"os"
)

type Args struct {
	pathToEnvDir string
	command      string
	arg          []string
}

func main() {
	args := os.Args //getArgs()
	fmt.Println(args)
	envMap, _ := ReadDir(args[1])
	RunCmd(args, envMap)
	//fmt.Println(args)
	//envMap := make(Environment)
	//envMap, _ = ReadDir("/Users/ekaterina/GolandProjects/otus-go/hw08_envdir_tool/testdata/env")
	//fmt.Println(envMap)

}

func getArgs() *Args {
	args := Args{
		pathToEnvDir: os.Args[1],
		command:      os.Args[2],
		arg:          os.Args[3:],
	}

	if args.pathToEnvDir == "" {
		log.Fatal("Specify the path to the folder with environment variables")
	}
	if args.command == "" {
		log.Fatal("Specify the command to execute")
	}

	return &args
}
