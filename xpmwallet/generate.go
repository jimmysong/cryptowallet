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
	addr, err := btcutil.NewAddressPubKey(pk.PubKey().SerializeUncompressed(), primeNet)
	if err != nil {
		fmt.Println("Error generating pay-to-pubkey address:", err)
		safeClose(dump)
		return
	}
	go func(dump chan Data) {
		code, err := qr.Encode(addr.EncodeAddress(), qr.H)
		if err != nil {
			fmt.Println("Cannot create QR code for pay-to-pubkey address:", err)
			safeClose(dump)
			return
		}
		dump <- &addrPubKey{qrCode: code, value: addr}
	}(dump)
}

func safeClose(ch chan Data) {
	select {
	case <-ch:
	default:
		close(ch)
	}
}

func newPaperWallet(wif, addr Data) {
	// Create QR code images
	highQuality := 100
	wifRGBA := image.NewRGBA(image.Rect(0, 0, 41, 41))
	draw.Draw(wifRGBA, wifRGBA.Bounds(), wif.QR(), image.Point{0, 0}, draw.Src)
	wifImg, err := os.Create("wifCode.jpeg")
	if err != nil {
		fmt.Println("Cannot create wifCode.jpeg:", err)
		return
	}
	defer wifImg.Close()
	if err := jpeg.Encode(wifImg, wifRGBA, &jpeg.Options{highQuality}); err != nil {
		fmt.Println("Cannot encode data to wifCode.jpeg:", err)
		return
	}
	addrRGBA := image.NewRGBA(image.Rect(0, 0, 33, 33))
	draw.Draw(addrRGBA, addrRGBA.Bounds(), addr.QR(), image.Point{0, 0}, draw.Src)
	addrImg, err := os.Create("addrCode.jpeg")
	if err != nil {
		fmt.Println("Cannot create addrCode.jpeg:", err)
		return
	}
	defer addrImg.Close()
	if err := jpeg.Encode(addrImg, addrRGBA, &jpeg.Options{highQuality}); err != nil {
		fmt.Println("Cannot encode data to addrCode.jpeg:", err)
		return
	}
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
	write(fmt.Sprintf("PrivKey: %s", wif.String()))
	write(fmt.Sprintf("Address: %s", addr.String()))
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

	// TOOD: Clean-up
}
