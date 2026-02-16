###############################################################################
# Database Targets
###############################################################################

###############################################################################
# Docker Compose Targets
###############################################################################

# run consolidated database
.PHONY: reset-docker-mysql
reset-docker-mysql:
	docker compose down -v
	docker compose --profile mysql up

.PHONY: reset-docker-postgres
reset-docker-postgres:
	docker compose down -v
	docker compose --profile postgres up

###############################################################################
# DB Related Targets
###############################################################################
# Atlas-related targets
include make/db_atlas.mk
# sqlc-related targets
include make/db_sqlc.mk

###############################################################################
# Full Schema Regeneration Workflow
###############################################################################

# Regenerate everything after Atlas HCL schema changes
# This runs all steps needed after modifying tools/atlas/schemas/*.hcl files:
# 1. atlas-dev-reset: Regenerate migration files from HCL schemas
# 2. reset-docker-mysql: Restart DB container to apply migrations
# 3. Wait for DB to be ready
# 4. extract-sqlc-schema-all: Extract sqlc schemas from running DB
# 5. sqlc: Generate Go code from schemas and queries
#
# Usage: After modifying tools/atlas/schemas/*.hcl, run:
#   make regenerate-all-from-atlas
.PHONY: regenerate-all-from-atlas
regenerate-all-from-atlas:
	@echo "=== Step 1/5: Regenerating Atlas migrations ==="
	@$(MAKE) atlas-dev-reset
	@echo ""
	@echo "=== Step 2/5: Restarting Docker DB ==="
	docker compose down -v
	docker compose --profile mysql up -d
	@echo ""
	@echo "=== Step 3/5: Waiting for DB and migrations to be ready ==="
	@# Use 'docker compose wait' to wait for all migration services to complete.
	@# Migration services depend on wallet-mysql with 'condition: service_healthy',
	@# so MySQL is guaranteed to be ready when migrations complete.
	@# This command blocks until each specified service exits and returns their exit codes.
	@# If any migration fails (non-zero exit), the command will fail and stop the workflow.
	@# Requires Docker Compose v2.21 or later.
	docker compose wait wallet-mysql-migrate-watch wallet-mysql-migrate-keygen wallet-mysql-migrate-sign
	@echo "✓ Database is ready and all migrations completed"
	@echo ""
	@echo "=== Step 4/5: Extracting sqlc schemas from DB ==="
	@$(MAKE) extract-sqlc-schema-all
	@echo ""
	@echo "=== Step 5/5: Generating sqlc code ==="
	@$(MAKE) sqlc
	@echo ""
	@echo "✓ All done! Schema regeneration complete."


###############################################################################
# sqlfluff Linting
###############################################################################

###############################################################################
# SQLFluff Targets (SQL Formatting and Linting)
###############################################################################
# SQLFluff is used for formatting and linting SQL files used by sqlc
# Note: SQLFluff may show PRS (parsing) errors for MySQL ? placeholders,
# but these are acceptable as sqlc handles them correctly.

# Format SQL files
# Formats SQL files according to .sqlfluff configuration
.PHONY: sqlfluff-format
sqlfluff-format:
	@echo "Formatting SQL files..."
	@sqlfluff format tools/sqlc/queries/*.sql
	@echo "✓ SQL files formatted"

# Lint SQL files
# Lints SQL files and reports issues (excluding parsing errors for ? placeholders)
.PHONY: sqlfluff-lint
sqlfluff-lint:
	@echo "Linting SQL files..."
	@sqlfluff lint tools/sqlc/queries/*.sql || true
	@echo "Note: PRS (parsing) errors for ? placeholders are acceptable for sqlc"

# Fix SQL files (format and auto-fix issues)
# Formats SQL files and automatically fixes linting issues where possible
.PHONY: sqlfluff-fix
sqlfluff-fix:
	@echo "Formatting and fixing SQL files..."
	@sqlfluff fix tools/sqlc/queries/*.sql
	@echo "✓ SQL files formatted and fixed"
