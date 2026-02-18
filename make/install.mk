###############################################################################
# Installation Targets
###############################################################################

.PHONY: check-ver
check-ver:
	#echo $(modVer)
	#echo $(currentVer)
	@if [ ${currentVer} -lt ${modVer} ]; then\
		echo go version ${modVer}++ is required but your go version is ${currentVer};\
	fi

.PHONY: install-mac-tools
install-mac-tools:
	brew tap hiromaily/tap \
	brew install yaml-lint \
		jq \
		mkcert \
		ariga/tap/atlas \
		protobuf \
		buf \
		sqlc \
		sqlfluff \
		shellcheck \
		clang-format

.PHONY: install-bun
install-bun:
	curl -fsSL https://bun.sh/install | bash

#------------------------------------------------------------------------------
# xrpl-grpc-server dependencies
#------------------------------------------------------------------------------

# Install all xrpl-grpc-server dependencies
.PHONY: install-xrpl-deps
install-xrpl-deps:
	cd apps/xrpl-grpc-server && bun install

# Check protoc version (required >= 33.0 for Edition 2024)
.PHONY: check-protoc
check-protoc:
	@command -v protoc >/dev/null 2>&1 || { \
		echo "Error: protoc is not installed"; \
		echo "Install via: make install-mac-tools (macOS) or see https://grpc.io/docs/protoc-installation/"; \
		exit 1; \
	}
	@PROTOC_VERSION=$$(protoc --version | sed 's/libprotoc //'); \
	if [ "$$(printf '%s\n' '33' "$$PROTOC_VERSION" | sort -V | head -n 1)" != "33" ]; then \
		echo "Error: protoc version $$PROTOC_VERSION is too old. Required: >= 33.0"; \
		echo "Upgrade via: brew upgrade protobuf (macOS)"; \
		exit 1; \
	fi; \
	echo "protoc version OK: $$PROTOC_VERSION"

.PHONY: install-ssl
install-ssl:
	mkcert -install
	mkcert localhost 127.0.0.1

.PHONY: install-tools-by-gomod
install-tools-by-gomod:
	go get -tool github.com/ethereum/go-ethereum/cmd/abigen@latest
	go get -tool github.com/vektra/mockery/v3@latest
	go get -tool github.com/evilmartians/lefthook@latest
	go get -tool gotest.tools/gotestsum@latest
	go get -tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION)
	go get -tool golang.org/x/vuln/cmd/govulncheck@latest
	go get -tool honnef.co/go/tools/cmd/staticcheck@latest
	go get -tool github.com/icholy/gomajor@latest
	go get -tool mvdan.cc/sh/v3/cmd/gosh@latest
	go get -tool mvdan.cc/sh/v3/cmd/shfmt@latest
	go get -tool github.com/mrtazz/checkmake/cmd/checkmake@latest
	go get -tool github.com/google/yamlfmt/cmd/yamlfmt@latest

# Note: somehow checkmake couldn't manage by go get -tool, so use go install instead.
.PHONY: install-tools
install-tools:
	go install github.com/mrtazz/checkmake/cmd/checkmake@latest

.PHONY: setup-lefthook
setup-lefthook:
	go tool lefthook install

.PHONY: goget
goget:
	go mod download

# For Ethereum between execution and beacon client
.PHONY:jwt
jwt:
	openssl rand -hex 32 | tr -d "\n" > "jwtsecret"
	mv jwtsecret ./docker/nodes/eth/configs/
