#!/bin/sh

set -u

CLI_WATCH="docker exec -it btc-watch bitcoin-cli"
CLI_KEYGEN="docker exec -it btc-keygen bitcoin-cli"
CLI_SIGN1="docker exec -it btc-sign1 bitcoin-cli"
CLI_SIGN2="docker exec -it btc-sign2 bitcoin-cli"

$CLI_WATCH createwallet watch
$CLI_KEYGEN createwallet keygen
$CLI_SIGN1 createwallet sign1
$CLI_SIGN2 createwallet sign2
