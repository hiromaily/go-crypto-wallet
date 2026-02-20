
###############################################################################
# Atlas Configuration
###############################################################################
# Atlas configuration file and environments
# DB_DIALECT selects the database dialect (mysql or postgresql)
# Usage: make <target> DB_DIALECT=postgresql
DB_DIALECT ?= mysql
ATLAS_CONFIG := file://tools/atlas/atlas.hcl
ATLAS_SCHEMAS := watch keygen sign

###############################################################################
# Atlas commands validation targets
###############################################################################

# Check that SCHEMA parameter is provided
.PHONY: _check-schema
_check-schema:
ifndef SCHEMA
	@echo "Error: SCHEMA parameter is required"
	@echo "Usage: make <target> SCHEMA=watch|keygen|sign"
	@exit 1
endif

# Check that NAME parameter is provided
.PHONY: _check-name
_check-name:
ifndef NAME
	@echo "Error: NAME parameter is required"
	@echo "Usage: make <target> NAME=<migration_name>"
	@exit 1
endif

###############################################################################
# Atlas Format and Lint Targets
###############################################################################

# Format all HCL schema files
# `atlas schema fmt`
.PHONY: atlas-fmt
atlas-fmt:
	@echo "Formatting Atlas HCL schema files..."
	@atlas schema fmt tools/atlas/schemas/
	@echo "✓ HCL files formatted successfully"

# Check if HCL files are formatted (CI mode)
# Format files and check for git changes
.PHONY: atlas-fmt-check
atlas-fmt-check:
	@echo "Checking Atlas HCL schema file formatting..."
	@atlas schema fmt tools/atlas/schemas/
	@if ! git diff --exit-code tools/atlas/schemas/ > /dev/null 2>&1; then \
		echo "❌ HCL files are not formatted. Run 'make atlas-fmt' to fix."; \
		git diff tools/atlas/schemas/; \
		exit 1; \
	fi
	@echo "✓ All HCL files are properly formatted"

