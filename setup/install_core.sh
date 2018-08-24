#!/bin/sh

# http://simoothie-cafe.hatenablog.com/entry/2018/05/06/002539

# clone git
mkdir work
cd ./work

git clone https://github.com/bitcoin/bitcoin.git
cd bitcoin

# setup required pkg
brew install automake libtool pkg-config
./autogen.sh
# configure: error: PKG_PROG_PKG_CONFIG macro not found. Please install pkg-config and re-run autogen.sh

brew install berkeley-db4 boost libevent
brew link berkeley-db4 --force
./configure

make
make check
make install

