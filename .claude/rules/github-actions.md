---
paths: [".github/workflows/*.yml"]
---

# GitHub Actions Workflow Rules

## Overview

Rules for modifying GitHub Actions workflow files (`.github/workflows/*.yml`) in go-crypto-wallet.

## Runner Configuration

### Required Runner

**Always use `ubuntu-slim` as the runner** unless there is a specific reason to use another image.

```yaml
# Good
runs-on: ubuntu-slim

# Avoid (unless specifically needed)
runs-on: ubuntu-latest
runs-on: ubuntu-24.04
```

### Why ubuntu-slim?

- Faster startup time (smaller image)
- Lower resource consumption
- Contains essential tools for most CI tasks
- Cost-effective for GitHub Actions minutes

Reference: [GitHub Actions Runner Images](https://github.com/actions/runner-images)

### Available Runners

| Image | YAML Label | Use Case |
|-------|------------|----------|
| Ubuntu Slim | `ubuntu-slim` | **Default choice** - most CI tasks |
| Ubuntu 24.04 | `ubuntu-latest` or `ubuntu-24.04` | When additional tools are required |
| Ubuntu 22.04 | `ubuntu-22.04` | Legacy compatibility |

## Performance Best Practices

### Avoid Duplication

1. **Reusable Workflows**: Extract common patterns into reusable workflows
2. **Composite Actions**: Create composite actions for repeated steps
3. **Matrix Strategy**: Use matrix for parallel execution instead of duplicate jobs

```yaml
# Good: Matrix strategy for parallel execution
strategy:
  matrix:
    go-version: ['1.22', '1.23']

# Avoid: Duplicate jobs for each version
jobs:
  test-go-122:
    ...
  test-go-123:
    ...
```

### Caching

Always use caching for dependencies:

```yaml
- uses: actions/cache@v4
  with:
    path: |
      ~/.cache/go-build
      ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
```

### Conditional Execution

Use path filters to avoid unnecessary runs:

```yaml
on:
  push:
    paths:
      - '**.go'
      - 'go.mod'
      - 'go.sum'
```

## Branch Configuration

### Default Branch

This repository uses `main` as the default branch. Always use `main` (not `master`):

```yaml
# Good
on:
  push:
    branches:
      - main

# Wrong
on:
  push:
    branches:
      - master
```

## Workflow Structure

### Standard Workflow Template

```yaml
name: Workflow Name

on:
  push:
    branches:
      - main
    paths:
      - 'relevant/paths/**'
  pull_request:
    branches:
      - main
  workflow_dispatch:  # Manual trigger

jobs:
  job-name:
    name: Descriptive Job Name
    runs-on: ubuntu-slim
    timeout-minutes: 10  # Always set timeout
    permissions:
      contents: read     # Minimal permissions
    steps:
      - uses: actions/checkout@v4
      # ... steps
```

### Required Elements

| Element | Required | Description |
|---------|----------|-------------|
| `timeout-minutes` | Yes | Prevent runaway jobs |
| `permissions` | Yes | Principle of least privilege |
| `runs-on: ubuntu-slim` | Yes | Standard runner |

## Verification

Before committing workflow changes:

1. Validate YAML syntax
2. Check for branch name (`main` not `master`)
3. Verify runner is `ubuntu-slim`
4. Ensure timeout is set
5. Review permissions (minimal required)

```bash
# Validate YAML syntax
make yaml-lint
```

## Quick Checklist

- [ ] Runner is `ubuntu-slim`
- [ ] Branch is `main` (not `master`)
- [ ] `timeout-minutes` is set
- [ ] `permissions` block is defined
- [ ] Path filters are used where appropriate
- [ ] Caching is implemented for dependencies
- [ ] No duplicate jobs (use matrix instead)

## Related Documentation

- @.github/workflows/ - Existing workflows
- @.github/labels.yml - Repository labels

## Related Skills

- `devops` - CI/CD and DevOps workflow
