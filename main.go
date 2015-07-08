package main

import (
	"bufio"
	"flag"
	"os"

	"fmt"

	f "github.com/stricture/fmt-validator/fmt"
)

func croak(m string) {
	fmt.Println(m)
	os.Exit(1)
}

const usage = `
  Usage: fmt-validate [options] <hash> <pattern>

  Options:
  -h, --help          Show usage and exit.
  -a                  Lookup pattern for algorithm by name.
  -i                  Match each line in a file against a pattern
`

func main() {
	name := flag.String("a", "", "Lookup pattern for algorithm by name")
	fname := flag.String("i", "", "Match each line in a file against a pattern")
	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	if *name == "" && *fname == "" && len(flag.Args()) < 2 {
		croak("Missing required argument")
	}

	var hash string
	var pattern string
	if len(flag.Args()) < 2 {
		hash = flag.Arg(0)
		pattern = flag.Arg(1)
		ok, err := f.T(hash, pattern)
		if err != nil {
			croak("Invalid pattern")
		}
		if ok != "" {
			fmt.Println(ok)
		}
	}
	if *name != "" {
		format, ok := f.Formats()[*name]
		if !ok {
			croak("Algorithm not found")
		}
		pattern = format.Pattern
	} else {
		pattern = flag.Arg(0)
	}
	if *fname != "" {
		fh, err := os.Open(*fname)
		if err != nil {
			croak("Could not open file")
		}
		defer fh.Close()
		scanner := bufio.NewScanner(fh)
		for scanner.Scan() {
			ok, err := f.T(scanner.Text(), pattern)
			if err != nil {
				croak("Invalid pattern")
			}
			if ok != "" {
				fmt.Println(ok)
			}
		}
	} else {
		hash = flag.Arg(0)
		ok, _ := f.T(hash, pattern)
		if ok != "" {
			fmt.Println(ok)
		}
	}
}
