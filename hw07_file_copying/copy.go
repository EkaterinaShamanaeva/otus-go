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

func Copy(fromPath, toPath string, offset, limit int64, bar *pb.ProgressBar) error {
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
	switch limit {
	// copy all file
	case 0:
		for fileSize-sum > step {
			written, errCopy := io.CopyN(copiedFile, file, step)
			sum += written
			if errCopy != nil {
				if errCopy == io.EOF {
					break
				}
				return errCopy
			}
			bar.Add(int(step))
		}
		if fileSize-sum != 0 {
			step = fileSize - sum
			_, errCopy := io.CopyN(copiedFile, file, step)
			if errCopy != nil {
				if errCopy == io.EOF {
					break
				}
				return errCopy
			}
			bar.Add(int(step))
		}
		bar.Finish()
	default:
		// copy limit number of bytes in file
		for limit-sum < step {
			written, errCopy := io.CopyN(copiedFile, file, step)
			sum += written
			if errCopy != nil {
				if errCopy == io.EOF {
					break
				}
				return errCopy
			}
			bar.Add(int(step))
		}
		if limit-sum != 0 {
			step = limit - sum
			_, errCopy := io.CopyN(copiedFile, file, step)
			if errCopy != nil {
				if errCopy == io.EOF {
					break
				}
				return errCopy
			}
			bar.Add(int(step))
		}
		bar.Finish()
	}
	return nil
}
