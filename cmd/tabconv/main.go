package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/nk2ge5k/tabconv"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	outdirp := flag.String("o", wd, "Output directory")
	delimp := flag.String("d", "\t", "CSV delimiter")
	flag.Parse()

	patterns := flag.Args()

	if len(patterns) == 0 {
		log.Fatal("nothing to convert")
	}

	if len(*delimp) != 1 {
		log.Fatal("delimiter must be single byte character")
	}

	outdir, err := tabconv.Expand(*outdirp)
	if err != nil {
		log.Fatal(err)
	}

	// check if output directory exists
	if exist, err := tabconv.FileExists(*outdirp); err != nil {
		log.Fatal(err)
	} else if !exist {
		log.Fatalf("directory %q does not exist", *outdirp)
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			matches = []string{pattern}
		}

		for _, file := range matches {
			if err := tabconv.Convert(file, outdir, []rune(*delimp)[0]); err != nil {
				fmt.Fprintf(os.Stdout, "failed to convert file %q: %v\n", file, err)
			}
		}
	}
}
