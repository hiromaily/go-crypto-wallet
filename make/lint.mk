###############################################################################
# Linter and Code Quality Targets
###############################################################################

# Note: Build tags (e.g., //go:build integration) are supported via .golangci.yml run.build-tags setting
.PHONY: format
format:
	go tool golangci-lint fmt

# format imports
.PHONY: imports
imports:
	./scripts/imports.sh

# lint by golangci-lint
.PHONY: lint
lint:
	go tool golangci-lint run

# lint and fix
.PHONY: lint-fix
lint-fix:
	go tool golangci-lint run --fix

# clean golangci-lint cache
.PHONY: clean-lint-cache
clean-lint-cache:
	go tool golangci-lint cache clean

# staticcheck
.PHONY: staticcheck
staticcheck:
	go tool staticcheck ./...

# check for upgrade
.PHONY: check-upgrade
check-upgrade:
	go tool gomajor list

# check for vulnerabilities
.PHONY: check-vuln
check-vuln:
	go tool govulncheck ./...

# format shell scripts
.PHONY: shfmt
shfmt:
	shfmt -l -w scripts/*.sh
	shfmt -l -w scripts/*/**.sh

# lint makefile
.PHONY: lint-makefile
lint-makefile:
	checkmake Makefile make/*.mk 
