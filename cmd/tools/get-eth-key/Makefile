mnemonic="math razor capable expose worth grape metal sunset metal sudden usage scheme"

.PHONY: build
build:
	go build -o ${GOPATH}/bin/get-eth-key ./*.go

# {dir}/{address} file must be existing
.PHONY: exec1
exec1:
	get-eth-key --addr 0xXXXXXX --dir ./keystore

.PHONY: exec2
exec2:
	get-eth-key --mnemonic $(mnemonic) --hdpath "m/44'/60'/0'/0/0"
