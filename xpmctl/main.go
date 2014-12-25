// Copyright (C) 2014-15 Michail Kargakis
// This source code is subject to the terms
// of the MIT License
package main

import (
	"fmt"
	"os"
	"path/filepath"

	pdf "code.google.com/p/gofpdf"
)

func main() {
	dump := make(chan Data, 1)
	// Start a goroutine to generate a private key and
	// a pay-to-pubkey address.
	go newPrivKeyandAddr(dump)
	wif := <-dump
	addrPubKey := <-dump
	if wif == nil || addrPubKey == nil {
		fmt.Println("Corrupted data")
		return
	}

	// Create pdf
	paperWallet := pdf.New("P", "mm", "A4", "")
	paperWallet.AddPage()
	paperWallet.SetFont("Arial", "B", 8)
	content := fmt.Sprintf("PrivKey: %s Address: %s\n", wif.String(), addrPubKey.String())
	paperWallet.Cell(40, 10, content)
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Cannot get current working directory:", err)
	}
	fileStr := filepath.Join(dir, "wallet.pdf")

	if err := paperWallet.OutputFileAndClose(fileStr); err == nil {
		fmt.Println("Successfully generated wallet.pdf")
	} else {
		fmt.Println(err)
	}
}

func safeClose(ch chan Data) {
	select {
	case <-ch:
	default:
		close(ch)
	}
}
