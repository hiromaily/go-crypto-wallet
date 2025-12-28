###############################################################################
# Atlas Migration Targets
###############################################################################

# Check if Atlas CLI is installed
.PHONY: check-atlas
check-atlas:
	@which atlas > /dev/null || (echo "Error: Atlas CLI not found. Install with: brew install arigaio/tap/atlas" && exit 1)

# Apply all pending migrations for all schemas
.PHONY: atlas-migrate
atlas-migrate: check-atlas
	@echo "Applying migrations for watch schema..."
	atlas migrate apply \
		--dir file://tools/atlas/migrations/watch \
		--url "mysql://root:root@127.0.0.1:3306/watch?charset=utf8mb4&parseTime=True&loc=Local"
	@echo "Applying migrations for keygen schema..."
	atlas migrate apply \
		--dir file://tools/atlas/migrations/keygen \
		--url "mysql://root:root@127.0.0.1:3306/keygen?charset=utf8mb4&parseTime=True&loc=Local"
	@echo "Applying migrations for sign schema..."
	atlas migrate apply \
		--dir file://tools/atlas/migrations/sign \
		--url "mysql://root:root@127.0.0.1:3306/sign?charset=utf8mb4&parseTime=True&loc=Local"
	@echo "All migrations applied successfully!"

# Wait for database to be ready
.PHONY: wait-db-ready
wait-db-ready:
	@echo "Waiting for database to be ready..."
	@timeout=60; \
	while [ $$timeout -gt 0 ]; do \
		if docker compose exec -T wallet-db mysqladmin ping -uroot -proot --silent 2>/dev/null; then \
			echo "Database is ready!"; \
			break; \
		fi; \
		echo "Database is not ready yet, waiting... ($$timeout seconds remaining)"; \
		sleep 2; \
		timeout=$$((timeout - 2)); \
	done; \
	if [ $$timeout -le 0 ]; then \
		echo "Error: Database did not become ready within 60 seconds"; \
		exit 1; \
	fi

# Apply migrations for Docker environment using migration service
# Note: Migrations are automatically applied when docker compose up is executed.
# This target is for manual re-application if needed.
.PHONY: atlas-migrate-docker
atlas-migrate-docker: wait-db-ready
	@echo "Applying migrations for watch schema (Docker)..."
	@docker compose run --rm wallet-db-migrate-watch migrate apply \
		--dir file://migrations/watch \
		--url "mysql://root:root@wallet-db:3306/watch?charset=utf8mb4&parseTime=True&loc=Local" || true
	@echo "Applying migrations for keygen schema (Docker)..."
	@docker compose run --rm wallet-db-migrate-keygen migrate apply \
		--dir file://migrations/keygen \
		--url "mysql://root:root@wallet-db:3306/keygen?charset=utf8mb4&parseTime=True&loc=Local" || true
	@echo "Applying migrations for sign schema (Docker)..."
	@docker compose run --rm wallet-db-migrate-sign migrate apply \
		--dir file://migrations/sign \
		--url "mysql://root:root@wallet-db:3306/sign?charset=utf8mb4&parseTime=True&loc=Local" || true
	@echo "All migrations applied successfully!"

# Baseline existing schemas (for databases initialized with SQL files)
# Usage: make atlas-baseline-docker
.PHONY: atlas-baseline-docker
atlas-baseline-docker:
	@echo "Setting baseline for watch schema..."
	@docker compose run --rm wallet-db-migrate-watch migrate apply \
		--baseline 20240101000000 \
		--dir file://migrations/watch \
		--url "mysql://root:root@wallet-db:3306/watch?charset=utf8mb4&parseTime=True&loc=Local" || true
	@echo "Setting baseline for keygen schema..."
	@docker compose run --rm wallet-db-migrate-keygen migrate apply \
		--baseline 20240101000000 \
		--dir file://migrations/keygen \
		--url "mysql://root:root@wallet-db:3306/keygen?charset=utf8mb4&parseTime=True&loc=Local" || true
	@echo "Setting baseline for sign schema..."
	@docker compose run --rm wallet-db-migrate-sign migrate apply \
		--baseline 20240101000000 \
		--dir file://migrations/sign \
		--url "mysql://root:root@wallet-db:3306/sign?charset=utf8mb4&parseTime=True&loc=Local" || true
	@echo "Baseline set for all schemas!"

# Show migration status for all schemas
.PHONY: atlas-status
atlas-status: check-atlas
	@echo "=== Watch Schema ==="
	@atlas migrate status \
		--dir file://tools/atlas/migrations/watch \
		--url "mysql://root:root@127.0.0.1:3306/watch?charset=utf8mb4&parseTime=True&loc=Local" || true
	@echo "\n=== Keygen Schema ==="
	@atlas migrate status \
		--dir file://tools/atlas/migrations/keygen \
		--url "mysql://root:root@127.0.0.1:3306/keygen?charset=utf8mb4&parseTime=True&loc=Local" || true
	@echo "\n=== Sign Schema ==="
	@atlas migrate status \
		--dir file://tools/atlas/migrations/sign \
		--url "mysql://root:root@127.0.0.1:3306/sign?charset=utf8mb4&parseTime=True&loc=Local" || true

