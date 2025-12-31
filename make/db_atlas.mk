
###############################################################################
# Atlas Configuration
###############################################################################
# Atlas configuration file and environments
ATLAS_CONFIG := file://tools/atlas/atlas.hcl
ATLAS_ENV_WATCH := local_watch
ATLAS_ENV_KEYGEN := local_keygen
ATLAS_ENV_SIGN := local_sign
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

# Lint all HCL schema files with validation
# `atlas schema lint`
.PHONY: atlas-lint
atlas-lint:
	@echo "Linting Atlas HCL schema files..."
	@for schema in $(ATLAS_SCHEMAS); do \
		echo "=== Linting $$schema schema ==="; \
		(cd tools/atlas && atlas schema lint --config file://atlas.hcl --env local_$$schema) || exit 1; \
	done
	@echo "✓ All schemas passed linting"

# Validate Atlas configuration
.PHONY: atlas-validate
atlas-validate:
	@echo "Validating Atlas configuration..."
	@for schema in $(ATLAS_SCHEMAS); do \
		(cd tools/atlas && atlas migrate validate --config file://atlas.hcl --env local_$$schema) || exit 1; \
	done
	@echo "✓ Atlas configuration is valid"

###############################################################################
# Development Workflow Targets
###############################################################################

# Regenerate migrations from HCL schemas (from scratch)
# This target must run after atlas schema files `watch.hcl`, `keygen.hcl`, `sign.hcl` are changed
# WARNING: This deletes all existing migrations and creates new ones
# `atlas migrate diff initial_schema`
.PHONY: atlas-dev-reset
atlas-dev-reset:
	@echo "⚠️  WARNING: This will delete all existing migrations!"
	@echo "⚠️  This is intended for development only."
	@read -p "Are you sure you want to continue? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo "Regenerating migrations from HCL schemas..."; \
		rm -f tools/atlas/migrations/watch/*.sql tools/atlas/migrations/watch/atlas.sum; \
		rm -f tools/atlas/migrations/keygen/*.sql tools/atlas/migrations/keygen/atlas.sum; \
		rm -f tools/atlas/migrations/sign/*.sql tools/atlas/migrations/sign/atlas.sum; \
		echo "=== Generating watch migration ==="; \
		(cd tools/atlas && atlas migrate diff initial_schema --config file://atlas.hcl --env $(ATLAS_ENV_WATCH)) || exit 1; \
		echo "=== Generating keygen migration ==="; \
		(cd tools/atlas && atlas migrate diff initial_schema --config file://atlas.hcl --env $(ATLAS_ENV_KEYGEN)) || exit 1; \
		echo "=== Generating sign migration ==="; \
		(cd tools/atlas && atlas migrate diff initial_schema --config file://atlas.hcl --env $(ATLAS_ENV_SIGN)) || exit 1; \
		echo "✓ Migrations regenerated. Run 'make reset-docker-db' to apply."; \
	else \
		echo "Cancelled."; \
		exit 1; \
	fi

# Clean databases and reapply from HCL schemas
# Uses admin_* environments which allow schema-level operations (drop/create schema)
# `atlas schema clean`
.PHONY: atlas-dev-clean
atlas-dev-clean:
	@echo "⚠️  WARNING: This will drop all databases and recreate them from HCL schemas!"
	@read -p "Are you sure you want to continue? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo "Cleaning and recreating databases..."; \
		for schema in $(ATLAS_SCHEMAS); do \
			echo "=== Cleaning $$schema schema ==="; \
			(cd tools/atlas && atlas schema clean --config file://atlas.hcl --env admin_$$schema --auto-approve) || exit 1; \
		done; \
		echo "Applying HCL schemas..."; \
		$(MAKE) atlas-schema-apply-all; \
		echo "✓ Databases cleaned and recreated"; \
	else \
		echo "Cancelled."; \
		exit 1; \
	fi

###############################################################################
# Atlas utility targets
###############################################################################

# Show migration status for all schemas
# `atlas migrate status`
.PHONY: atlas-migrate-status
atlas-migrate-status:
	@for schema in $(ATLAS_SCHEMAS); do \
		echo "=== $$schema Schema Migration Status ==="; \
		(cd tools/atlas && atlas migrate status --config file://atlas.hcl --env local_$$schema) || true; \
	done


###############################################################################
# Atlas Schema Management Targets
###############################################################################

# Apply HCL schema directly to database (all schemas)
# `atlas schema apply`
.PHONY: atlas-schema-apply-all
atlas-schema-apply-all:
	@echo "Applying HCL schemas directly to databases..."
	@for schema in $(ATLAS_SCHEMAS); do \
		$(MAKE) atlas-schema-apply SCHEMA=$$schema || exit 1; \
	done
	@echo "✓ All schemas applied successfully"

# Apply HCL schema for a specific schema
# Usage: make atlas-schema-apply SCHEMA=watch
# `atlas schema apply`
.PHONY: atlas-schema-apply
atlas-schema-apply: _check-schema
	@echo "=== Applying $(SCHEMA) schema ==="
	@cd tools/atlas && atlas schema apply --config file://atlas.hcl --env local_$(SCHEMA) --auto-approve

###############################################################################
# Atlas Migration Targets
# Note: This target would not be used while developing phase
###############################################################################

# Apply all pending migrations for all schemas
# `atlas migrate apply`
.PHONY: atlas-migrate-apply-all
atlas-migrate-apply-all:
	@echo "Applying migrations for all schemas..."
	@for schema in $(ATLAS_SCHEMAS); do \
		echo "=== Applying $$schema migrations ==="; \
		(cd tools/atlas && atlas migrate apply --config file://atlas.hcl --env local_$$schema) || exit 1; \
	done
	@echo "✓ All migrations applied successfully"

# Generate new migration from HCL schema diff
# Usage: make atlas-migrate-diff SCHEMA=watch NAME=add_new_column
# `atlas migrate diff`
.PHONY: atlas-migrate-diff
atlas-migrate-diff: check-schema check-name
	@echo "=== Generating migration for $(SCHEMA) schema ==="
	@cd tools/atlas && atlas migrate diff $(NAME) --config file://atlas.hcl --env local_$(SCHEMA)
	@echo "✓ Migration generated: tools/atlas/migrations/$(SCHEMA)/"

# Hash migration directory for production readiness
.PHONY: atlas-migrate-hash-all
atlas-migrate-hash-all:
	@echo "Hashing migration directories..."
	@for schema in $(ATLAS_SCHEMAS); do \
		(cd tools/atlas && atlas migrate hash --config file://atlas.hcl --env local_$$schema) || exit 1; \
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
# 		(cd tools/atlas && atlas migrate hash --config file://atlas.hcl --env local_$$schema) || exit 1; \
# 	done
# 	@echo "✓ Production migration history initialized"
# 	@echo "You can now create incremental migrations using: make atlas-migrate-diff SCHEMA=<schema> NAME=<name>"
