// Copyright (C) 2014-15 Michail Kargakis
// This source code is subject to the terms
// of the MIT License

package main

const (
	defaultDumpString = false
	defaultDebug      = false
	defaultTestnet    = false
	defaultCoinType   = "btc"
	defaultSupport    = false
)

type config struct {
	DumpString bool   `long:"dump" description:"Dump WIF and pay-to-pubkey address as strings"`
	Debug      bool   `long:"debug" description:"Enable debug logging"`
	Testnet    bool   `long:"testnet" description:"Testnet network"`
	CoinType   string `long:"coin" description:"Coin type"`
	Support    bool   `long:"support" description:"Show supported cryptocurrencies"`
}

var conf = &config{
	DumpString: defaultDumpString,
	Debug:      defaultDebug,
	Testnet:    defaultTestnet,
	CoinType:   defaultCoinType,
	Support:    defaultSupport,
}
