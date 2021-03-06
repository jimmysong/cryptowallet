// Copyright (C) 2014-15 Michail Kargakis
// This source code is subject to the terms
// of the MIT License

package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"github.com/btcsuite/btcec"
	"github.com/btcsuite/btcutil"

	pdf "code.google.com/p/gofpdf"
	"code.google.com/p/rsc/qr"
)

// PrivKey is the private key of a cryptocoin public address
// in WIF and QR code format.
type PrivKey struct {
	qrCode *qr.Code
	value  *btcutil.WIF
}

// QR returns the QR code of a private key.
func (pk *PrivKey) QR() image.Image { return pk.qrCode.Image() }
func (pk *PrivKey) String() string  { return fmt.Sprint(pk.value.String()) }

// AddrPubKey is a cryptocoin public address of a private key
// in pay-to-pubkey and QR code format.
type AddrPubKey struct {
	qrCode *qr.Code
	value  *btcutil.AddressPubKey
}

// QR returns the QR code of a public address.
func (a *AddrPubKey) QR() image.Image { return a.qrCode.Image() }
func (a *AddrPubKey) String() string  { return a.value.EncodeAddress() }

const highQuality = 100

// NewPrivKey returns a new private key in WIF and QR code format.
func NewPrivKey() *PrivKey {
	// Generate new private key
	pk, err := btcec.NewPrivateKey(btcec.S256())
	debug(err, "Cannot generate new private key")
	wif, err := btcutil.NewWIF(pk, netParams, false)
	debug(err, "Cannot encode private key to WIF")
	pkCode, err := qr.Encode(wif.String(), qr.H)
	debug(err, "Cannot encode WIF to QR code")
	return &PrivKey{qrCode: pkCode, value: wif}
}

// NewAddress returns a new public address derived from the
// passed private key.
func NewAddress(pk *btcutil.WIF) *AddrPubKey {
	// Extract public from private key, serialize it, and create a new pay-to-pubkey address
	addr, err := btcutil.NewAddressPubKey(pk.PrivKey.PubKey().SerializeUncompressed(), netParams)
	debug(err, "Cannot extract public address from private key")
	addrCode, err := qr.Encode(addr.EncodeAddress(), qr.H)
	debug(err, "Cannot encode public address to QR code")
	return &AddrPubKey{qrCode: addrCode, value: addr}
}

// NewPaperWallet accepts a private key and generates a pdf
// paper wallet.
func NewPaperWallet(pk *PrivKey) {
	dir, err := os.Getwd()
	debug(err, "Cannot get current working directory")

	// A wallet.pdf already exists in the current directory.
	// Do not overwrite it so abort new wallet generation.
	if _, err := os.Open("wallet.pdf"); !os.IsNotExist(err) {
		fmt.Println("wallet.pdf already exists!")
		os.Exit(1)
	}

	addr := NewAddress(pk.value)

	// Create QR code for the private key
	pkRGBA := image.NewRGBA(image.Rect(0, 0, 41, 41))
	draw.Draw(pkRGBA, pkRGBA.Bounds(), pk.QR(), image.Point{0, 0}, draw.Src)
	pkImg, err := os.Create("pkCode.jpeg")
	debug(err, "Cannot create pkCode.jpeg")
	debug(jpeg.Encode(pkImg, pkRGBA, &jpeg.Options{Quality: highQuality}), "Cannot encode private key QR into pkCode.jpeg")

	// Create QR code for the public address
	addrRGBA := image.NewRGBA(image.Rect(0, 0, 33, 33))
	draw.Draw(addrRGBA, addrRGBA.Bounds(), addr.QR(), image.Point{0, 0}, draw.Src)
	addrImg, err := os.Create("addrCode.jpeg")
	debug(err, "Cannot create addrCode.jpeg")
	debug(jpeg.Encode(addrImg, addrRGBA, &jpeg.Options{Quality: highQuality}), "Cannot encode public address QR into addrCode.jpeg")

	// Create pdf
	paperWallet := pdf.New("P", "mm", "A4", "")
	paperWallet.AddPage()
	paperWallet.SetFont("Helvetica", "B", 10.0)
	tr := paperWallet.UnicodeTranslatorFromDescriptor("") // "" defaults to "cp1252"
	paperWallet.CellFormat(190, 20, tr(fmt.Sprintf("PrivKey: %s", pk.String())), "", 1, "C", false, 0, "")
	paperWallet.Image(pkImg.Name(), 80, 25, 50, 50, false, "JPEG", 0, "")
	logoPath := coinLogo(dir)
	paperWallet.Image(logoPath, 90, 90, 100, 100, false, "", 0, "")
	paperWallet.CellFormat(190, 230, tr(fmt.Sprintf("Address: %s", addr.String())), "", 1, "C", false, 0, "")
	paperWallet.Image(addrImg.Name(), 80, 150, 50, 50, false, "JPEG", 0, "")
	walletPath := filepath.Join(dir, "wallet.pdf")
	debug(paperWallet.OutputFileAndClose(walletPath), "Cannot generate wallet.pdf")
	fmt.Println("Successfully generated wallet.pdf")

	// Clean-up
	pkImg.Close()
	addrImg.Close()
	debug(os.Remove(filepath.Join(dir, pkImg.Name())), "Cannot remove pkCode.jpeg")
	debug(os.Remove(filepath.Join(dir, addrImg.Name())), "Cannot remove addrCode.jpeg")
	debug(os.Remove(logoPath), "Cannot remove logo image")
}

func coinLogo(dir string) string {
	logoData, err := Logo("logo.png")
	debug(err, "Cannot find embedded logo data")
	buf := bytes.NewBuffer(logoData)
	logo, err := png.Decode(buf)
	debug(err, "Cannot decode embedded logo data into png")
	logoRGBA := image.NewRGBA(image.Rect(0, 0, 900, 900))
	draw.Draw(logoRGBA, logoRGBA.Bounds(), logo, image.Point{0, 0}, draw.Src)
	logoImg, err := os.Create("logo.png")
	debug(err, "Cannot create logo.png")
	debug(png.Encode(logoImg, logoRGBA), "Cannot encode logo data into logo.png")
	logoImg.Close()
	return filepath.Join(dir, "logo.png")
}
