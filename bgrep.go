package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func errorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, args...)
}

var searchBin []byte

func searchFile(path string) {
		fileData, err := ioutil.ReadFile(path)
		if err != nil {
			errorf("Failed to open %s: %v\n", path, err)
		} else if i := bytes.Index(fileData, searchBin); i != -1 {
			fmt.Printf("%s: found at byte %d\n", path, i)
		} else {
			fmt.Printf("%s: not found\n", path)
		}
}

func walker (path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}
	searchFile(path)
	return nil
}

func main() {
	recurse := flag.Bool("r", false, "recursivel search")
	flag.Parse()

	if flag.NArg() < 2 {
		return
	}

	var err error
	searchBin, err = hex.DecodeString(flag.Arg(0))
	if err != nil {
		errorf("Failed to parse search string: %v\n", err)
		os.Exit(1)
	}

	for _, fileName := range flag.Args()[1:] {
		if *recurse {
			err = filepath.Walk(fileName, walker)
			if err != nil {
				errorf("Walk failed: %v\n", err)
			}
		} else {
			searchFile(fileName)
		}
	}
}
