package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("wrong path to the directory", func(t *testing.T) {
		path, _ := os.Getwd()
		pathToDir := filepath.Join(path, "/testdata/envIncorrect")

		envMap, errReadDir := ReadDir(pathToDir)

		require.Empty(t, envMap)
		require.NotEqual(t, nil, errReadDir)
	})

	t.Run("case with environment variables", func(t *testing.T) {
		path, _ := os.Getwd()
		pathToDir := filepath.Join(path, "/testdata/env")

		envMap, errReadDir := ReadDir(pathToDir)

		expectedEnvMap := make(Environment)
		expectedEnvMap["BAR"] = EnvValue{Value: "bar", NeedRemove: false}
		expectedEnvMap["EMPTY"] = EnvValue{Value: "", NeedRemove: false}
		expectedEnvMap["FOO"] = EnvValue{Value: "   foo\nwith new line", NeedRemove: false}
		expectedEnvMap["HELLO"] = EnvValue{Value: `"hello"`, NeedRemove: false}
		expectedEnvMap["UNSET"] = EnvValue{Value: "", NeedRemove: true}

		require.Equal(t, nil, errReadDir)
		require.Equal(t, expectedEnvMap["BAR"], envMap["BAR"])
		require.Equal(t, expectedEnvMap["EMPTY"], envMap["EMPTY"])
		require.Equal(t, expectedEnvMap["FOO"], envMap["FOO"])
		require.Equal(t, expectedEnvMap["HELLO"], envMap["HELLO"])
		require.Equal(t, expectedEnvMap["UNSET"], envMap["UNSET"])
	})
}
