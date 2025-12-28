#!/bin/sh

set -u

CLI_WATCH="docker exec btc-watch bitcoin-cli"
CLI_KEYGEN="docker exec btc-keygen bitcoin-cli"
CLI_SIGN1="docker exec btc-sign1 bitcoin-cli"
CLI_SIGN2="docker exec btc-sign2 bitcoin-cli"

$CLI_WATCH createwallet watch
$CLI_KEYGEN createwallet keygen
$CLI_SIGN1 createwallet sign1
$CLI_SIGN2 createwallet sign2
