###############################################################################
# Database Targets
###############################################################################

###############################################################################
# DB Dialect Configuration
###############################################################################
# DB_DIALECT selects the database dialect (mysql or postgres)
# Usage: make <target> DB_DIALECT=postgres
# Default: postgres
DB_DIALECT ?= postgres

###############################################################################
# DB Related Targets
###############################################################################
# Atlas-related targets
include make/db_atlas.mk
# sqlc-related targets
include make/db_sqlc.mk

###############################################################################
# Docker Compose Targets
# Usage: make reset-docker [DB_DIALECT=mysql|postgres]
###############################################################################

# Reset and restart database Docker containers
# Usage: make reset-docker [DB_DIALECT=mysql|postgres]
.PHONY: reset-docker
reset-docker:
	docker compose --profile $(DB_DIALECT) down -v
	docker compose --profile $(DB_DIALECT) up

# Convenience aliases for dialect-specific reset
.PHONY: reset-docker-mysql
reset-docker-mysql:
	@$(MAKE) reset-docker DB_DIALECT=mysql

# Remove all of database containers
.PHONY: remove-all-dbs
remove-all-dbs:
	COMPOSE_PROFILES=mysql,postgres docker compose down -v


###############################################################################
# Full Schema Regeneration Workflow
###############################################################################

# Regenerate everything after Atlas HCL schema changes
# This runs all steps needed after modifying tools/atlas/schemas/{db_dialect}/*.hcl files:
# 1. atlas-dev-reset: Regenerate migration files from HCL schemas
# 2. reset-docker: Restart DB container to apply migrations
# 3. Wait for DB to be ready
# 4. extract-sqlc-schema-all: Extract sqlc schemas from running DB
# 5. sqlc: Generate Go code from schemas and queries
#
# Usage: make regenerate-all-from-atlas [DB_DIALECT=mysql|postgres]
.PHONY: regenerate-all-from-atlas
regenerate-all-from-atlas:
	@echo "=== Step 1/5: Regenerating Atlas migrations ($(DB_DIALECT)) ==="
	@$(MAKE) atlas-dev-reset DB_DIALECT=$(DB_DIALECT)
	@echo ""
	@echo "=== Step 2/5: Restarting Docker DB ($(DB_DIALECT)) ==="
	docker compose --profile $(DB_DIALECT) down -v
	docker compose --profile $(DB_DIALECT) up -d
	@echo ""
	@echo "=== Step 3/5: Waiting for DB and migrations to be ready ==="
	@# Use 'docker compose wait' to wait for all migration services to complete.
	@# Migration services depend on the DB with 'condition: service_healthy',
	@# so the DB is guaranteed to be ready when migrations complete.
	@# This command blocks until each specified service exits and returns their exit codes.
	@# If any migration fails (non-zero exit), the command will fail and stop the workflow.
	@# Requires Docker Compose v2.21 or later.
	docker compose wait wallet-$(DB_DIALECT)-migrate-watch wallet-$(DB_DIALECT)-migrate-keygen wallet-$(DB_DIALECT)-migrate-sign
	@echo "✓ Database is ready and all migrations completed"
	@echo ""
	@echo "=== Step 4/5: Extracting sqlc schemas from DB ==="
	@$(MAKE) extract-sqlc-schema-all DB_DIALECT=$(DB_DIALECT)
	@echo ""
	@echo "=== Step 5/5: Generating sqlc code ==="
	@$(MAKE) sqlc-$(DB_DIALECT)
	@echo ""
	@echo "✓ All done! Schema regeneration complete ($(DB_DIALECT))."

# Convenience aliases for dialect-specific regeneration
.PHONY: regenerate-all-from-atlas-mysql
regenerate-all-from-atlas-mysql:
	@$(MAKE) regenerate-all-from-atlas DB_DIALECT=mysql


###############################################################################
# sqlfluff Linting
###############################################################################
###############################################################################
# SQLFluff Targets (SQL Formatting and Linting)
###############################################################################
# SQLFluff is used for formatting and linting SQL files used by sqlc
# Note: SQLFluff may show PRS (parsing) errors for MySQL ? placeholders,
# but these are acceptable as sqlc handles them correctly.

# Format SQL files for all dialects
.PHONY: sqlfluff-format
sqlfluff-format:
	@for dialect in mysql postgres; do \
		echo "Formatting SQL files ($$dialect)..."; \
		sqlfluff format tools/sqlc/queries/$$dialect/*.sql; \
	done
	@echo "✓ SQL files formatted"

# Lint SQL files for all dialects
.PHONY: sqlfluff-lint
sqlfluff-lint:
	@for dialect in mysql postgres; do \
		echo "Linting SQL files ($$dialect)..."; \
		sqlfluff lint tools/sqlc/queries/$$dialect/*.sql || true; \
	done
	@echo "Note: PRS (parsing) errors for ? placeholders are acceptable for sqlc"

# Fix SQL files for all dialects (format and auto-fix issues)
.PHONY: sqlfluff-fix
sqlfluff-fix:
	@for dialect in mysql postgres; do \
		echo "Formatting and fixing SQL files ($$dialect)..."; \
		sqlfluff fix tools/sqlc/queries/$$dialect/*.sql; \
	done
	@echo "✓ SQL files formatted and fixed"
