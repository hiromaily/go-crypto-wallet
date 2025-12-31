###############################################################################
# Database Targets
###############################################################################

###############################################################################
# Docker Compose Targets
###############################################################################

# run consolidated database
.PHONY: reset-docker-db
reset-docker-db:
	docker compose down -v
	docker compose up


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
# 2. reset-docker-db: Restart DB container to apply migrations
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
	docker compose up -d
	@echo ""
	@echo "=== Step 3/5: Waiting for DB and migrations to be ready ==="
	@# Use 'docker compose wait' to wait for all migration services to complete.
	@# Migration services depend on wallet-db with 'condition: service_healthy',
	@# so MySQL is guaranteed to be ready when migrations complete.
	@# This command blocks until each specified service exits and returns their exit codes.
	@# If any migration fails (non-zero exit), the command will fail and stop the workflow.
	@# Requires Docker Compose v2.21 or later.
	docker compose wait wallet-db-migrate-watch wallet-db-migrate-keygen wallet-db-migrate-sign
	@echo "✓ Database is ready and all migrations completed"
	@echo ""
	@echo "=== Step 4/5: Extracting sqlc schemas from DB ==="
	@$(MAKE) extract-sqlc-schema-all
	@echo ""
	@echo "=== Step 5/5: Generating sqlc code ==="
	@$(MAKE) sqlc
	@echo ""
	@echo "✓ All done! Schema regeneration complete."
