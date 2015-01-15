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

var netParams = &btcnet.Params{}

func init() {
	_, err := flag.Parse(conf)
	debug(err, "Error while parsing flags")

	if id, supported := coinID[conf.CoinType]; supported {
		if conf.Testnet {
			netParams.PubKeyHashAddrID = id.isOnTestNet()
		} else {
			netParams.PubKeyHashAddrID = id.isOnMainNet()
		}
	} else {
		fmt.Println("Coin type " + conf.CoinType + " not supported!")
		os.Exit(1)
	}
	netParams.PrivateKeyID = netParams.PubKeyHashAddrID + 128
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
func debug(err error, reason string) {
	if err == nil {
		return
	}
	fmt.Println(reason+":", err)
	os.Exit(1)
}
