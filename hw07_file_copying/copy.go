package main

import (
	"errors"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOpenFile              = errors.New("file was not opened")
)

func Copy(fromPath, toPath string, offset, limit int64, pb *pb.ProgressBar) error {
	// open file <from>
	var file *os.File
	file, err := os.OpenFile(fromPath, os.O_RDONLY, 0)
	defer file.Close()
	if err != nil {
		return ErrOpenFile
	}

	// count of bytes in file <from>
	fileInfo, errFileSize := file.Stat()
	if errFileSize != nil {
		return ErrUnsupportedFile
	}
	fileSize := fileInfo.Size()
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	// create output file
	copiedFile, _ := os.Create(toPath)
	defer copiedFile.Close()

	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// copy
	switch limit {
	case 0:
		written, errCopy := io.Copy(copiedFile, file)
		if errCopy != nil {
			return errCopy
		}
		fmt.Println(written)
	default:
		written, errCopyN := io.CopyN(copiedFile, file, limit)
		if errCopyN != nil {
			return errCopyN
		}
		fmt.Println(written)
	}
	pb.Increment()
	return nil
}
