###############################################################################
# Test Targets
###############################################################################

.PHONY: gotest
gotest:
	go test -v ./...

.PHONY: gotest-addr
gotest-addr:
	go test -tags=integration -v -run pkg/wallets/api/btc/...
	go test -tags=integration -v -run GetAddressInfo pkg/wallets/api/btc/...
	go test -v pkg/wallets/api/btc/... -run GetAddressInfo

.PHONY: gotest-integration
gotest-integration:
	go test -v -tags=integration ./...
