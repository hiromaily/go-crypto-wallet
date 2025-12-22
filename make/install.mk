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
	brew install protobuf prototool
	brew tap bufbuild/buf                                                                                                                                                                      (git)-[master]
	brew install buf

.PHONY: install-ssl
install-ssl:
	mkcert -install
	mkcert localhost 127.0.0.1

.PHONY: install-tools
install-tools:
	go install github.com/ethereum/go-ethereum/cmd/abigen@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION)
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/icholy/gomajor@latest
	go install mvdan.cc/sh/v3/cmd/gosh@latest
	go install mvdan.cc/sh/v3/cmd/shfmt@latest

.PHONY: install-tools-by-gomod
install-tools-by-gomod:
	go get -tool github.com/ethereum/go-ethereum/cmd/abigen@latest
	go get -tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION)
	go get -tool golang.org/x/vuln/cmd/govulncheck@latest
	go get -tool honnef.co/go/tools/cmd/staticcheck@latest
	go get -tool github.com/icholy/gomajor@latest
	go get -tool mvdan.cc/sh/v3/cmd/gosh@latest

.PHONY: install-proto-plugin
install-proto-plugin:
	# refer to https://developers.google.com/protocol-buffers/docs/reference/go-generated
	# provides a protoc-gen-go binary which protoc uses when invoked with the --go_out command-line flag.
	# The --go_out flag
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	# The go-grpc_out flag
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	# gogo protobuf
	go install github.com/gogo/protobuf/protoc-gen-gogo@latest
	#go get github.com/gogo/protobuf/proto
	#go get github.com/gogo/protobuf/jsonpb
	#go get github.com/gogo/protobuf/gogoproto

.PHONY: goget
goget:
	go mod download

# For Ethereum between execution and beacon client
.PHONY:jwt
jwt:
	openssl rand -hex 32 | tr -d "\n" > "jwtsecret"
	mv jwtsecret ./docker/nodes/eth/configs/