# Show migration status for Docker environment
.PHONY: atlas-status-docker
atlas-status-docker:
	@echo "=== Watch Schema ==="
	@docker compose run --rm wallet-db-migrate-watch migrate status \
		--dir file://migrations/watch \
		--url "mysql://root:root@wallet-db:3306/watch?charset=utf8mb4&parseTime=True&loc=Local" || true
	@echo "\n=== Keygen Schema ==="
	@docker compose run --rm wallet-db-migrate-keygen migrate status \
		--dir file://migrations/keygen \
		--url "mysql://root:root@wallet-db:3306/keygen?charset=utf8mb4&parseTime=True&loc=Local" || true
	@echo "\n=== Sign Schema ==="
	@docker compose run --rm wallet-db-migrate-sign migrate status \
		--dir file://migrations/sign \
		--url "mysql://root:root@wallet-db:3306/sign?charset=utf8mb4&parseTime=True&loc=Local" || true

# Rollback last migration for a specific schema
# Usage: make atlas-rollback SCHEMA=watch
.PHONY: atlas-rollback
atlas-rollback: check-atlas
	@if [ -z "$(SCHEMA)" ]; then \
		echo "Error: SCHEMA not specified. Usage: make atlas-rollback SCHEMA=watch"; \
		exit 1; \
	fi
	@echo "Rolling back last migration for $(SCHEMA) schema..."
	atlas migrate down \
		--dir file://tools/atlas/migrations/$(SCHEMA) \
		--url "mysql://root:root@127.0.0.1:3306/$(SCHEMA)?charset=utf8mb4&parseTime=True&loc=Local" 1

# Rollback last migration for Docker environment
# Usage: make atlas-rollback-docker SCHEMA=watch
.PHONY: atlas-rollback-docker
atlas-rollback-docker:
	@if [ -z "$(SCHEMA)" ]; then \
		echo "Error: SCHEMA not specified. Usage: make atlas-rollback-docker SCHEMA=watch"; \
		exit 1; \
	fi
	@echo "Rolling back last migration for $(SCHEMA) schema (Docker)..."
	docker compose run --rm wallet-db-migrate-$(SCHEMA) migrate down \
		--dir file://migrations/$(SCHEMA) \
		--url "mysql://root:root@wallet-db:3306/$(SCHEMA)?charset=utf8mb4&parseTime=True&loc=Local" 1

# Validate migration files
.PHONY: atlas-validate
atlas-validate: check-atlas
	@echo "Validating watch schema migrations..."
	atlas migrate validate \
		--dir file://tools/atlas/migrations/watch
	@echo "Validating keygen schema migrations..."
	atlas migrate validate \
		--dir file://tools/atlas/migrations/keygen
	@echo "Validating sign schema migrations..."
	atlas migrate validate \
		--dir file://tools/atlas/migrations/sign
	@echo "All migrations are valid!"

# Validate migration files for Docker environment
.PHONY: atlas-validate-docker
atlas-validate-docker:
	@echo "Validating watch schema migrations..."
	docker compose run --rm wallet-db-migrate-watch migrate validate \
		--dir file://migrations/watch
	@echo "Validating keygen schema migrations..."
	docker compose run --rm wallet-db-migrate-keygen migrate validate \
		--dir file://migrations/keygen
	@echo "Validating sign schema migrations..."
	docker compose run --rm wallet-db-migrate-sign migrate validate \
		--dir file://migrations/sign
	@echo "All migrations are valid!"

# Generate a new migration file
# Usage: make atlas-new SCHEMA=watch NAME=add_new_table
.PHONY: atlas-new
atlas-new: check-atlas
	@if [ -z "$(SCHEMA)" ] || [ -z "$(NAME)" ]; then \
		echo "Error: SCHEMA and NAME must be specified. Usage: make atlas-new SCHEMA=watch NAME=add_new_table"; \
		exit 1; \
	fi
	@echo "Creating new migration for $(SCHEMA) schema: $(NAME)..."
	atlas migrate new \
		--dir file://tools/atlas/migrations/$(SCHEMA) \
		--name $(NAME)

# Generate a new migration file for Docker environment
# Usage: make atlas-new-docker SCHEMA=watch NAME=add_new_table
.PHONY: atlas-new-docker
atlas-new-docker:
	@if [ -z "$(SCHEMA)" ] || [ -z "$(NAME)" ]; then \
		echo "Error: SCHEMA and NAME must be specified. Usage: make atlas-new-docker SCHEMA=watch NAME=add_new_table"; \
		exit 1; \
	fi
	@echo "Creating new migration for $(SCHEMA) schema: $(NAME) (Docker)..."
	docker compose run --rm wallet-db-migrate-$(SCHEMA) migrate new \
		--dir file://migrations/$(SCHEMA) \
		--name $(NAME)

