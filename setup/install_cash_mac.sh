#!/bin/sh

#https://github.com/Bitcoin-ABC/bitcoin-abc/blob/master/doc/build-osx.md

brew install automake berkeley-db libtool boost --c++11 miniupnpc openssl pkg-config protobuf --c++11 qt5 libevent

# clone git
mkdir work
cd ./work

git clone https://github.com/Bitcoin-ABC/bitcoin-abc.git
cd bitcoin-abc

./autogen.sh

brew link berkeley-db4 --force
#Warning: Already linked: /usr/local/Cellar/berkeley-db@4/4.8.30
./configure
#configure: error: libdb_cxx headers missing, Bitcoin ABC requires this library for wallet functionality (--disable-wallet to disable wallet functionality)

make

make check

make deploy

