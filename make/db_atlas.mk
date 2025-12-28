
###############################################################################
# Atlas Configuration
###############################################################################
# Atlas configuration file and environments
ATLAS_CONFIG := file://tools/atlas/atlas.hcl
ATLAS_ENV_WATCH := local_watch
ATLAS_ENV_KEYGEN := local_keygen
ATLAS_ENV_SIGN := local_sign

###############################################################################
# Atlas Format and Lint Targets
###############################################################################

# Format all HCL schema files
.PHONY: atlas-fmt
atlas-fmt:
	@echo "Formatting Atlas HCL schema files..."
	@atlas schema fmt tools/atlas/schemas/
	@echo "✓ HCL files formatted successfully"

# Lint all HCL schema files with validation
.PHONY: atlas-lint
atlas-lint:
	@echo "Linting Atlas HCL schema files..."
	@echo "=== Linting watch schema ==="
	@atlas schema lint --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_WATCH)
	@echo "\n=== Linting keygen schema ==="
	@atlas schema lint --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_KEYGEN)
	@echo "\n=== Linting sign schema ==="
	@atlas schema lint --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_SIGN)
	@echo "✓ All schemas passed linting"

###############################################################################
# Atlas Schema Management Targets
###############################################################################

# Show diff between database and HCL schema for a specific schema
# Usage: make atlas-schema-diff SCHEMA=watch
.PHONY: atlas-schema-diff
atlas-schema-diff:
ifndef SCHEMA
	@echo "Error: SCHEMA parameter is required"
	@echo "Usage: make atlas-schema-diff SCHEMA=watch|keygen|sign"
	@exit 1
endif
	@echo "=== Schema diff for $(SCHEMA) ==="
	@atlas schema diff --config $(ATLAS_CONFIG) --env local_$(SCHEMA)

# Show diff for all schemas
.PHONY: atlas-schema-diff-all
atlas-schema-diff-all:
	@echo "=== Watch Schema Diff ==="
	@atlas schema diff --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_WATCH) || true
	@echo "\n=== Keygen Schema Diff ==="
	@atlas schema diff --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_KEYGEN) || true
	@echo "\n=== Sign Schema Diff ==="
	@atlas schema diff --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_SIGN) || true

# Apply HCL schema directly to database (all schemas)
.PHONY: atlas-schema-apply
atlas-schema-apply:
	@echo "Applying HCL schemas directly to databases..."
	@echo "=== Applying watch schema ==="
	@atlas schema apply --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_WATCH) --auto-approve
	@echo "=== Applying keygen schema ==="
	@atlas schema apply --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_KEYGEN) --auto-approve
	@echo "=== Applying sign schema ==="
	@atlas schema apply --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_SIGN) --auto-approve
	@echo "✓ All schemas applied successfully"

# Apply HCL schema for a specific schema
# Usage: make atlas-schema-apply-one SCHEMA=watch
.PHONY: atlas-schema-apply-one
atlas-schema-apply-one:
ifndef SCHEMA
	@echo "Error: SCHEMA parameter is required"
	@echo "Usage: make atlas-schema-apply-one SCHEMA=watch|keygen|sign"
	@exit 1
endif
	@echo "=== Applying $(SCHEMA) schema ==="
	@atlas schema apply --config $(ATLAS_CONFIG) --env local_$(SCHEMA) --auto-approve

###############################################################################
# Atlas Migration Targets
###############################################################################

# Show migration status for all schemas
.PHONY: atlas-migrate-status
atlas-migrate-status:
	@echo "=== Watch Schema Migration Status ==="
	@atlas migrate status --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_WATCH) || true
	@echo "\n=== Keygen Schema Migration Status ==="
	@atlas migrate status --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_KEYGEN) || true
	@echo "\n=== Sign Schema Migration Status ==="
	@atlas migrate status --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_SIGN) || true

# Apply all pending migrations for all schemas
.PHONY: atlas-migrate-apply
atlas-migrate-apply:
	@echo "Applying migrations for all schemas..."
	@echo "=== Applying watch migrations ==="
	@atlas migrate apply --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_WATCH)
	@echo "=== Applying keygen migrations ==="
	@atlas migrate apply --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_KEYGEN)
	@echo "=== Applying sign migrations ==="
	@atlas migrate apply --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_SIGN)
	@echo "✓ All migrations applied successfully"

# Generate new migration from HCL schema diff
# Usage: make atlas-migrate-diff SCHEMA=watch NAME=add_new_column
.PHONY: atlas-migrate-diff
atlas-migrate-diff:
ifndef SCHEMA
	@echo "Error: SCHEMA parameter is required"
	@echo "Usage: make atlas-migrate-diff SCHEMA=watch|keygen|sign NAME=migration_name"
	@exit 1
endif
ifndef NAME
	@echo "Error: NAME parameter is required"
	@echo "Usage: make atlas-migrate-diff SCHEMA=watch|keygen|sign NAME=migration_name"
	@exit 1
endif
	@echo "=== Generating migration for $(SCHEMA) schema ==="
	@atlas migrate diff $(NAME) --config $(ATLAS_CONFIG) --env local_$(SCHEMA)
	@echo "✓ Migration generated: tools/atlas/migrations/$(SCHEMA)/"

