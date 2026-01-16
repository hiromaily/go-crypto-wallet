###############################################################################
# Linter and Code Quality Targets
###############################################################################
###############################################################################
# Golang Linting
###############################################################################

# Note: Build tags (e.g., //go:build integration) are supported via .golangci.yml run.build-tags setting
.PHONY: go-fmt
go-fmt:
	go tool golangci-lint fmt

# format imports
.PHONY: go-imports
imports:
	./scripts/imports.sh

# lint by golangci-lint
.PHONY: go-lint-check
go-lint-check:
	go tool golangci-lint run

# lint and fix
.PHONY: go-lint
go-lint:
	go tool golangci-lint run --fix

# clean golangci-lint cache
.PHONY: go-clean-lint-cache
go-clean-lint-cache:
	go tool golangci-lint cache clean

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

###############################################################################
# Shell Script Linting
###############################################################################
# format shell scripts
.PHONY: shfmt
shfmt:
	shfmt -l -w scripts/*.sh
	shfmt -l -w scripts/*/**.sh

.PHONY: shellcheck
shellcheck:
	shellcheck scripts/*.sh
	shellcheck scripts/*/**.sh

###############################################################################
# Makefile Linting
###############################################################################
# lint makefile
.PHONY: mk-lint
mk-lint:
	checkmake Makefile make/*.mk

###############################################################################
# YAML Linting
###############################################################################
.PHONY: yaml-lint
yaml-lint:
	yaml-lint .github/workflows
	yaml-lint .devcontainer
	yaml-lint config

###############################################################################
# SQL Linting (Aliases to db.mk sqlfluff targets)
###############################################################################
.PHONY: sql-fmt
sql-fmt: sqlfluff-format

.PHONY: sql-lint
sql-lint: sqlfluff-lint

.PHONY: sql-fix
sql-fix: sqlfluff-fix

# Lint SQL queries
# This target runs vet to check queries for potential issues
.PHONY: sqlc-lint
sqlc-lint: sqlc-validate


###############################################################################
# HCL Linting (Aliases to db_atlas.mk targets)
###############################################################################
.PHONY: hcl-fmt
hcl-fmt: atlas-fmt

.PHONY: hcl-lint
hcl-lint: atlas-lint

.PHONY: hcl-validate
hcl-validate: atlas-validate

###############################################################################
# Proto Linting (Aliases to codegen.mk targets)
###############################################################################
.PHONY: proto-fmtlint
proto-fmtlint: proto-fmt proto-lint
