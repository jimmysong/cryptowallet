// Copyright (C) 2014-15 Michail Kargakis
// This source code is subject to the terms
// of the MIT License

package main

// ID is a struct containing ids of each coin
// for both mainnet and testnet networks.
type ID struct {
	mainNet uint8
	testNet uint8
}

func (id *ID) isOnMainNet() uint8 {
	return id.mainNet
}

func (id *ID) isOnTestNet() uint8 {
	return id.testNet
}

var coinID = map[string]*ID{
	"btc":  &ID{0, 111},
	"ltc":  &ID{48, 111},
	"xpm":  &ID{23, 112},
	"ppc":  &ID{55, 112},
	"nmc":  &ID{53, 112},
	"doge": &ID{30, 122},
}
