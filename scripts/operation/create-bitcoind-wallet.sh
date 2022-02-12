#!/bin/sh

set -eu

bitcoin-cli-watch createwallet watch
bitcoin-cli-keygen createwallet keygen
bitcoin-cli-sign createwallet sign1
bitcoin-cli-sign createwallet sign2
bitcoin-cli-sign createwallet sign3
bitcoin-cli-sign createwallet sign4
bitcoin-cli-sign createwallet sign5
