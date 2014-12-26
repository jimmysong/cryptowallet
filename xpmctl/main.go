// Copyright (C) 2014-15 Michail Kargakis
// This source code is subject to the terms
// of the MIT License
package main

import (
	"fmt"
)

func main() {
	dump := make(chan Data, 1)
	// Start a goroutine to generate a private key and
	// a pay-to-pubkey address.
	go newPrivKeyandAddr(dump)
	wif := <-dump
	addr := <-dump
	if wif == nil || addr == nil {
		fmt.Println("Corrupted data")
		return
	}

	newPaperWallet(wif, addr)
}
