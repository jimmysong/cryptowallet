package main

import (
	"fmt"
	"os"

	flag "github.com/conformal/go-flags"
)

const (
	defaultDumpString = false
	defaultDebug      = false
)

// config defines the configuration options for xpmctl
type config struct {
	DumpString bool `long:"dump" description:"Dump WIF and pay-to-pubkey address as strings"`
	Debug      bool `long:"debug" description:"Enable debug logging"`
}

var conf = &config{
	DumpString: defaultDumpString,
	Debug:      defaultDebug,
}

func init() {
	if _, err := flag.Parse(conf); err != nil {
		usage()
		os.Exit(1)
	}
}

// usage prints out how this software has to be used
func usage() {
	fmt.Println("Usage: xpmctl [OPTION]...")
}
