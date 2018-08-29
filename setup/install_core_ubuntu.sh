#!/bin/sh

sudo apt-add-repository ppa:bitcoin/bitcoin
sudo apt-get update
sudo apt-get install bitcoind

mkdir ~/.bitcoin
touch ~/.bitcoin/bitcoin.conf
