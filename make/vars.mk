###############################################################################
# Variable Definitions
###############################################################################

# Go version check
modVer=$(shell cat go.mod | head -n 3 | tail -n 1 | awk '{print $2}' | cut -d'.' -f2)
currentVer=$(shell go version | awk '{print $3}' | sed -e "s/go//" | cut -d'.' -f2)

# Tool versions
GOLANGCI_VERSION=v2.7.2
# Note: PROTOC_BIN is deprecated. Use 'buf' directly for Protocol Buffer operations.

# ETH Variables
GETH_HTTP_PORT=8546
BEACON_HTTP_PORT=9596
GETH_VERSION=v1.10.26
LODESTAR_VERSION=v1.4.3
#ETH_CHAIN_ID=11155111 # used in docker-compose.eth.yml.
TARGET_NETWORK=sepolia
# https://eth-clients.github.io/checkpoint-sync-endpoints/
CHECKPOINT_SYNC_URL=https://beaconstate-${TARGET_NETWORK}.chainsafe.io

# Timestamp calculation (OS-dependent)
OS=$(shell uname -s)
timestamp=""
ifeq ($(OS), Darwin)
	timestamp=$(shell date -v+10S '+%s')
else
	timestamp=$(shell date -d'+10second' +%s)
endif
