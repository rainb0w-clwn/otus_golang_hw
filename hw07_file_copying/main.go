package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy >= 0")
	flag.Int64Var(&offset, "offset", 0, "offset in input file >= 0")
}

func main() {
	flag.Parse()

	if from == "" {
		fmt.Println("The 'from' flag is required.")
		flag.PrintDefaults()
		return
	}

	if to == "" {
		fmt.Println("The 'to' flag is required.")
		flag.PrintDefaults()
		return
	}

	if offset < 0 {
		fmt.Println("The 'offset' flag is unsupported.")
		flag.PrintDefaults()
		return
	}

	if limit < 0 {
		fmt.Println("The 'limit' flag is unsupported.")
		flag.PrintDefaults()
		return
	}

	_, err := Copy(from, to, offset, limit)
	if err != nil {
		if errors.Is(err, ErrUnsupportedOffsetLimit) {
			fmt.Printf("Error: %s\n", err.Error())
			flag.PrintDefaults()
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
		os.Exit(1)
	}
}
