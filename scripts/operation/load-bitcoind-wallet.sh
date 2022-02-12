#!/bin/sh

set -eu

bitcoin-cli-watch loadwallet watch
bitcoin-cli-keygen loadwallet keygen
bitcoin-cli-sign loadwallet sign1
bitcoin-cli-sign loadwallet sign2
bitcoin-cli-sign loadwallet sign3
bitcoin-cli-sign loadwallet sign4
bitcoin-cli-sign loadwallet sign5
