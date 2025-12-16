package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedOffsetLimit = errors.New("unsupported offset or limit")
	ErrFromFileNotExists      = errors.New("input file not exists")
	ErrFromFileOpen           = errors.New("input file opening error")
	ErrFromFileUnsupported    = errors.New("input file not supported")
	ErrOffsetExceedsFileSize  = errors.New("offset exceeds file size")
	ErrToFileDirNotExists     = errors.New("out file dir not exists")
	ErrToFileOpen             = errors.New("out file opening error")
	ErrToFileWrite            = errors.New("out file writing error")
)

func Copy(fromPath, toPath string, offset, limit int64) (int64, error) {
	if offset < 0 || limit < 0 {
		return 0, ErrUnsupportedOffsetLimit
	}

	fromFile, err := os.Open(fromPath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, ErrFromFileNotExists
		}
		return 0, ErrFromFileOpen
	}
	defer fromFile.Close()

	fromFileStat, err := fromFile.Stat()
	switch {
	case err != nil:
		return 0, err
	case fromFileStat.IsDir():
		return 0, ErrFromFileUnsupported
	}

	fromFileRegular := fromFileStat.Mode().IsRegular()
	fromFileSize := fromFileStat.Size()

	if fromFileRegular {
		if offset > fromFileSize {
			return 0, ErrOffsetExceedsFileSize
		}
	} else if limit == 0 {
		return 0, ErrFromFileUnsupported
	}

	toFile, err := os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, ErrToFileDirNotExists
		}
		return 0, ErrToFileOpen
	}
	defer toFile.Close()

	reader := io.Reader(fromFile)
	if offset > 0 {
		if _, err = fromFile.Seek(offset, io.SeekStart); err != nil {
			return 0, err
		}
	}

	bytesToRead := fromFileSize - offset

	if limit > 0 && (!fromFileRegular || fromFileSize >= offset+limit) {
		bytesToRead = limit
	}

	bar := pb.Full.Start64(bytesToRead)
	defer bar.Finish()

	barReader := bar.NewProxyReader(reader)

	return io.CopyN(toFile, barReader, bytesToRead)
}
