// Copyright (C) 2014-15 Michail Kargakis
// This source code is subject to the terms
// of the MIT License

package main

import (
	"fmt"
	"os"

	flag "github.com/conformal/go-flags"
)

func init() {
	if _, err := flag.Parse(conf); err != nil {
		usage()
		os.Exit(1)
	}
}

func main() {
	// TODO: Search for existing paper wallet in the
	// current working directory. If so, abort.
	pk, addr := NewPrivKeyAndAddr()

	if !conf.DumpString {
		NewPaperWallet(pk, addr)
	}
	if conf.DumpString {
		fmt.Println(pk)
		fmt.Println(addr)
	}
}

// debug is a conveniece function for handling errors;
// Since every error in the process is regarded as fatal,
// this function can be used everywhere in this tool.
func debug(err error) {
	if err == nil {
		return
	}
	fmt.Println(err)
	os.Exit(1)
}
