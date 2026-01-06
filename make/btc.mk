###############################################################################
# Bitcoin Core Targets
###############################################################################

.PHONY: bitcoin-run
bitcoin-run:
	bitcoind -daemon

.PHONY: bitcoin-stop
bitcoin-stop:
	bitcoin-cli stop

# MacOS only
.PHONY: cd-btc-dir
cd-btc-dir:
	cd ~/Library/Application\ Support/Bitcoin

###############################################################################
# Docker Compose Targets
###############################################################################

# run bitcoin core server
.PHONY: up-docker-btc
up-docker-btc:
	docker compose -f compose.btc.yaml up btc-watch btc-keygen btc-sign

# run bitcoin cash core server
.PHONY: up-docker-bch
up-docker-bch:
	docker compose -f compose.bch.yaml up bch-watch

###############################################################################
# auto key generator
###############################################################################
.PHONY: generate-btc-key-local
generate-btc-key-local:
	./scripts/operation/generate-btc-key.sh btc false 5

.PHONY: generate-bch-key-local
generate-bch-key-local:
	./scripts/operation/generate-btc-key.sh bch false 5

###############################################################################
# E2E Testing
###############################################################################
# Run complete Bitcoin end-to-end workflow (regression test)
.PHONY: btc-e2e-test
btc-e2e-test:
	./scripts/operation/btc/e2e-workflow.sh

# Run Bitcoin E2E workflow with verbose output
.PHONY: btc-e2e-test-verbose
btc-e2e-test-verbose:
	./scripts/operation/btc/e2e-workflow.sh --verbose

# Run Bitcoin E2E workflow in non-interactive mode (for CI/CD)
.PHONY: btc-e2e-test-ci
btc-e2e-test-ci:
	./scripts/operation/btc/e2e-workflow.sh --non-interactive

# Cleanup Bitcoin E2E test environment
.PHONY: btc-e2e-cleanup
btc-e2e-cleanup:
	./scripts/operation/btc/e2e-workflow.sh --cleanup
