---
name: docs-update
description: Documentation update workflow. Use when modifying files in docs/ directory or any markdown files (*.md).
---

# Documentation Update Workflow

Workflow for documentation changes.

## Prerequisites

**Use `git-workflow` Skill** for branch, commit, and PR workflow.

## Applicable Files

| Path | Description |
|------|-------------|
| `docs/` | All documentation |
| `*.md` | Root markdown files (README, AGENTS, etc.) |
| `internal/AGENTS.md` | Internal package guidelines |
| `pkg/AGENTS.md` | Package guidelines |

## Guidelines

### Markdown Style

- Use ATX-style headers (`#`, `##`, `###`)
- One blank line between sections
- Code blocks with language specifier
- Tables for structured data

### Documentation Types

| Type | Location | Purpose |
|------|----------|---------|
| AI Agent docs | `docs/ai-agents/` | Agent guidelines |
| Standards | `docs/standards/` | SSOT for standards |
| Architecture | `ARCHITECTURE.md` | System design |
| Crypto-specific | `docs/crypto/` | Chain documentation |

### Links

- Use relative links within docs
- Verify links work after changes
- Update cross-references if moving files

## Verification

No automated verification, but check:

- [ ] Markdown renders correctly
- [ ] Links work
- [ ] No typos or grammar issues
- [ ] Consistent formatting
- [ ] Tables align properly

## Commit Format

```
docs: {brief description}

- {change 1}
- {change 2}

Closes #{issue_number}
```

## Related Skills

- `git-workflow` - Branch, commit, PR workflow
