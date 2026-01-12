---
name: shell-scripts
description: Shell script development workflow. Use when modifying files in scripts/ directory or any *.sh files.
---

# Shell Scripts Workflow

Workflow for shell script changes.

## Prerequisites

**Use `git-workflow` Skill** for branch, commit, and PR workflow.

## Applicable Files

| Path | Description |
|------|-------------|
| `scripts/` | All shell scripts |
| `*.sh` | Shell scripts anywhere |

## Verification Commands

```bash
make shfmt    # Format shell scripts
```

### Manual Checks

```bash
# Check syntax
bash -n scripts/{script}.sh

# Shellcheck (if installed)
shellcheck scripts/{script}.sh
```

## Guidelines

### Style

- Use `#!/bin/bash` or `#!/usr/bin/env bash`
- Quote variables: `"${var}"`
- Use `set -euo pipefail` for strict mode
- Add comments for complex logic

### Best Practices

- [ ] Script is executable (`chmod +x`)
- [ ] Has shebang line
- [ ] Uses strict mode
- [ ] Variables are quoted
- [ ] Error handling exists
- [ ] No hardcoded paths (use variables)

### Example Header

```bash
#!/usr/bin/env bash
set -euo pipefail

# Description: What this script does
# Usage: ./script.sh [options]
```

## Verification Checklist

- [ ] `make shfmt` passes
- [ ] Script runs without errors
- [ ] No shellcheck warnings (if available)
- [ ] Proper error handling

## Commit Format

```
chore(scripts): {brief description}

- {change 1}
- {change 2}

Closes #{issue_number}
```

## Related Skills

- `git-workflow` - Branch, commit, PR workflow
