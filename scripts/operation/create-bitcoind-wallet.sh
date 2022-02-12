#!/bin/sh

set -u

CLI_WATCH="docker exec -it btc-watch bitcoin-cli"
CLI_KEYGEN="docker exec -it btc-keygen bitcoin-cli"
CLI_SIGN="docker exec -it btc-sign bitcoin-cli"

$CLI_WATCH createwallet watch
$CLI_KEYGEN createwallet keygen
$CLI_SIGN createwallet sign1
$CLI_SIGN createwallet sign2
$CLI_SIGN createwallet sign3
$CLI_SIGN createwallet sign4
$CLI_SIGN createwallet sign5
