## cryptowallet [![Build Status](https://travis-ci.org/kargakis/cryptowallet.svg?branch=master)](https://travis-ci.org/kargakis/cryptowallet)
```cryptowallet``` is a command-line paper wallet. Use it off-line to safely generate a pdf containing a private key in WIF and QR code format and a pay-to-pubkey address in base58 and QR code format. Supported coins are Bitcoin, Litecoin, Primecoin, Peercoin, Namecoin, and Dogecoin.

### Install
Install by building from source for now:

	$ go get github.com/kargakis/cryptowallet

Update in a similar fashion:

	$ go get -u -v github.com/kargakis/cryptowallet

### Use
Run the following command and expect a ```wallet.pdf``` to be generated in your current directory:
	
	$ cryptowallet

The generated wallet above will be a Bitcoin paper wallet. If you want to generate a paper wallet for a different supported cryptocurrency, use the ```--coin``` flag. For example, the following commnad will generate a Peercoin paper wallet:

	$ cryptowallet --coin ppc

Supported cryptocurrencies can be seen by running:

	$ cryptowallet --support

For more information run:

	$ cryptowallet --help

### License
See [LICENSE](https://github.com/kargakis/cryptowallet/blob/master/LICENSE).
