package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const outFile string = "out.txt"

func TestCopy(t *testing.T) {
	t.Run("offset exceeds file size", func(t *testing.T) {
		fromFile := "testdata/input.txt"
		var offsetFile int64 = 7000
		var limitFile int64
		require.Equal(t, ErrOffsetExceedsFileSize, Copy(fromFile, outFile, offsetFile, limitFile))
	})

	t.Run("file not found", func(t *testing.T) {
		fromFile := "testdata/input_invalid.txt"
		var offsetFile int64 = 6000
		var limitFile int64
		require.Equal(t, ErrOpenFile, Copy(fromFile, outFile, offsetFile, limitFile))
	})

	t.Run("copy with limit 10", func(t *testing.T) {
		fromFile := "testdata/input.txt"
		var offsetFile int64
		var limitFile int64 = 10
		require.Equal(t, nil, Copy(fromFile, outFile, offsetFile, limitFile))

		// open out.txt
		var file *os.File
		file, _ = os.OpenFile(outFile, os.O_RDONLY, 0)

		fileInfo, _ := file.Stat()
		require.Equal(t, int64(10), fileInfo.Size())
		_ = file.Close()
		_ = os.Remove(outFile)
	})

	t.Run("copy with hieroglyphs", func(t *testing.T) {
		fromFile := "testdata/input2.txt"
		var offsetFile int64
		var limitFile int64
		require.Equal(t, nil, Copy(fromFile, outFile, offsetFile, limitFile))

		var file *os.File
		file, _ = os.OpenFile(outFile, os.O_RDONLY, 0)

		fileInfo, _ := file.Stat()
		require.Equal(t, int64(18), fileInfo.Size())
		_ = file.Close()
		_ = os.Remove(outFile)
	})
}
