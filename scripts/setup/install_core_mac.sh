#!/bin/sh

# https://github.com/bitcoin/bitcoin/blob/master/doc/build-osx.md

# Clone git
mkdir work
cd ./work

git clone https://github.com/bitcoin/bitcoin.git
cd bitcoin

# Install required pkg
brew install automake berkeley-db4 libtool boost miniupnpc pkg-config python qt libevent qrencode

# Berkeley DB
./contrib/install_db4.sh .
$()$(
	When compiling bitcoind, run
)./configure$(
	in the following way:

	export BDB_PREFIX='/Users/hy/work/btc/bitcoin/db4'
	./configure BDB_LIBS="-L${BDB_PREFIX}/lib -ldb_cxx-4.8" BDB_CFLAGS="-I${BDB_PREFIX}/include" ...
)$()

# Build Bitcoin Core
./autogen.sh
make
./configure
make check

# Link
ln -s ${HOME}/work/btc/bitcoin/src/bitcoind /usr/local/bin/bitcoind
ln -s ${HOME}/work/btc/bitcoin/src/bitcoin-cli /usr/local/bin/bitcoin-cli
ln -s ${HOME}/work/btc/bitcoin/src/bitcoin-tx /usr/local/bin/bitcoin-tx

# daemon
bitcoind -daemon

# cli
bitcoin-cli
