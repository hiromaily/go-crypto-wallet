###############################################################################
# Golang Linting
###############################################################################
# 2.5s as of Feb, 2026
.PHONY: go-fmt
go-fmt:
	go tool golangci-lint fmt

# lint
# 65.38s as of Feb, 2026
.PHONY: go-lint-check
go-lint-check:
	go tool golangci-lint run

# lint and fix
.PHONY: go-lint
go-lint:
	go tool golangci-lint run --fix

# 6.12s as of Feb, 2026
.PHONY: go-lint-fast-check
go-lint-fast-check:
	go tool golangci-lint run --fast-only

.PHONY: go-lint-fast
go-lint-fast:
	go tool golangci-lint run --fast-only --fix

# verify golangci-lint configuration
# Note: run after modifying .golangci.yml
.PHONY: go-lint-verify-config
go-lint-verify-config:
	go tool golangci-lint config verify

# clean golangci-lint cache
.PHONY: go-clean-lint-cache
go-clean-lint-cache:
	go tool golangci-lint cache clean


###############################################################################
# Others
###############################################################################
# format imports
.PHONY: go-imports
imports:
	./scripts/lint/imports.sh

# staticcheck
.PHONY: go-staticcheck
go-staticcheck:
	go tool staticcheck ./...

# check for upgrade
.PHONY: go-check-upgrade
go-check-upgrade:
	go tool gomajor list

# check for vulnerabilities
.PHONY: go-check-vuln
go-check-vuln:
	go tool govulncheck ./...

# Use for refactoring plan
.PHONY: go-lint-deps
go-lint-deps:
	./scripts/lint/deps.sh

.PHONY: go-lint-visibility
go-lint-visibility:
	./scripts/lint/visibility.sh
