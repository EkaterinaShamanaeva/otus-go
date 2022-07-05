package main

import (
	"errors"
	"github.com/cheggaaa/pb/v3"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOpenFile              = errors.New("file was not opened")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
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

	var sum int64 = 0
	var step int64 = 4

	// copy
	var bytesToCopy int64
	if limit == 0 {
		bytesToCopy = fileSize
	} else {
		bytesToCopy = limit
	}

	var count int
	if fileSize-offset-limit < 0 {
		count = int(fileSize - offset)
	} else if limit == 0 {
		count = int(fileSize)
	} else {
		count = int(limit)
	}
	bar := pb.StartNew(count)

	for bytesToCopy-sum > step {
		written, errCopy := io.CopyN(copiedFile, file, step)
		sum += written
		bar.Add(int(step))
		//time.Sleep(time.Millisecond)
		if errCopy != nil {
			if errCopy == io.EOF {
				break
			}
			return errCopy
		}
	}
	if bytesToCopy-sum != 0 {
		step = bytesToCopy - sum
		_, errCopy := io.CopyN(copiedFile, file, step)
		bar.Add(int(step))
		//time.Sleep(time.Millisecond)
		if errCopy != nil {
			if errCopy == io.EOF {
				//break
				return nil
			}
			return errCopy
		}
	}
	bar.Finish()
	return nil
}
