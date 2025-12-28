#!/bin/sh

set -u

CLI_WATCH="docker exec -it btc-watch bitcoin-cli"
CLI_KEYGEN="docker exec -it btc-keygen bitcoin-cli"
CLI_SIGN="docker exec -it btc-sign bitcoin-cli"

$CLI_WATCH loadwallet watch
$CLI_KEYGEN loadwallet keygen
$CLI_SIGN loadwallet sign1
$CLI_SIGN loadwallet sign2
$CLI_SIGN loadwallet sign3
$CLI_SIGN loadwallet sign4
$CLI_SIGN loadwallet sign5
