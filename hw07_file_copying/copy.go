package main

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
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
	if err != nil {
		return ErrOpenFile
	}
	defer file.Close()

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

	// move pointer according to offset
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// define the number of bytes to copy
	bytesToCopy := defineBytesToCopy(fileSize, offset, limit)

	// new progress bar
	bar := pb.StartNew(int(bytesToCopy))

	// copy
	var sum int64
	var buffer int64 = 4

	for bytesToCopy-sum > buffer {
		written, errCopy := io.CopyN(copiedFile, file, buffer)
		sum += written
		bar.Add(int(buffer))
		time.Sleep(time.Millisecond)
		if errCopy != nil {
			if errors.Is(errCopy, io.EOF) {
				break
			}
			return errCopy
		}
	}
	if bytesToCopy-sum != 0 {
		buffer = bytesToCopy - sum
		_, errCopy := io.CopyN(copiedFile, file, buffer)
		bar.Add(int(buffer))
		time.Sleep(time.Millisecond)
		if errCopy != nil {
			if errors.Is(errCopy, io.EOF) {
				return nil
			}
			return errCopy
		}
	}
	bar.Finish()
	return nil
}

func defineBytesToCopy(fileSize, offset, limit int64) int64 {
	var bytesToCopy int64
	switch {
	case limit == 0 && offset == 0:
		bytesToCopy = fileSize
	case limit == 0 && offset != 0:
		bytesToCopy = fileSize - offset
	case limit != 0 && offset == 0:
		if limit > fileSize {
			bytesToCopy = fileSize
		} else {
			bytesToCopy = limit
		}
	case limit != 0 && offset != 0:
		if limit > fileSize-offset {
			bytesToCopy = fileSize - offset
		} else {
			bytesToCopy = limit
		}
	}
	return bytesToCopy
}
