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

var netParams = &btcnet.Params{
	PubKeyHashAddrID: 23,
	PrivateKeyID:     151,
}

func init() {
	_, err := flag.Parse(conf)
	debug(err)
	if conf.Testnet {
		netParams.PubKeyHashAddrID = 112
		netParams.PrivateKeyID = 240
	}
}

func main() {
	pk := NewPrivKey()
	if !conf.DumpString {
		NewPaperWallet(pk)
	}
	if conf.DumpString {
		fmt.Println(pk)
		fmt.Println(NewAddress(pk.value))
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
