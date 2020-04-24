package main

import (
	"io"
	"os"
	"syscall"
)

const (
	DefaultBlockSize = 4 * 1024 * 1024

	// TODO: use constants from C header files?
	fallocFlKeepSize  = 0x01
	fallocFlPunchHole = 0x02
)

func DestructiveCat(inputFileNames []string, output io.WriteCloser, blockSize int64) error {
	for _, inputFileName := range inputFileNames {
		inputFile, err := os.OpenFile(inputFileName, os.O_RDWR, 0)
		if err != nil {
			return err
		}
		for off := int64(0); ; off += blockSize {
			_, err = io.CopyN(output, inputFile, blockSize)
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}
			err = syscall.Fallocate(int(inputFile.Fd()), fallocFlKeepSize|fallocFlPunchHole, off, blockSize)
			if err != nil {
				return err
			}
		}
		err = os.Remove(inputFileName)
		if err != nil {
			return err
		}
		err = inputFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
