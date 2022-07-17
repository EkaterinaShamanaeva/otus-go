package main

import (
	"bufio"
	"errors"
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

var ErrUnsupportedFile = errors.New("file name contains '='")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// create map of env variables
	envMap := make(Environment)
	// read directory
	files, errReadDir := ioutil.ReadDir(dir)
	if errReadDir != nil {
		return nil, errReadDir
	}
	// read files in the directory
	for _, file := range files {
		// checks
		if file.IsDir() {
			continue
		}

		if strings.IndexByte(file.Name(), '=') != -1 {
			return nil, ErrUnsupportedFile
		}

		if file.Size() == 0 {
			envMap[file.Name()] = EnvValue{Value: "", NeedRemove: true}
			continue
		}

		// read the first line
		openedFile, _ := os.Open(filepath.Join(dir, file.Name()))
		text, _, errReadFile := bufio.NewReader(openedFile).ReadLine()
		if errReadFile != nil {
			return nil, errReadFile
		}
		_ = openedFile.Close()
		fileString := string(text)

		if strings.IndexByte(fileString, '\x00') != -1 {
			fileString = strings.ReplaceAll(fileString, "\x00", "\n")
		}

		fileString = strings.TrimRight(fileString, " ")
		fileString = strings.TrimRight(fileString, "\n")

		if len(fileString) == 0 {
			envMap[file.Name()] = EnvValue{Value: "", NeedRemove: false}
			continue
		}

		envMap[file.Name()] = EnvValue{Value: fileString, NeedRemove: false}
	}

	return envMap, nil
}