# Hash migration directory for production readiness
.PHONY: atlas-migrate-hash
atlas-migrate-hash:
	@echo "Hashing migration directories..."
	@atlas migrate hash --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_WATCH)
	@atlas migrate hash --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_KEYGEN)
	@atlas migrate hash --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_SIGN)
	@echo "✓ Migration directories hashed"

###############################################################################
# Development Workflow Targets
###############################################################################

# Development Mode: Regenerate migrations from HCL schemas (from scratch)
# WARNING: This deletes all existing migrations and creates new ones
.PHONY: atlas-dev-reset
atlas-dev-reset:
	@echo "⚠️  WARNING: This will delete all existing migrations!"
	@echo "⚠️  This is intended for development only."
	@read -p "Are you sure you want to continue? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo "Regenerating migrations from HCL schemas..."; \
		rm -rf tools/atlas/migrations/watch/*.sql; \
		rm -rf tools/atlas/migrations/keygen/*.sql; \
		rm -rf tools/atlas/migrations/sign/*.sql; \
		echo "=== Generating watch migration ==="; \
		atlas migrate diff initial_schema --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_WATCH); \
		echo "=== Generating keygen migration ==="; \
		atlas migrate diff initial_schema --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_KEYGEN); \
		echo "=== Generating sign migration ==="; \
		atlas migrate diff initial_schema --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_SIGN); \
		echo "✓ Migrations regenerated. Run 'docker compose up' or 'make atlas-migrate-apply' to apply."; \
	else \
		echo "Cancelled."; \
		exit 1; \
	fi

# Development Mode: Clean databases and reapply from HCL schemas
.PHONY: atlas-dev-clean
atlas-dev-clean:
	@echo "⚠️  WARNING: This will drop all databases and recreate them from HCL schemas!"
	@read -p "Are you sure you want to continue? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo "Cleaning and recreating databases..."; \
		atlas schema clean --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_WATCH) --auto-approve; \
		atlas schema clean --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_KEYGEN) --auto-approve; \
		atlas schema clean --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_SIGN) --auto-approve; \
		echo "Applying HCL schemas..."; \
		$(MAKE) atlas-schema-apply; \
		echo "✓ Databases cleaned and recreated"; \
	else \
		echo "Cancelled."; \
		exit 1; \
	fi

###############################################################################
# Production Workflow Targets
###############################################################################

# Production Mode: Initialize migration history from current database state
.PHONY: atlas-prod-init
atlas-prod-init:
	@echo "Initializing production migration history..."
	@echo "=== Hashing watch migrations ==="
	@atlas migrate hash --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_WATCH)
	@echo "=== Hashing keygen migrations ==="
	@atlas migrate hash --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_KEYGEN)
	@echo "=== Hashing sign migrations ==="
	@atlas migrate hash --config $(ATLAS_CONFIG) --env $(ATLAS_ENV_SIGN)
	@echo "✓ Production migration history initialized"
	@echo "You can now create incremental migrations using: make atlas-migrate-diff SCHEMA=<schema> NAME=<name>"

###############################################################################
# Utility Targets
###############################################################################

# Validate Atlas configuration
.PHONY: atlas-validate
atlas-validate:
	@echo "Validating Atlas configuration..."
	@cd tools/atlas && atlas migrate validate --config file://atlas.hcl --env $(ATLAS_ENV_WATCH)
	@cd tools/atlas && atlas migrate validate --config file://atlas.hcl --env $(ATLAS_ENV_KEYGEN)
	@cd tools/atlas && atlas migrate validate --config file://atlas.hcl --env $(ATLAS_ENV_SIGN)
	@echo "✓ Atlas configuration is valid"

# Show Atlas help
.PHONY: atlas-help
atlas-help:
	@echo "Atlas Development Workflow Targets:"
	@echo ""
	@echo "Format and Lint:"
	@echo "  atlas-fmt                Format HCL schema files"
	@echo "  atlas-lint               Lint HCL schema files"
	@echo ""
	@echo "Schema Management:"
	@echo "  atlas-schema-diff SCHEMA=<name>    Show diff for specific schema"
	@echo "  atlas-schema-diff-all              Show diff for all schemas"
	@echo "  atlas-schema-apply                 Apply HCL schemas to all databases"
	@echo "  atlas-schema-apply-one SCHEMA=<name>  Apply HCL schema to specific database"
	@echo ""
	@echo "Migration Management:"
	@echo "  atlas-migrate-status               Show migration status for all schemas"
	@echo "  atlas-migrate-apply                Apply all pending migrations"
	@echo "  atlas-migrate-diff SCHEMA=<name> NAME=<name>  Generate new migration"
	@echo "  atlas-migrate-hash                 Hash migrations for production"
	@echo ""
	@echo "Development Workflow:"
	@echo "  atlas-dev-reset          Regenerate migrations from scratch (deletes existing)"
	@echo "  atlas-dev-clean          Clean databases and recreate from HCL schemas"
	@echo ""
	@echo "Production Workflow:"
	@echo "  atlas-prod-init          Initialize migration history for production"
	@echo ""
	@echo "Utilities:"
	@echo "  atlas-validate           Validate Atlas configuration"
	@echo "  atlas-help               Show this help message"

# Backward compatibility aliases
.PHONY: fmt-atlas lint-atlas atlas-status atlas-migrate
fmt-atlas: atlas-fmt
lint-atlas: atlas-lint
atlas-status: atlas-migrate-status
atlas-migrate: atlas-migrate-apply
