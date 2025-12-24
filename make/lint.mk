###############################################################################
# Linter and Code Quality Targets
###############################################################################

.PHONY: imports
imports:
	./scripts/imports.sh

.PHONY: format
format:
	go tool golangci-lint fmt

.PHONY: lint
lint:
	go tool golangci-lint run

# Bug: format doesn't work on files which has tags
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
