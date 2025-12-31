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
	@echo "=== Step 3/5: Waiting for DB to be ready ==="
	@echo "Waiting for MySQL to be healthy..."
	@until docker exec wallet-db mysqladmin ping -u root -proot --silent 2>/dev/null; do \
		echo "  Waiting for database..."; \
		sleep 2; \
	done
	@echo "✓ Database is ready"
	@echo ""
	@echo "=== Step 4/5: Extracting sqlc schemas from DB ==="
	@$(MAKE) extract-sqlc-schema-all
	@echo ""
	@echo "=== Step 5/5: Generating sqlc code ==="
	@$(MAKE) sqlc
	@echo ""
	@echo "✓ All done! Schema regeneration complete."
