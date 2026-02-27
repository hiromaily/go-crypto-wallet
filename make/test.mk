###############################################################################
# Test Targets
###############################################################################

# Run all unit tests with gotestsum
.PHONY: go-test
go-test:
	go tool gotestsum --format testname -- -v ./...

# Run integration tests with gotestsum
.PHONY: go-test-integration
go-test-integration:
	go tool gotestsum --format testname -- -v -tags=integration ./...
