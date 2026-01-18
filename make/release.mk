# Release targets
# See docs/guidelines/release.md for detailed release documentation
#
# Quick Example (releasing v6.0.4):
# ---------------------------------
# Step 1: Create and push tag
#   $ make release-tag VERSION=v6.0.4
#   Creating tag v6.0.4...
#   Pushing tag v6.0.4 to origin...
#   ✅ Tag v6.0.4 pushed. GitHub Actions will now run GoReleaser.
#
# Step 2: Monitor GitHub Actions
#   https://github.com/hiromaily/go-crypto-wallet/actions/workflows/release.yml
#
# Or manually:
#   $ git tag -a v6.0.4 -m "Release v6.0.4"
#   $ git push origin v6.0.4

.PHONY: release-check
release-check: ## Check release prerequisites
	@echo "=== Release Prerequisites Check ==="
	@echo ""
	@echo "Current branch:"
	@git branch --show-current
	@echo ""
	@echo "Current version (latest tag):"
	@git describe --tags --abbrev=0 2>/dev/null || echo "(no tags found)"
	@echo ""
	@echo "Commits since last tag:"
	@git log $$(git describe --tags --abbrev=0 2>/dev/null)..HEAD --oneline 2>/dev/null || echo "(no previous tag)"
	@echo ""
	@echo "Uncommitted changes:"
	@git diff-index --quiet HEAD -- && echo "(clean)" || git status --short
	@echo ""
	@echo "=== Checklist ==="
	@echo "[ ] On main branch"
	@echo "[ ] All changes committed and pushed"
	@echo "[ ] CI/CD checks passed"
	@echo "[ ] Version bump determined (major/minor/patch)"

.PHONY: release-dry-run
release-dry-run: ## Dry run GoReleaser locally (requires goreleaser installed)
	goreleaser release --snapshot --clean --skip=publish

.PHONY: release-tag
release-tag: ## Create and push release tag (usage: make release-tag VERSION=v1.2.3)
ifndef VERSION
	@echo "ERROR: VERSION is required"
	@echo "Usage: make release-tag VERSION=v1.2.3"
	@exit 1
endif
	@echo "Creating tag $(VERSION)..."
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@echo "Pushing tag $(VERSION) to origin..."
	@git push origin $(VERSION)
	@echo ""
	@echo "✅ Tag $(VERSION) pushed. GitHub Actions will now run GoReleaser."
	@echo "Monitor: https://github.com/hiromaily/go-crypto-wallet/actions/workflows/release.yml"

.PHONY: release-list
release-list: ## List recent releases
	@echo "=== Recent Tags ==="
	@git tag -l --sort=-version:refname | head -10
	@echo ""
	@echo "=== GitHub Releases ==="
	@gh release list --limit 5 2>/dev/null || echo "(gh CLI not available or not authenticated)"

.PHONY: release-help
release-help: ## Show release documentation
	@echo "Release Process:"
	@echo ""
	@echo "1. Ensure you're on main branch with latest changes"
	@echo "   git checkout main && git pull origin main"
	@echo ""
	@echo "2. Check release prerequisites"
	@echo "   make release-check"
	@echo ""
	@echo "3. (Optional) Dry run locally"
	@echo "   make release-dry-run"
	@echo ""
	@echo "4. Create and push tag"
	@echo "   make release-tag VERSION=vX.Y.Z"
	@echo ""
	@echo "5. Monitor release workflow"
	@echo "   https://github.com/hiromaily/go-crypto-wallet/actions/workflows/release.yml"
	@echo ""
	@echo "Version Guidelines (Semantic Versioning):"
	@echo "  MAJOR: Breaking changes"
	@echo "  MINOR: New features (backward compatible)"
	@echo "  PATCH: Bug fixes"
	@echo ""
	@echo "See docs/guidelines/release.md for detailed documentation."
