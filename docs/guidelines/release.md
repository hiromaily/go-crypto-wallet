# Release Guide

This document describes how to create a new release for go-crypto-wallet.

## Overview

Releases are automated using [GoReleaser](https://goreleaser.com/) and GitHub Actions.
When a tag matching `v*` is pushed, the release workflow automatically:

1. Builds binaries for Linux/macOS (amd64/arm64)
2. Generates changelog from commit messages
3. Creates a GitHub Release with artifacts

## Quick Start

```bash
# 1. Ensure on main branch with latest changes
git checkout main && git pull origin main

# 2. Check prerequisites
make release-check

# 3. Create and push tag
make release-tag VERSION=v6.1.0

# 4. Monitor workflow
# https://github.com/hiromaily/go-crypto-wallet/actions/workflows/release.yml
```

## Detailed Process

### Step 1: Prepare for Release

Ensure all changes are merged to `main`:

```bash
git checkout main
git pull origin main
```

Verify there are no uncommitted changes:

```bash
git status
```

### Step 2: Determine Version

This project follows [Semantic Versioning](https://semver.org/):

| Change Type | Version Bump | Example |
|-------------|--------------|---------|
| Breaking changes | MAJOR | v5.0.0 → v6.0.0 |
| New features (backward compatible) | MINOR | v6.0.0 → v6.1.0 |
| Bug fixes | PATCH | v6.1.0 → v6.1.1 |

Check current version and commits since last release:

```bash
make release-check
```

### Step 3: Create Release Tag

Create an annotated tag and push:

```bash
# Using Makefile (recommended)
make release-tag VERSION=v6.1.0

# Or manually
git tag -a v6.1.0 -m "Release v6.1.0"
git push origin v6.1.0
```

### Step 4: Monitor Release Workflow

The GitHub Actions workflow will automatically run:

1. **Check workflow status**:

   ```bash
   gh run list --workflow=release.yml --limit=3
   ```

2. **Watch live progress**:

   ```bash
   gh run watch <run-id>
   ```

3. **View on GitHub**:
   <https://github.com/hiromaily/go-crypto-wallet/actions/workflows/release.yml>

### Step 5: Verify Release

Once the workflow completes:

1. **Check release page**:
   <https://github.com/hiromaily/go-crypto-wallet/releases>

2. **Verify artifacts**:

   ```bash
   gh release view v6.1.0
   ```

## Release Artifacts

Each release includes:

| Artifact | Description |
|----------|-------------|
| `go-crypto-wallet_X.Y.Z_linux_amd64.tar.gz` | Linux x64 |
| `go-crypto-wallet_X.Y.Z_linux_arm64.tar.gz` | Linux ARM64 |
| `go-crypto-wallet_X.Y.Z_darwin_amd64.tar.gz` | macOS x64 |
| `go-crypto-wallet_X.Y.Z_darwin_arm64.tar.gz` | macOS ARM64 (Apple Silicon) |
| `checksums.txt` | SHA256 checksums |

Each archive contains:

- `watch` - Watch wallet binary (online)
- `keygen` - Keygen wallet binary (offline)
- `sign-auth1` - Sign wallet binary for auth1 (offline)
- `sign-auth2` - Sign wallet binary for auth2 (offline)
- `LICENSE`
- `README.md`
- `ARCHITECTURE.md`

## Changelog Generation

Changelog is automatically generated from commit messages using Conventional Commits:

| Commit Type | Changelog Section |
|-------------|-------------------|
| `feat` | 🚀 Features |
| `fix` | 🐛 Bug Fixes |
| `refactor` | 🔄 Refactoring |
| `deps`, `build` | 📦 Dependencies |
| Other | Others |

Excluded from changelog:

- `docs:` commits
- `test:` commits
- `chore:` commits
- `ci:` commits
- Merge commits

## Makefile Targets

| Target | Description |
|--------|-------------|
| `make release-check` | Check release prerequisites |
| `make release-dry-run` | Dry run GoReleaser locally |
| `make release-tag VERSION=vX.Y.Z` | Create and push release tag |
| `make release-list` | List recent releases |
| `make release-help` | Show release documentation |

## Local Testing

Test GoReleaser locally before pushing a tag:

```bash
# Install GoReleaser (if not installed)
go install github.com/goreleaser/goreleaser/v2@latest

# Dry run (builds but doesn't publish)
make release-dry-run

# Or directly
goreleaser release --snapshot --clean --skip=publish
```

## Troubleshooting

### Workflow Failed

1. Check workflow logs:

   ```bash
   gh run view <run-id> --log-failed
   ```

2. Common issues:
   - **Build errors**: Run `make check-build` locally
   - **CGO issues**: Ensure code is CGO-free compatible
   - **Tag already exists**: Delete and recreate if needed

### Delete a Tag

If you need to delete a tag (e.g., wrong version):

```bash
# Delete local tag
git tag -d v6.1.0

# Delete remote tag
git push origin :refs/tags/v6.1.0
```

**Note**: If a release was created, delete it first on GitHub.

### Re-run Failed Release

If the workflow failed after tag was pushed:

1. Fix the issue in code
2. Delete the tag (local and remote)
3. Create a new commit with the fix
4. Create and push the tag again

## Configuration Files

| File | Purpose |
|------|---------|
| `.goreleaser.yml` | GoReleaser configuration |
| `.github/workflows/release.yml` | GitHub Actions workflow |
| `make/release.mk` | Makefile release targets |

## Related Documentation

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Semantic Versioning](https://semver.org/)
- [GoReleaser Documentation](https://goreleaser.com/)
- [CONTRIBUTING.md](../../CONTRIBUTING.md) - Commit message conventions
