#!/bin/sh

set -u

CLI_WATCH="docker exec btc-watch bitcoin-cli"
CLI_KEYGEN="docker exec btc-keygen bitcoin-cli"
CLI_SIGN1="docker exec btc-sign1 bitcoin-cli"
CLI_SIGN2="docker exec btc-sign2 bitcoin-cli"

# createwallet "wallet_name" ( disable_private_keys blank "passphrase" avoid_reuse descriptors load_on_startup )
# Create descriptor wallets (Bitcoin Core v0.21.0+) for all wallets

# Watch wallet: watch-only with descriptors (disable_private_keys=true, descriptors=true)
$CLI_WATCH createwallet watch true false "" false true

# Keygen wallet: full wallet with descriptors (disable_private_keys=false, descriptors=true)
$CLI_KEYGEN createwallet keygen false false "" false true

# Sign wallets: full wallets with descriptors (disable_private_keys=false, descriptors=true)
$CLI_SIGN1 createwallet sign1 false false "" false true
$CLI_SIGN2 createwallet sign2 false false "" false true
