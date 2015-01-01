// Copyright (C) 2014-15 Michail Kargakis
// This source code is subject to the terms
// of the MIT License

package main

import (
	"fmt"

	_ "github.com/conformal/go-flags"
)

type config struct {
	// TODO: ImportPrivKey (sendto), sweepprivkey
	DumpString bool `long:"dump" description:"Dump WIF and pay-to-pubkey address as strings" default:"false"`
	Debug      bool `long:"debug" description:"Enable debug logging" default:"false"`
}

var conf = &config{}

// usage prints out how this software has to be used
func usage() {
	fmt.Println("Usage: xpmwallet [OPTION]...")
}
