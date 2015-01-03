// Copyright (C) 2014-15 Michail Kargakis
// This source code is subject to the terms
// of the MIT License

package main

import (
	"fmt"
	"os"

	"github.com/conformal/btcnet"
	flag "github.com/conformal/go-flags"
)

var primeNet = &btcnet.Params{
	Name:             "Primecoin",
	PubKeyHashAddrID: 23,
	ScriptHashAddrID: 53,
	HDCoinType:       23,
}

func init() {
	if _, err := flag.Parse(conf); err != nil {
		usage()
		os.Exit(1)
	}
}

func main() {
	// Print out usage message if --help is one of the flags
	if conf.Help {
		usage()
		os.Exit(1)
	}
	// TODO: Search for existing paper wallet in the
	// current working directory. If so, abort.
	pk := NewPrivKey()
	addr := NewAddress(pk.value)

	if !conf.DumpString {
		NewPaperWallet(pk, addr)
	}
	if conf.DumpString {
		fmt.Println(pk)
		fmt.Println(addr)
	}
}

// debug is a conveniece function for handling errors.
// TODO: This function has to accept a string as well
// for more informative logging. The functionality is
// going to be added as soon as the --debug flag will
// be fully implemented.
func debug(err error) {
	if err == nil {
		return
	}
	fmt.Println(err)
	os.Exit(1)
}
