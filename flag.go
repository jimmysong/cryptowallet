// Copyright (C) 2014-15 Michail Kargakis
// This source code is subject to the terms
// of the MIT License

package main

// TODO: Use constants as default values in flags

type config struct {
	// TODO: ImportPrivKey (sendto), sweepprivkey
	DumpString   bool   `long:"dump" description:"Dump WIF and pay-to-pubkey address as strings" default:"false"`
	Debug        bool   `long:"debug" description:"Enable debug logging" default:"false"`
	SweepPrivKey string `long:"sweepprivkey" description:"Allow sweeping of funds to another address" default:""`
}

var conf = &config{}
