#!/bin/sh

# http://simoothie-cafe.hatenablog.com/entry/2018/05/06/002539

# clone git
mkdir work
cd ./work

git clone https://github.com/bitcoin/bitcoin.git
cd bitcoin

# setup required pkg
brew install automake libtool
./autogen.sh

brew install berkeley-db4 boost libevent
brew link berkeley-db4 --force
./configure
make
make check
make install

