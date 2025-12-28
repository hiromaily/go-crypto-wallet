
###############################################################################
# Atlas Lint/Format Targets
###############################################################################
.PHONY: fmt-atlas
fmt-atlas: 
	atlas schema fmt tools/atlas/schemas/

.PHONY: lint-atlas
lint-atlas: 
	@atlas schema lint --config file://tools/atlas/atlas.hcl --env local_watch && \
	 atlas schema lint --config file://tools/atlas/atlas.hcl --env local_keygen && \
	 atlas schema lint --config file://tools/atlas/atlas.hcl --env local_sign

###############################################################################
# Atlas Migrate Targets
###############################################################################

# Show migration status for all schemas
.PHONY: atlas-status
atlas-status: 
	@echo "=== Watch Schema ==="
	@atlas migrate status --config file://tools/atlas/atlas.hcl --env local_watch || true
	@echo "\n=== Keygen Schema ==="
	@atlas migrate status --config file://tools/atlas/atlas.hcl --env local_keygen || true
	@echo "\n=== Sign Schema ==="
	@atlas migrate status --config file://tools/atlas/atlas.hcl --env local_sign || true


# Apply all pending migrations for all schemas
.PHONY: atlas-migrate
atlas-migrate:
	@echo "Applying migrations for watch schema..."
	@atlas migrate apply --config file://tools/atlas/atlas.hcl --env local_watch
	@echo "Applying migrations for keygen schema..."
	@atlas migrate apply --config file://tools/atlas/atlas.hcl --env local_keygen
	@echo "Applying migrations for sign schema..."
	@atlas migrate apply --config file://tools/atlas/atlas.hcl --env local_sign
	@echo "All migrations applied successfully!"
