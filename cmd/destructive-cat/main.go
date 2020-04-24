package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	var blockSize int64

	flag.Int64Var(&blockSize, "block-size", DefaultBlockSize, "block size used for hole punching")
	flag.Parse()

	if flag.NArg() > 0 {
		files := flag.Args()

		fmt.Fprintf(os.Stderr, "You are about to run destructive-cat on the follwing files:\n")
		for _, file := range files {
			fmt.Fprintf(os.Stderr, "  * %q\n", file)
		}
		fmt.Fprintf(os.Stderr, "These files will be DELETED. Are you sure? Type YES to continue: ")
		response, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		if response != "YES\n" {
			fmt.Fprintf(os.Stderr, "\nYou did not type YES, bailing out.\n")
			os.Exit(1)
		}

		err := DestructiveCat(files, os.Stdout, blockSize)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
