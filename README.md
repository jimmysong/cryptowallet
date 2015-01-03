## xpmwallet [![Build Status](https://travis-ci.org/kargakis/xpmwallet.svg?branch=master)](https://travis-ci.org/kargakis/xpmwallet)
```xpmwallet``` is a command-line [Primecoin](http://primecoin.io/) paper wallet. Use it off-line to safely generate a pdf containing a private key in WIF and QR code format and a pay-to-pubkey address in base58 and QR code format.

### Install
Install by building from source for now:

	$ go get github.com/kargakis/xpmwallet/...

### Use
Run the following command and expect a ```wallet.pdf``` to be generated in your current directory:
	
	$ xpmwallet

For more information run:

	$ xpmwallet --help

### License
See [LICENSE](https://github.com/kargakis/prime/blob/master/LICENSE).
