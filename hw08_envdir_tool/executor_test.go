package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("run command", func(t *testing.T) {
		path, _ := os.Getwd()
		cmd := []string{
			"/bin/bash", filepath.Join(path, "/testdata/echo.sh"),
			"arg1=1", "arg2=2",
		}

		envMap := make(Environment)

		code := RunCmd(cmd, envMap)

		require.Equal(t, exitCodeSuccessful, code)
	})

	t.Run("incorrect environmental variable value", func(t *testing.T) {
		path, _ := os.Getwd()
		cmd := []string{
			"/bin/bash", filepath.Join(path, "/testdata/echo.sh"),
			"arg1=1", "arg2=2",
		}

		envMap := make(Environment)
		envMap["INCORRECT_VAR"] = EnvValue{Value: "foo\x00with new line", NeedRemove: false}

		code := RunCmd(cmd, envMap)

		require.Equal(t, exitCodeUnsuccessful, code)
	})

	t.Run("wrong path to the directory with the command", func(t *testing.T) {
		path, _ := os.Getwd()
		cmd := []string{
			"/bin/bash", filepath.Join(path, "/testdata/echoIncorrect.sh"),
			"arg1=1", "arg2=2",
		}

		envMap := make(Environment)

		code := RunCmd(cmd, envMap)

		require.Equal(t, exitCodeUnsuccessful, code)
	})

	t.Run("delete environment variable", func(t *testing.T) {
		path, _ := os.Getwd()
		cmd := []string{
			"/bin/bash", filepath.Join(path, "/testdata/echo.sh"),
			"arg1=1", "arg2=2",
		}

		_ = os.Setenv("EMPTY", "should be removed")

		envMap := make(Environment)
		envMap["EMPTY"] = EnvValue{Value: "", NeedRemove: true}

		code := RunCmd(cmd, envMap)

		_, ok := os.LookupEnv("EMPTY")

		require.Equal(t, false, ok)
		require.Equal(t, exitCodeSuccessful, code)
	})
}
