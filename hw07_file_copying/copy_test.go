package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const out string = "out.txt"

func TestCopy(t *testing.T) {
	t.Run("offset exceeds file size", func(t *testing.T) {
		from := "testdata/input.txt"
		var offset int64 = 7000
		var limit int64
		require.Equal(t, ErrOffsetExceedsFileSize, Copy(from, out, offset, limit))
	})

	t.Run("file not found", func(t *testing.T) {
		from := "testdata/input_invalid.txt"
		var offset int64 = 6000
		var limit int64
		require.Equal(t, ErrOpenFile, Copy(from, out, offset, limit))
	})

	t.Run("copy with limit 10", func(t *testing.T) {
		from := "testdata/input.txt"
		var offset int64
		var limit int64 = 10
		require.Equal(t, nil, Copy(from, out, offset, limit))

		var file *os.File
		file, _ = os.OpenFile(to, os.O_RDONLY, 0)
		defer file.Close()

		fileInfo, _ := file.Stat()
		require.Equal(t, int64(10), fileInfo.Size())
		os.Remove(out)
	})

	t.Run("copy with hieroglyphs", func(t *testing.T) {
		from := "testdata/input2.txt"
		var offset int64
		var limit int64
		require.Equal(t, nil, Copy(from, out, offset, limit))

		var file *os.File
		file, _ = os.OpenFile(to, os.O_RDONLY, 0)
		defer file.Close()

		fileInfo, _ := file.Stat()
		require.Equal(t, int64(18), fileInfo.Size())
		os.Remove(out)
	})
}
