package main

import (
	"fmt"
	"image"
	"os"

	"github.com/conformal/btcec"
	"github.com/conformal/btcnet"
	"github.com/conformal/btcutil"

	"code.google.com/p/rsc/qr"
)

// Data defines the functionality needed to create
// a paper wallet.
type Data interface {
	QR() image.Image
	String() string
}

// TODO: These structures could use some pretty-printing
type privKey struct {
	qrCode *qr.Code
	value  *btcutil.WIF
}

func (pk *privKey) QR() image.Image { return pk.qrCode.Image() }
func (pk *privKey) String() string  { return fmt.Sprint(pk.value.String()) }

type addr struct {
	qrCode *qr.Code
	value  *btcutil.AddressPubKey
}

func (a *addr) QR() image.Image { return a.qrCode.Image() }
func (a *addr) String() string  { return a.value.EncodeAddress() }

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

func newPrivKeyandAddr(dump chan Data) {
	// Generate new private key
	pk, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		fmt.Println("Error generating private key:", err)
		safeClose(dump)
		return
	}
	wif, err := btcutil.NewWIF(pk, primeNet, false)
	if err != nil {
		fmt.Println("Error generating WIF:", err)
		safeClose(dump)
		return
	}
	go func(dump chan Data) {
		code, err := qr.Encode(wif.String(), qr.H)
		if err != nil {
			fmt.Println("Cannot create QR code for WIF:", err)
			safeClose(dump)
			return
		}
		dump <- &privKey{qrCode: code, value: wif}
	}(dump)

	// Extract public from private key, serialize it, and create a new pay-to-pubkey address
	addrPubKey, err := btcutil.NewAddressPubKey(pk.PubKey().SerializeUncompressed(), primeNet)
	if err != nil {
		fmt.Println("Error generating pay-to-pubkey address:", err)
		safeClose(dump)
		return
	}
	go func(dump chan Data) {
		code, err := qr.Encode(addrPubKey.EncodeAddress(), qr.H)
		if err != nil {
			fmt.Println("Cannot create QR code for pay-to-pubkey address:", err)
			safeClose(dump)
			return
		}
		dump <- &addr{qrCode: code, value: addrPubKey}
	}(dump)
}
