package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envMap := make(Environment)

	// read directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(111)
	}
	// read files in the directory
	for _, file := range files {

		if file.IsDir() {
			continue
		}
		if file.Size() == 0 {
			envMap[file.Name()] = EnvValue{Value: "", NeedRemove: true}
			continue
		}
		openedFile, _ := os.Open(filepath.Join(dir, file.Name()))

		text, _, errReadFile := bufio.NewReader(openedFile).ReadLine()
		if errReadFile != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(111)
		}

		_ = openedFile.Close()
		fileString := string(text)

		if stringHasNullCharacter(fileString) {
			fileString = strings.Replace(fileString, "\\x00", "\n", -1)
		}

		fileString = strings.TrimRight(fileString, "\t\n")

		if len(fileString) == 0 {
			envMap[file.Name()] = EnvValue{Value: "", NeedRemove: false}
			continue
		}

		envMap[file.Name()] = EnvValue{Value: fileString, NeedRemove: false}
	}

	return envMap, nil
}

func stringHasNullCharacter(s string) bool {
	i := strings.IndexByte(s, '\x00')
	return i != -1
}
