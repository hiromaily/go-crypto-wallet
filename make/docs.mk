###############################################################################
# Documentation SSOT (docs-ssot)
###############################################################################

# Install docs-ssot and required dependencies
.PHONY: install-docs
install-docs:
	brew tap hiromaily/tap && brew install hiromaily/tap/docs-ssot jq yq

# Build all documentation outputs from templates
.PHONY: docs
docs:
	docs-ssot build
	docs-ssot index

# Validate that all includes resolve (dry-run, no files written)
.PHONY: docs-validate
docs-validate:
	docs-ssot validate

# Check for SSOT violations (near-duplicate sections)
.PHONY: docs-check
docs-check:
	docs-ssot check
