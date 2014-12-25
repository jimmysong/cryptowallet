package main

import (
	"fmt"

	pdf "code.google.com/p/gofpdf"
)

func main() {
	dump := make(chan Data, 1)
	go newPrivKeyandAddr(dump)

	wif := <-dump
	addrPubKey := <-dump
	if wif == nil || addrPubKey == nil {
		fmt.Println("Corrupted data")
		return
	}

	_ = addrPubKey
	_ = wif

	// TODO: Create pdf
	_ = pdf.New("P", "mm", "A4", "")
}

func safeClose(ch chan Data) {
	select {
	case <-ch:
	default:
		close(ch)
	}
}
