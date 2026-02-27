###############################################################################
# Variable Definitions
###############################################################################

# App version from git tag (fallback to "dev" if no tags)
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")
# For development builds, include commit info: git describe --tags --always --dirty
LDFLAGS := -X main.appVersion=$(VERSION)

# Go version check
modVer=$(shell cat go.mod | head -n 3 | tail -n 1 | awk '{print $2}' | cut -d'.' -f2)
currentVer=$(shell go version | awk '{print $3}' | sed -e "s/go//" | cut -d'.' -f2)

# Tool versions
GOLANGCI_VERSION=v2.10.1
# Note: PROTOC_BIN is deprecated. Use 'buf' directly for Protocol Buffer operations.

# Timestamp calculation (OS-dependent)
OS=$(shell uname -s)
timestamp=""
ifeq ($(OS), Darwin)
	timestamp=$(shell date -v+10S '+%s')
else
	timestamp=$(shell date -d'+10second' +%s)
endif
