**THIS PROJECT IS STILL UNDER EARLY DEVELOPMENT SO THE DOCUMENTATION WILL BE VERY BASIC FOR SOMETIME**

## xpmwallet
```xpmwallet``` is a command-line [Primecoin](http://primecoin.io/) paper wallet. Use it off-line to safely generate a pdf containing a private key in WIF and QR code format and a pay-to-pubkey address in base58 and QR code format.

### Install
There are going to be two ways of installing ```xpmwallet```, by building from source or by downloading a binary file. If you wish to follow the former way, you will have to have pre-installed the [Go](https://golang.org/) language, in order to use ```go get```, the only dependency required when building from source.

### Use
It will be something along the lines of the following command:
	
	$ xpmwallet

If everything goes well, a ```wallet.pdf``` containing your private key and primecoin address will be created in the current directory. Simple enough!

### License
See [LICENSE](https://github.com/kargakis/prime/blob/master/LICENSE).
