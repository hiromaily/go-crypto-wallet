#!/bin/sh

set -eu

COIN="${1:?eth}"

###############################################################################
# keygen wallet
###############################################################################
# create seed
keygen create seed

# create hdkey for client, deposit, payment account
keygen -coin ${COIN} create hdkey -account client -keynum 10
keygen -coin ${COIN} create hdkey -account deposit -keynum 10
keygen -coin ${COIN} create hdkey -account payment -keynum 10
keygen -coin ${COIN} create hdkey -account stored -keynum 10

# import generated private key into keygen wallet
keygen -coin ${COIN} import privkey -account client
keygen -coin ${COIN} import privkey -account deposit
keygen -coin ${COIN} import privkey -account payment
keygen -coin ${COIN} import privkey -account stored

# export addresses to watch wallet
