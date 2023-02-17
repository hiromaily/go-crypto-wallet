mnemonic="math razor capable expose worth grape metal sunset metal sudden usage scheme"

.PHONY: build
build:
	go build -o ${GOPATH}/bin/get-eth-key ./get-eth-key/*.go

# {dir}/{address} file must be existing
run:
	get-eth-key --addr 0xXXXXXX --dir ./keystore

run2:
	get-eth-key --mnemonic $(mnemonic) --hdpath "m/44'/60'/0'/0/0"
