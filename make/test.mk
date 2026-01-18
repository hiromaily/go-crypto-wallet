###############################################################################
# Test Targets
###############################################################################

# Run all unit tests with gotestsum
.PHONY: gotest
gotest:
	go tool gotestsum --format testname -- -v ./...

# Run integration tests with gotestsum
.PHONY: gotest-integration
gotest-integration:
	go tool gotestsum --format testname -- -v -tags=integration ./...
