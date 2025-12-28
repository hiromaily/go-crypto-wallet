#!/bin/sh

set -u

CLI_WATCH="docker exec btc-watch bitcoin-cli"
CLI_KEYGEN="docker exec btc-keygen bitcoin-cli"
CLI_SIGN1="docker exec btc-sign1 bitcoin-cli"
CLI_SIGN2="docker exec btc-sign2 bitcoin-cli"

$CLI_WATCH loadwallet watch
$CLI_KEYGEN loadwallet keygen
$CLI_SIGN1 loadwallet sign1
$CLI_SIGN2 loadwallet sign2
