---
paths: ["**/*.yaml", "**/*.yml"]
---

# YAML File Rules

## Overview

Rules for modifying YAML files (`*.yaml`, `*.yml`) in go-crypto-wallet.

## Applicable Directories

| Directory | Description |
|-----------|-------------|
| `config/` | Application configuration files |
| `.github/workflows/` | GitHub Actions workflows |
| `.devcontainer/` | Dev container configuration |
| Root (`*.yaml`) | Docker Compose, tool configurations |

## Verification Commands

| Command | Purpose | Required |
|---------|---------|----------|
| `make yaml-lint` | Lint YAML files | Recommended |

### Target Directories for yaml-lint

The `yaml-lint` command checks:

- `.github/workflows/`
- `.devcontainer/`
- `config/`

## Configuration File Types

### Application Config (`config/`)

Wallet configuration files for BTC, BCH, ETH, XRP:

```yaml
# Example: config/wallet/btc_watch.yaml
bitcoin:
  network: testnet
  host: "127.0.0.1:18332"
```

### Docker Compose (`compose*.yaml`)

Service definitions:

| File | Description |
|------|-------------|
| `compose.yaml` | Base services (database, etc.) |
| `compose.btc.yaml` | Bitcoin services |
| `compose.bch.yaml` | Bitcoin Cash services |
| `compose.eth.yaml` | Ethereum services |
| `compose.xrp.yaml` | Ripple services |

### GitHub Actions (`.github/workflows/`)

CI/CD workflow definitions.

## Validation

### Manual Validation

```bash
# Validate YAML syntax
python -c "import yaml; yaml.safe_load(open('file.yaml'))"

# Or using yq
yq eval '.' file.yaml > /dev/null
```

### Docker Compose Validation

```bash
docker compose config
```

## Best Practices

### Formatting

- Use 2 spaces for indentation
- Use quotes for strings that could be misinterpreted
- Keep lines under 120 characters

### Structure

```yaml
# Good: Explicit quotes for version strings
version: "3.8"

# Good: Explicit quotes for port mappings
ports:
  - "8080:80"

# Good: Anchors for repeated values
defaults: &defaults
  restart: unless-stopped

services:
  app:
    <<: *defaults
```

## Quick Checklist

- [ ] Valid YAML syntax
- [ ] Consistent indentation (2 spaces)
- [ ] Sensitive values use environment variables, not hardcoded
- [ ] Port mappings are quoted
- [ ] Comments explain non-obvious configurations

## Related Documentation

- @compose.yaml - Main Docker Compose configuration
- @config/ - Application configuration files
- @.github/workflows/ - CI/CD workflows

## Related Skills

- `devops` - CI/CD and DevOps workflow
