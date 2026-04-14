### Bitcoind Setup without container

1. install `bitcoind` on macOS directly if needed

- see [bitcoin core installation](https://github.com/bitcoin/bitcoin/blob/master/doc/build-osx.md)

1. run bitcoind `$ bitcoind`
2. create wallets separately (if only one node used)

    ```
    $ bitcoin-cli createwallet watch
    $ bitcoin-cli createwallet keygen
    $ bitcoin-cli createwallet sign1
    $ bitcoin-cli createwallet sign2
    $ bitcoin-cli createwallet sign3
    $ bitcoin-cli createwallet sign4
    $ bitcoin-cli createwallet sign5
    $ bitcoin-cli listwallets
    [
      "",
      "watch",
      "keygen",
      "sign1",
      "sign2",
      "sign3",
      "sign4",
      "sign5"
    ]
    ```
