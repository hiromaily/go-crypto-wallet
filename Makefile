goget:
	go get -u -d -v ./...

bld:
	go build -o wallet ./cmd/wallet/main.go

run: bld
	./wallet

.PHONY: clean
clean:
	rm -rf detect