# Lint all HCL schema files with validation
# Usage: make atlas-lint [DB_DIALECT=mysql|postgresql]
# `atlas schema lint`
.PHONY: atlas-lint
atlas-lint:
	@echo "Linting Atlas HCL schema files ($(DB_DIALECT))..."
	@for schema in $(ATLAS_SCHEMAS); do \
		echo "=== Linting $$schema schema ($(DB_DIALECT)) ==="; \
		(cd tools/atlas && atlas schema lint --config file://atlas.hcl --env local_$(DB_DIALECT)_$$schema) || exit 1; \
	done
	@echo "✓ All schemas passed linting"

# Convenience aliases for dialect-specific lint
.PHONY: atlas-lint-mysql
atlas-lint-mysql:
	@$(MAKE) atlas-lint DB_DIALECT=mysql

.PHONY: atlas-lint-postgresql
atlas-lint-postgresql:
	@$(MAKE) atlas-lint DB_DIALECT=postgresql

# Validate Atlas configuration
# Usage: make atlas-validate [DB_DIALECT=mysql|postgresql]
.PHONY: atlas-validate
atlas-validate:
	@echo "Validating Atlas configuration ($(DB_DIALECT))..."
	@for schema in $(ATLAS_SCHEMAS); do \
		(cd tools/atlas && atlas migrate validate --config file://atlas.hcl --env local_$(DB_DIALECT)_$$schema) || exit 1; \
	done
	@echo "✓ Atlas configuration is valid"

###############################################################################
# Development Workflow Targets
#
# Typical workflow after modifying HCL schema files:
#   1. make atlas-dev-reset     # Regenerate migration SQL files from HCL
#   2. make reset-docker-mysql  # Reset Docker volumes and apply migrations
#      OR
#      make atlas-dev-clean     # Reset DB only (keep Docker running)
#
# Target comparison:
#   - atlas-dev-reset:    Operates on FILES (.sql, atlas.sum)
#   - atlas-dev-clean:    Operates on DATABASE (drop and recreate from HCL)
#   - reset-docker-mysql: Operates on DOCKER (recreate volumes and apply migrations)
#
###############################################################################

# Regenerate migrations from HCL schemas (from scratch)
# This target must run after atlas schema files `watch.hcl`, `keygen.hcl`, `sign.hcl` are changed
# WARNING: This deletes all existing migrations and creates new ones
# If no actual content changes are detected, original files are restored
# Usage: make atlas-dev-reset [DB_DIALECT=mysql|postgresql]
# `atlas migrate diff initial_schema`
.PHONY: atlas-dev-reset
atlas-dev-reset:
	@DB_DIALECT=$(DB_DIALECT) \
	./scripts/db/atlas-dev-reset.sh

# Clean databases and reapply from HCL schemas
# Uses admin_* environments which allow schema-level operations (drop/create schema)
# Usage: make atlas-dev-clean [DB_DIALECT=mysql|postgresql]
# `atlas schema clean`
.PHONY: atlas-dev-clean
atlas-dev-clean:
	@echo "⚠️  WARNING: This will drop all databases and recreate them from HCL schemas!"
	@read -p "Are you sure you want to continue? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo "Cleaning and recreating databases ($(DB_DIALECT))..."; \
		for schema in $(ATLAS_SCHEMAS); do \
			echo "=== Cleaning $$schema schema ==="; \
			(cd tools/atlas && atlas schema clean --config file://atlas.hcl --env admin_$(DB_DIALECT)_$$schema --auto-approve) || exit 1; \
		done; \
		echo "Applying HCL schemas..."; \
		$(MAKE) atlas-schema-apply-all DB_DIALECT=$(DB_DIALECT); \
		echo "✓ Databases cleaned and recreated"; \
	else \
		echo "Cancelled."; \
		exit 1; \
	fi

###############################################################################
# Atlas utility targets
###############################################################################

# Show migration status for all schemas
# Usage: make atlas-migrate-status [DB_DIALECT=mysql|postgresql]
# `atlas migrate status`
.PHONY: atlas-migrate-status
atlas-migrate-status:
	@for schema in $(ATLAS_SCHEMAS); do \
		echo "=== $$schema Schema Migration Status ($(DB_DIALECT)) ==="; \
		(cd tools/atlas && atlas migrate status --config file://atlas.hcl --env local_$(DB_DIALECT)_$$schema) || true; \
	done


###############################################################################
# Atlas Schema Management Targets
###############################################################################

# Apply HCL schema directly to database (all schemas)
# Usage: make atlas-schema-apply-all [DB_DIALECT=mysql|postgresql]
# `atlas schema apply`
.PHONY: atlas-schema-apply-all
atlas-schema-apply-all:
	@echo "Applying HCL schemas directly to databases ($(DB_DIALECT))..."
	@for schema in $(ATLAS_SCHEMAS); do \
		$(MAKE) atlas-schema-apply SCHEMA=$$schema DB_DIALECT=$(DB_DIALECT) || exit 1; \
	done
	@echo "✓ All schemas applied successfully"

# Apply HCL schema for a specific schema
# Usage: make atlas-schema-apply SCHEMA=watch [DB_DIALECT=mysql|postgresql]
# `atlas schema apply`
.PHONY: atlas-schema-apply
atlas-schema-apply: _check-schema
	@echo "=== Applying $(SCHEMA) schema ($(DB_DIALECT)) ==="
	@cd tools/atlas && atlas schema apply --config file://atlas.hcl --env local_$(DB_DIALECT)_$(SCHEMA) --auto-approve

###############################################################################
# Atlas Migration Targets
# Note: This target would not be used while developing phase
###############################################################################

# Apply all pending migrations for all schemas
# Usage: make atlas-migrate-apply-all [DB_DIALECT=mysql|postgresql]
# `atlas migrate apply`
.PHONY: atlas-migrate-apply-all
atlas-migrate-apply-all:
	@echo "Applying migrations for all schemas ($(DB_DIALECT))..."
	@for schema in $(ATLAS_SCHEMAS); do \
		echo "=== Applying $$schema migrations ==="; \
		(cd tools/atlas && atlas migrate apply --config file://atlas.hcl --env local_$(DB_DIALECT)_$$schema) || exit 1; \
	done
	@echo "✓ All migrations applied successfully"

# Generate new migration from HCL schema diff
# Usage: make atlas-migrate-diff SCHEMA=watch NAME=add_new_column [DB_DIALECT=mysql|postgresql]
# `atlas migrate diff`
.PHONY: atlas-migrate-diff
atlas-migrate-diff: _check-schema _check-name
	@echo "=== Generating migration for $(SCHEMA) schema ($(DB_DIALECT)) ==="
	@cd tools/atlas && atlas migrate diff $(NAME) --config file://atlas.hcl --env local_$(DB_DIALECT)_$(SCHEMA)
	@echo "✓ Migration generated: tools/atlas/migrations/$(DB_DIALECT)/$(SCHEMA)/"

# Hash migration directory for production readiness
# Usage: make atlas-migrate-hash-all [DB_DIALECT=mysql|postgresql]
.PHONY: atlas-migrate-hash-all
atlas-migrate-hash-all:
	@echo "Hashing migration directories ($(DB_DIALECT))..."
	@for schema in $(ATLAS_SCHEMAS); do \
		(cd tools/atlas && atlas migrate hash --config file://atlas.hcl --env local_$(DB_DIALECT)_$$schema) || exit 1; \
	done
	@echo "✓ Migration directories hashed"

###############################################################################
# Production Workflow Targets
###############################################################################

# Production Mode: Initialize migration history from current database state
# .PHONY: atlas-prod-init
# atlas-prod-init:
# 	@echo "Initializing production migration history..."
# 	@for schema in $(ATLAS_SCHEMAS); do \
# 		echo "=== Hashing $$schema migrations ==="; \
# 		(cd tools/atlas && atlas migrate hash --config file://atlas.hcl --env local_$(DB_DIALECT)_$$schema) || exit 1; \
# 	done
# 	@echo "✓ Production migration history initialized"
# 	@echo "You can now create incremental migrations using: make atlas-migrate-diff SCHEMA=<schema> NAME=<name>"

# Usage: make atlas-migration-flow [DB_DIALECT=mysql|postgresql]
.PHONY: atlas-migration-flow
atlas-migration-flow:

