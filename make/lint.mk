###############################################################################
# Linter and Code Quality Targets
###############################################################################

# Note: Build tags (e.g., //go:build integration) are supported via .golangci.yml run.build-tags setting
.PHONY: format
format:
	go tool golangci-lint fmt

.PHONY: imports
imports:
	./scripts/imports.sh

.PHONY: lint
lint:
	go tool golangci-lint run

# Note: Build tags (e.g., //go:build integration) are now supported via .golangci.yml run.build-tags setting
.PHONY: lint-fix
lint-fix:
	go tool golangci-lint run --fix

.PHONY: staticcheck
staticcheck:
	go tool staticcheck ./...

.PHONY: check-upgrade
check-upgrade:
	go tool gomajor list

.PHONY: check-vuln
check-vuln:
	go tool govulncheck ./...

.PHONY: shfmt
shfmt:
	shfmt -l -w scripts/*.sh
	shfmt -l -w scripts/*/**.sh

.PHONY: lint-makefile
lint-makefile:
	checkmake Makefile make/*.mk 
