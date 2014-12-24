package main

import (
	"fmt"
	"os"

	"github.com/conformal/btcec"
	"github.com/conformal/btcnet"
	"github.com/conformal/btcutil"

	"code.google.com/p/rsc/qr"
)

var primeNet = &btcnet.Params{
	Name:             "Primecoin",
	PubKeyHashAddrID: 23,
	ScriptHashAddrID: 53,
	HDCoinType:       23,
}

func init() {
	if err := btcnet.Register(primeNet); err != nil {
		fmt.Println("Couldn't register Primecoin network parameters")
		os.Exit(1)
	}
}

func main() {
	ch := make(chan string, 1)
	qrChan := make(chan *qr.Code, 2)
	go newPrivKeyandAddr(ch)
	go qrCodes(ch, qrChan)

	wifCode := <-qrChan
	pkCode := <-qrChan

	fmt.Println(wifCode.Image(), pkCode.Image())
}

func newPrivKeyandAddr(ch chan string) {
	// Generate new private key
	pk, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		fmt.Println("Error generating private key: " + err.Error())
		close(ch)
		return
	}
	wif, err := btcutil.NewWIF(pk, primeNet, false)
	if err != nil {
		fmt.Println("Error generating WIF: " + err.Error())
		close(ch)
		return
	}
	ch <- wif.String()

	// Extract public from private key, serialize it, and create a new pay-to-pubkey address
	addrPubKey, err := btcutil.NewAddressPubKey(pk.PubKey().SerializeUncompressed(), primeNet)
	if err != nil {
		fmt.Println("Error generating pay-to-pubkey address: " + err.Error())
		close(ch)
		return
	}
	ch <- addrPubKey.EncodeAddress()
}

func qrCodes(ch chan string, qrChan chan *qr.Code) {
	for i := 0; i < 2; i++ {
		select {
		case str := <-ch:
			if str == "" {
				return
			}
			code, err := qr.Encode(str, qr.H)
			if err != nil {
				fmt.Println(err)
				return
			}
			qrChan <- code
		}
	}
}
