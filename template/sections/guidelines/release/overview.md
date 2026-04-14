## Release Guide

This document describes how to create a new release for go-crypto-wallet.

### Overview

Releases are automated using [GoReleaser](https://goreleaser.com/) and GitHub Actions.
When a tag matching `v*` is pushed, the release workflow automatically:

1. Builds binaries for Linux/macOS (amd64/arm64)
2. Generates changelog from commit messages
3. Creates a GitHub Release with artifacts

### Quick Start

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
