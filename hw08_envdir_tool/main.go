package main

import (
	"os"
)

func main() {
	args := os.Args //getArgs()
	//for i, elem := range args {
	//fmt.Println("ARGUMENTS: ", i, elem)
	//}

	envMap, _ := ReadDir(args[1])
	code := RunCmd(args, envMap)
	os.Exit(code)
	//fmt.Println(args)
	//envMap := make(Environment)
	//envMap, _ = ReadDir("/Users/ekaterina/GolandProjects/otus-go/hw08_envdir_tool/testdata/env")
	//fmt.Println(envMap)

}
