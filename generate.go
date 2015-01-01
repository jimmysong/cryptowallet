// Copyright (C) 2014-15 Michail Kargakis
// This source code is subject to the terms
// of the MIT License

package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/conformal/btcec"
	"github.com/conformal/btcnet"
	"github.com/conformal/btcutil"

	pdf "code.google.com/p/gofpdf"
	"code.google.com/p/rsc/qr"
)

type privKey struct {
	qrCode *qr.Code
	value  *btcutil.WIF
}

func (pk *privKey) QR() image.Image { return pk.qrCode.Image() }
func (pk *privKey) String() string  { return fmt.Sprint(pk.value.String()) }

type addrPubKey struct {
	qrCode *qr.Code
	value  *btcutil.AddressPubKey
}

func (a *addrPubKey) QR() image.Image { return a.qrCode.Image() }
func (a *addrPubKey) String() string  { return a.value.EncodeAddress() }

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

// NewPrivKeyAndAddr returns a new private key and a corresponding
// public address. If any error occurs during the process, xpmwallet
// exits with exit status one.
func NewPrivKeyAndAddr() (*privKey, *addrPubKey) {
	// Generate new private key
	pk, err := btcec.NewPrivateKey(btcec.S256())
	debug(err)
	wif, err := btcutil.NewWIF(pk, primeNet, false)
	debug(err)
	pkCode, err := qr.Encode(wif.String(), qr.H)
	debug(err)
	// Extract public from private key, serialize it, and create a new pay-to-pubkey address
	addr, err := btcutil.NewAddressPubKey(pk.PubKey().SerializeUncompressed(), primeNet)
	debug(err)
	addrCode, err := qr.Encode(addr.EncodeAddress(), qr.H)
	debug(err)
	return &privKey{qrCode: pkCode, value: wif}, &addrPubKey{qrCode: addrCode, value: addr}
}

// NewPaperWallet accepts a private key and a public address
// (presumably an address corresponding to the private key) and
// generates a pdf paper wallet.
func NewPaperWallet(pk *privKey, addr *addrPubKey) {
	dir, err := os.Getwd()
	debug(err)

	// Create QR code for the private key
	pkRGBA := image.NewRGBA(image.Rect(0, 0, 41, 41))
	draw.Draw(pkRGBA, pkRGBA.Bounds(), pk.QR(), image.Point{0, 0}, draw.Src)
	pkImg, err := os.Create("pkCode.jpeg")
	debug(err)
	defer pkImg.Close()
	highQuality := 100
	debug(jpeg.Encode(pkImg, pkRGBA, &jpeg.Options{Quality: highQuality}))

	// Create QR code for the public address
	addrRGBA := image.NewRGBA(image.Rect(0, 0, 33, 33))
	draw.Draw(addrRGBA, addrRGBA.Bounds(), addr.QR(), image.Point{0, 0}, draw.Src)
	addrImg, err := os.Create("addrCode.jpeg")
	debug(err)
	defer addrImg.Close()
	debug(jpeg.Encode(addrImg, addrRGBA, &jpeg.Options{Quality: highQuality}))

	// Create pdf
	paperWallet := pdf.New("P", "mm", "A4", "")
	paperWallet.AddPage()
	fontSize := 8.0
	paperWallet.SetFont("Helvetica", "B", fontSize)
	ht := paperWallet.PointConvert(fontSize)
	tr := paperWallet.UnicodeTranslatorFromDescriptor("") // "" defaults to "cp1252"
	write := func(str string) {
		paperWallet.CellFormat(190, ht, tr(str), "", 1, "C", false, 0, "")
		paperWallet.Ln(ht)
	}
	write(fmt.Sprintf("PrivKey: %s", pk.String()))
	write(fmt.Sprintf("Address: %s", addr.String()))
	walletPath := filepath.Join(dir, "wallet.pdf")

	debug(paperWallet.OutputFileAndClose(walletPath))
	fmt.Println("Successfully generated wallet.pdf")

	// TODO: Clean-up
}
