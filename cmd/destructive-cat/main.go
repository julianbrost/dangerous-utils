package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var blockSize int64

	flag.Int64Var(&blockSize, "block-size", DefaultBlockSize, "block size used for hole punching")
	flag.Parse()

	err := DestructiveCat(flag.Args(), os.Stdout, blockSize)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
