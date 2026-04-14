## SSOT Structure

**When modifying rules, skills, or documentation, always edit the SSOT location.**

### AI Agent Configuration

| Category | SSOT Location | Other Locations |
|----------|---------------|-----------------|
| Rules | `.claude/rules/*.md` | `.cursor/rules/*.mdc` (auto-generated) |
| Skills | `.claude/skills/*/SKILL.md` | `.cursor/skills/` (symlink) |
| Commands | `.claude/commands/` | `.cursor/commands/` (reference only) |

**Sync Process:**

- `.cursor/rules/` → Auto-generated via `make sync-cursor-rules`
- `.cursor/skills/` → Symlink to `.claude/skills/`

### Project Documentation

| Category | SSOT Location | Notes |
|----------|---------------|-------|
| Guidelines | `docs/guidelines/` | Coding, testing, security, workflow, database, code-generation |
| Architecture | `ARCHITECTURE.md` | System design |
| Agent behavior | `AGENTS.md` (this file) | Entry point for all agents |

### Key Principle

> **Don't Repeat Yourself (DRY)**: Define once, reference everywhere.
> When information exists in multiple places, update the SSOT and reference it from others.
