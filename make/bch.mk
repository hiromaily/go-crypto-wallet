###############################################################################
# Bitcoin Cash Node Targets
###############################################################################

###############################################################################
# Docker Compose Targets
###############################################################################

# run all bitcoin cash node servers
.PHONY: up-docker-bch
up-docker-bch:
	docker compose -f compose.bch.yaml up bch-watch bch-keygen bch-sign1 bch-sign2

# run bitcoin cash node servers in detached mode
.PHONY: up-docker-bch-d
up-docker-bch-d:
	docker compose -f compose.bch.yaml up -d bch-watch bch-keygen bch-sign1 bch-sign2

# stop bitcoin cash node servers
.PHONY: down-docker-bch
down-docker-bch:
	docker compose -f compose.bch.yaml down

# stop bitcoin cash node servers and remove volumes
.PHONY: down-docker-bch-v
down-docker-bch-v:
	docker compose -f compose.bch.yaml down -v

###############################################################################
# E2E Testing
###############################################################################
# Run Bitcoin Cash E2E workflow from completely fresh state (recommended)
.PHONY: bch-e2e-test-reset
bch-e2e-test-reset:
	./scripts/operation/bch/e2e-workflow.sh --reset

# Run complete Bitcoin Cash end-to-end workflow (regression test)
.PHONY: bch-e2e-test
bch-e2e-test:
	./scripts/operation/bch/e2e-workflow.sh

# Run Bitcoin Cash E2E workflow with verbose output
.PHONY: bch-e2e-test-verbose
bch-e2e-test-verbose:
	./scripts/operation/bch/e2e-workflow.sh --verbose

# Run Bitcoin Cash E2E workflow in non-interactive mode (for CI/CD)
.PHONY: bch-e2e-test-ci
bch-e2e-test-ci:
	./scripts/operation/bch/e2e-workflow.sh --non-interactive

# Cleanup Bitcoin Cash E2E test environment
.PHONY: bch-e2e-cleanup
bch-e2e-cleanup:
	./scripts/operation/bch/e2e-workflow.sh --cleanup

