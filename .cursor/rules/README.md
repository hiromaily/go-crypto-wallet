# Cursor Rules

> **Note**: Rules in this directory are auto-generated from `.claude/rules/`.
> Do not edit directly. Edit the source and run `make sync-cursor-rules`.

## Overview

Cursor Rules provide system-level instructions to AI Agents.
Rule contents are included at the beginning of the model context, providing consistent guidance for code generation and editing.

## Rule Types

| Type | frontmatter | Description |
|------|-------------|-------------|
| **Always Apply** | `alwaysApply: true` | Applied to all chat sessions |
| **Apply Intelligently** | `alwaysApply: false` + `description` | Agent determines relevance and applies |
| **Apply to Specific Files** | `globs: ["pattern"]` | Applied when file pattern matches |
| **Apply Manually** | (description only) | Applied when mentioned with `@rule-name` |

## File Format

### RULE.md Format (Recommended)

```
.cursor/rules/
  my-rule/
    RULE.md           # Main rule file
    scripts/          # Helper scripts (optional)
```

### .mdc Format (Legacy, still supported)

```yaml
---
description: Rule description
globs:
  - "**/*.go"
  - "**/*.ts"
alwaysApply: false
---

# Rule body
...
```

## frontmatter Properties

| Property | Type | Description |
|----------|------|-------------|
| `description` | string | Rule description (required for Apply Intelligently) |
| `globs` | string[] | File patterns to apply to |
| `alwaysApply` | boolean | `true`: always apply, `false`: conditional apply |

## Usage in This Project

### SSOT (Single Source of Truth)

- **Source**: `.claude/rules/*.md`
- **Output**: `.cursor/rules/*.mdc`
- **Conversion command**: `make sync-cursor-rules`

### Conversion Rules

| Claude (`paths:`) | Cursor output |
|-------------------|---------------|
| None | `alwaysApply: true` |
| Present | `globs: ...` + `alwaysApply: false` |
| First `# heading` | `description: ...` |
| `.md` extension | `.mdc` extension |

### Sync Commands

```bash
# Verify with dry-run
make sync-cursor-rules-dry

# Actually convert (overwrites existing files)
make sync-cursor-rules
```

## References

- [Cursor Rules Documentation](https://cursor.com/docs/context/rules)
- [AGENTS.md](AGENTS.md) - Project-wide guidelines
