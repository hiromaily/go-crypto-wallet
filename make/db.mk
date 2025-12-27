###############################################################################
# Database Targets
###############################################################################

###############################################################################
# Docker Compose Targets
###############################################################################

# run consolidated database
.PHONY: up-docker-db
up-docker-db:
	docker compose up wallet-db

# remove database volume
.PHONY: rm-db-volumes
rm-db-volumes:
	docker volume rm -f go-crypto-wallet_wallet-db

###############################################################################
# Schema Export Targets
###############################################################################

# Export watch schema from wallet-db container
.PHONY: dump-schema-watch
dump-schema-watch:
	mkdir -p data/dump/sql
	docker exec wallet-db mysqldump -u root -proot --no-data --skip-triggers watch > data/dump/sql/dump_watch.sql

# Export keygen schema from wallet-db container
.PHONY: dump-schema-keygen
dump-schema-keygen:
	mkdir -p data/dump/sql
	docker exec wallet-db mysqldump -u root -proot --no-data --skip-triggers keygen > data/dump/sql/dump_keygen.sql

# Export sign schema from wallet-db container
.PHONY: dump-schema-sign
dump-schema-sign:
	mkdir -p data/dump/sql
	docker exec wallet-db mysqldump -u root -proot --no-data --skip-triggers sign > data/dump/sql/dump_sign.sql

# Export all schemas from wallet-db container
.PHONY: dump-schema-all
dump-schema-all: dump-schema-watch dump-schema-keygen dump-schema-sign

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

# Apply migrations for Docker environment
.PHONY: atlas-migrate-docker
atlas-migrate-docker: check-atlas
	@echo "Applying migrations for watch schema (Docker)..."
	atlas migrate apply \
		--dir file://tools/atlas/migrations/watch \
		--url "mysql://root:root@wallet-db:3306/watch?charset=utf8mb4&parseTime=True&loc=Local"
	@echo "Applying migrations for keygen schema (Docker)..."
	atlas migrate apply \
		--dir file://tools/atlas/migrations/keygen \
		--url "mysql://root:root@wallet-db:3306/keygen?charset=utf8mb4&parseTime=True&loc=Local"
	@echo "Applying migrations for sign schema (Docker)..."
	atlas migrate apply \
		--dir file://tools/atlas/migrations/sign \
		--url "mysql://root:root@wallet-db:3306/sign?charset=utf8mb4&parseTime=True&loc=Local"
	@echo "All migrations applied successfully!"

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

# Apply HCL schema files to database
# This applies the declarative HCL schema definitions
.PHONY: atlas-schema-apply
atlas-schema-apply: check-atlas
	@echo "Applying watch schema from HCL..."
	cd tools/atlas && atlas schema apply \
		--url "mysql://root:root@127.0.0.1:3306/watch?charset=utf8mb4&parseTime=True&loc=Local" \
		--to file://schemas/watch.hcl \
		--auto-approve
	@echo "Applying keygen schema from HCL..."
	cd tools/atlas && atlas schema apply \
		--url "mysql://root:root@127.0.0.1:3306/keygen?charset=utf8mb4&parseTime=True&loc=Local" \
		--to file://schemas/keygen.hcl \
		--auto-approve
	@echo "Applying sign schema from HCL..."
	cd tools/atlas && atlas schema apply \
		--url "mysql://root:root@127.0.0.1:3306/sign?charset=utf8mb4&parseTime=True&loc=Local" \
		--to file://schemas/sign.hcl \
		--auto-approve
	@echo "All HCL schemas applied successfully!"

# Show diff between HCL schema and database
# Usage: make atlas-schema-diff SCHEMA=watch
.PHONY: atlas-schema-diff
atlas-schema-diff: check-atlas
	@if [ -z "$(SCHEMA)" ]; then \
		echo "Error: SCHEMA not specified. Usage: make atlas-schema-diff SCHEMA=watch"; \
		exit 1; \
	fi
	@echo "Showing diff for $(SCHEMA) schema..."
	cd tools/atlas && atlas schema diff \
		--from "mysql://root:root@127.0.0.1:3306/$(SCHEMA)?charset=utf8mb4&parseTime=True&loc=Local" \
		--to file://schemas/$(SCHEMA).hcl

# Generate migration from HCL schema diff
# Usage: make atlas-schema-diff-migration SCHEMA=watch
.PHONY: atlas-schema-diff-migration
atlas-schema-diff-migration: check-atlas
	@if [ -z "$(SCHEMA)" ]; then \
		echo "Error: SCHEMA not specified. Usage: make atlas-schema-diff-migration SCHEMA=watch"; \
		exit 1; \
	fi
	@echo "Generating migration from HCL schema diff for $(SCHEMA)..."
	cd tools/atlas && atlas migrate diff \
		--dir file://migrations/$(SCHEMA) \
		--to file://schemas/$(SCHEMA).hcl \
		--from "mysql://root:root@127.0.0.1:3306/$(SCHEMA)?charset=utf8mb4&parseTime=True&loc=Local"
