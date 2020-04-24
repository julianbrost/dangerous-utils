package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var args Args
	flag.StringVar(&args.Source, "source", "", "source file")
	flag.StringVar(&args.Destination, "destination", "", "destination file")
	flag.BoolVar(&args.Final, "final", false, "final pass (i.e. copy end of file smaller than a filesystem block")
	flag.Uint64Var(&args.Max, "max-size", 0, "maximum size to copy at once")
	flag.Parse()
	if args.Source == "" {
		log.Fatal("-source must be given")
	}
	if args.Destination == "" {
		log.Fatal("-destination must be given")
	}

	err := DestructiveGzip(&args)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
