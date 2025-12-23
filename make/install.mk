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
	brew install jq mkcert go-task/tap/go-task

.PHONY: install-protobuf
install-protobuf:
	brew install protobuf
	brew tap bufbuild/buf                                                                                                                                                                      (git)-[master]
	brew install buf

.PHONY: install-ssl
install-ssl:
	mkcert -install
	mkcert localhost 127.0.0.1

.PHONY: install-tools-by-gomod
install-tools-by-gomod:
	go get -tool github.com/ethereum/go-ethereum/cmd/abigen@latest
	go get -tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION)
	go get -tool golang.org/x/vuln/cmd/govulncheck@latest
	go get -tool honnef.co/go/tools/cmd/staticcheck@latest
	go get -tool github.com/icholy/gomajor@latest
	go get -tool mvdan.cc/sh/v3/cmd/gosh@latest
	go get -tool mvdan.cc/sh/v3/cmd/shfmt@latest

.PHONY: goget
goget:
	go mod download

# For Ethereum between execution and beacon client
.PHONY:jwt
jwt:
	openssl rand -hex 32 | tr -d "\n" > "jwtsecret"
	mv jwtsecret ./docker/nodes/eth/configs/
