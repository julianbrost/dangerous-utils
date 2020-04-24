package main

import (
	"compress/gzip"
	"io"
	"os"
	"syscall"
)

type Args struct {
	Source      string
	Destination string
	Final       bool
	Max         uint64
}

// TODO: use constants from C header files?
const (
	SeekData          = 3
	FallocFlKeepSize  = 0x01
	FallocFlPunchHole = 0x02
)

func DestructiveGzip(args *Args) error {
	sourceFile, err := os.OpenFile(args.Source, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	destinationFile, err := os.OpenFile(args.Destination, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	var sourceStat syscall.Stat_t
	err = syscall.Fstat(int(sourceFile.Fd()), &sourceStat)
	if err != nil {
		return err
	}
	blockSize := sourceStat.Blksize
	sourceSize := sourceStat.Size

	sourceStart, err := sourceFile.Seek(0, SeekData)
	sourceEnd := sourceSize
	if args.Max > 0 && sourceEnd-sourceStart > int64(args.Max) {
		sourceEnd = sourceStart + int64(args.Max)
	}
	if !args.Final {
		sourceEnd = sourceEnd / blockSize * blockSize
	}

	destinationWriter := gzip.NewWriter(destinationFile)
	_, err = io.CopyN(destinationWriter, sourceFile, sourceEnd-sourceStart)
	if err != nil {
		return err
	}
	err = destinationWriter.Close()
	if err != nil {
		return err
	}
	err = destinationFile.Close()
	if err != nil {
		return err
	}

	err = syscall.Fallocate(int(sourceFile.Fd()), FallocFlKeepSize|FallocFlPunchHole, 0, sourceEnd)
	if err != nil {
		return err
	}

	return nil
}
