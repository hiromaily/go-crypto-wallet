###############################################################################
# Linter and Code Quality Targets
###############################################################################

###############################################################################
# Lint Related Targets
###############################################################################
# go-related targets
include make/lint-go.mk


###############################################################################
# Shell Script Linting
###############################################################################
# format shell scripts
.PHONY: shfmt
shfmt:
	find scripts -name '*.sh' -exec shfmt -l -w {} +

# Excluded shellcheck codes:
#   SC2086: Double quote to prevent globbing (safe in controlled context)
#   SC2034: Unused variables (often used via source)
#   SC1091: Not following sourced files (shellcheck limitation)
#   SC2143: Use grep -q (style preference)
#   SC2181: Check exit code directly (style preference)
#   SC3043: POSIX sh 'local' undefined (scripts use bash)
#   SC3018: POSIX sh '++' undefined (scripts use bash)
#   SC2016: Expressions in single quotes (intentional)
#   SC2001: Use ${var//} instead of sed (style suggestion)
#   SC2295: Nested quoting style (style suggestion)
.PHONY: shellcheck
shellcheck:
	find scripts -name '*.sh' -exec shellcheck --exclude=SC2086,SC2034,SC1091,SC2143,SC2181,SC3043,SC3018,SC2016,SC2001,SC2295 {} +

###############################################################################
# Makefile Linting
###############################################################################
# lint makefile
.PHONY: mk-lint
mk-lint:
	@checkmake Makefile make/*.mk | awk '\
	/^  minphony/ { skip=1; next } \
	skip && /^  [a-z]/ { skip=0 } \
	!skip \
	'
	@#checkmake Makefile make/*.mk

###############################################################################
# Markdown Linting
###############################################################################
# FIXME: a lot of lint errors ...
.PHONY: md-lint
md-lint:
	markdownlint-cli2 --fix "**/*.md"


###############################################################################
# YAML Linting
###############################################################################
yaml-fmt:
	go tool yamlfmt .
	go tool yamlfmt .github/workflows
	go tool yamlfmt .github/actions
	go tool yamlfmt .devcontainer
	go tool yamlfmt config

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

# Lint sqlc configuration and queries
# This target runs sqlc-validate, which combines compile and vet checks.
.PHONY: sqlc-lint
sqlc-lint: sqlc-validate


###############################################################################
# HCL Linting (Aliases to db_atlas.mk targets)
###############################################################################
.PHONY: hcl-fmt
hcl-fmt: atlas-fmt

.PHONY: hcl-fmt-check
hcl-fmt-check: atlas-fmt-check

.PHONY: hcl-lint
hcl-lint: atlas-lint

.PHONY: hcl-validate
hcl-validate: atlas-validate

###############################################################################
# Proto Linting (Aliases to codegen.mk targets)
###############################################################################
.PHONY: proto-fmtlint
proto-fmtlint: proto-fmt-check proto-lint
