###############################################################################
# Build Targets
###############################################################################

# Build on local
# - appVersion is injected via LDFLAGS (defined in vars.mk)
# - authName on sign works as account name

#------------------------------------------------------------------------------
# Go Source Files and Binary Paths
#------------------------------------------------------------------------------
# GO_SRCS: All Go source files in the project
# - Used as dependencies for incremental builds
# - When any .go file is modified, Make detects it and triggers rebuild
# - Excludes vendor/ and .git/ directories
GO_SRCS := $(shell find . -name '*.go' -not -path './vendor/*' -not -path './.git/*')

# Binary output paths
# - These are the actual executable files that `go build` produces
# - Located in $GOPATH/bin/ so they're available in PATH
# - Use 'go env GOPATH' as fallback when GOPATH env var is not set (e.g., CI)
BIN_DIR    := $(or $(GOPATH),$(shell go env GOPATH))/bin
WATCH_BIN  := $(BIN_DIR)/watch
KEYGEN_BIN := $(BIN_DIR)/keygen
SIGN1_BIN  := $(BIN_DIR)/sign1
SIGN2_BIN  := $(BIN_DIR)/sign2
WALLET_BINS := $(WATCH_BIN) $(KEYGEN_BIN) $(SIGN1_BIN) $(SIGN2_BIN)

#------------------------------------------------------------------------------
# Incremental Build Targets (file-based dependencies)
#------------------------------------------------------------------------------
# How it works:
#   Make compares timestamps of the target file (binary) vs dependencies (source files).
#   - If any dependency is NEWER than the target → rebuild
#   - If target is NEWER than all dependencies → skip (already up-to-date)
#
# Example:
#   $(WATCH_BIN): $(GO_SRCS) go.mod go.sum
#   ↑ target      ↑ dependencies
#
#   1st run:  $GOPATH/bin/watch doesn't exist → build
#   2nd run:  watch exists, sources unchanged → "Nothing to be done" (skip)
#   3rd run:  internal/foo.go modified → rebuild
#
# Benefits:
#   - `make btc-e2e-reset P=2` skips build if sources haven't changed
#   - Saves ~10-30 seconds per E2E test run during development
#
# To force rebuild regardless of timestamps:
#   make build-all-force

# Common dependencies for all binaries
# - Include Makefiles so changes to build flags trigger rebuild
BUILD_DEPS := $(GO_SRCS) go.mod go.sum Makefile make/build.mk make/vars.mk

$(WATCH_BIN): $(BUILD_DEPS)
	@echo "Building watch..."
	@go build -ldflags "$(LDFLAGS)" -o $@ ./cmd/watch/

$(KEYGEN_BIN): $(BUILD_DEPS)
	@echo "Building keygen..."
	@go build -ldflags "$(LDFLAGS)" -o $@ ./cmd/keygen/

$(SIGN1_BIN): $(BUILD_DEPS)
	@echo "Building sign1..."
	@go build -ldflags "$(LDFLAGS) -X main.authName=auth1" -o $@ ./cmd/sign/

$(SIGN2_BIN): $(BUILD_DEPS)
	@echo "Building sign2..."
	@go build -ldflags "$(LDFLAGS) -X main.authName=auth2" -o $@ ./cmd/sign/

#------------------------------------------------------------------------------
# PHONY Build Targets
#------------------------------------------------------------------------------

.PHONY: tidy
tidy:
	# go mod verify
	go mod tidy

.PHONY: check-build
check-build: tidy
	go build -ldflags "$(LDFLAGS)" -v -o /dev/null ./cmd/watch/
	go build -ldflags "$(LDFLAGS)" -v -o /dev/null ./cmd/keygen/
	go build -ldflags "$(LDFLAGS) -X main.authName=auth1" -v -o /dev/null ./cmd/sign/

# build-all: Build all wallet binaries (incremental - only rebuilds when sources change)
.PHONY: build-all
build-all: $(WALLET_BINS)

# build-all-force: Force rebuild all binaries (ignores timestamps)
.PHONY: build-all-force
build-all-force: tidy
	go build -ldflags "$(LDFLAGS)" -v -o ${GOPATH}/bin/watch ./cmd/watch/
	go build -ldflags "$(LDFLAGS)" -v -o ${GOPATH}/bin/keygen ./cmd/keygen/
	go build -ldflags "$(LDFLAGS) -X main.authName=auth1" -v -o ${GOPATH}/bin/sign1 ./cmd/sign/
	go build -ldflags "$(LDFLAGS) -X main.authName=auth2" -v -o ${GOPATH}/bin/sign2 ./cmd/sign/

.PHONY: build-watch
build-watch: $(WATCH_BIN)

.PHONY: build-keygen
build-keygen: $(KEYGEN_BIN)

.PHONY: build-sign
build-sign: $(SIGN1_BIN) $(SIGN2_BIN)

# Build from inside docker container
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "$(LDFLAGS)" -o /go/bin/watch ./cmd/watch/main.go
	CGO_ENABLED=0 GOOS=linux go build -ldflags "$(LDFLAGS)" -o /go/bin/keygen ./cmd/keygen/main.go
	CGO_ENABLED=0 GOOS=linux go build -ldflags "$(LDFLAGS) -X main.authName=auth1" -o /go/bin/sign ./cmd/sign/main.go

# Show current version from git tag
.PHONY: show-version
show-version:
	@echo "VERSION: $(VERSION)"

.PHONY: run-watch
run-watch:
	go run ./cmd/watch/ -conf ./config/wallet/btc/watch.yaml
