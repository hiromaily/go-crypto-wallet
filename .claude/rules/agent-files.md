---
paths:
  - ".claude/**"
  - ".cursor/**"
  - ".codex/**"
  - ".github/copilot-instructions.md"
---

# Agent Files

This rule applies when editing `.claude/`, `.cursor/`, `.codex/`, or `.github/copilot-instructions.md`.

## Single Source of Truth (SSOT) Design

AI Agent configuration files are managed with **Claude Code's `.claude/` directory as the single source**. Other agent configurations are auto-generated or symlinked.

```
.claude/                    ← SSOT (Source of Truth)
├── commands/               ← Slash commands
├── rules/                  ← Rules (*.md)
└── skills/                 ← Skills (SKILL.md)

.cursor/                    ← Auto-generated / Symlinked
├── commands/README.md      ← Reference only (uses .claude/commands)
├── rules/*.mdc             ← Auto-generated (make sync-cursor-rules)
└── skills/                 ← Symlink → ../.claude/skills

.codex/                     ← TODO (Future support)
.github/copilot-instructions.md ← TODO (Future support)
```

## Supported AI Agents

| Agent | Version | Status | Config Location |
|-------|---------|--------|-----------------|
| Claude Code | v2 | ✅ Active | `.claude/` (SSOT) |
| Cursor | v2 | ✅ Active | `.cursor/` (Auto-generated) |
| Codex | v0.80 | 📋 TODO | `.codex/` |
| GitHub Copilot | 2026 | 📋 TODO | `.github/copilot-instructions.md` |

## File Reference Format

### Claude Code

```markdown
- Follow @AGENTS.md for guidelines
- Refer to @docs/guidelines/coding-conventions.md
```

### Cursor

Supports the same `@path` format. No conversion needed.

## Directory-Specific Rules

### `.claude/commands/`

Location for slash command definitions. Automatically loaded by Cursor as well.

**Editing Notes:**

- Create new commands in `.claude/commands/`
- Only place README.md in `.cursor/commands/`

### `.claude/rules/`

Rule files for Claude Code (`.md`).

**Format:**

```markdown
---
paths:
  - "**/*.go"
  - "**/*.ts"
---

# Rule Title

Rule content...
```

- Without `paths:` → Applies to all instructions (global rule)
- With `paths:` → Applies when editing files matching the specified patterns

### `.claude/skills/`

Location for Skills (MCP) definitions. Each skill has `SKILL.md` in its subdirectory.

```
.claude/skills/
├── go-development/SKILL.md
├── git-workflow/SKILL.md
└── db-migration/SKILL.md
```

### `.cursor/rules/`

**Auto-generated files - DO NOT EDIT DIRECTLY**

Auto-generated from Claude rules. Conversion rules:

| Claude | Cursor |
|--------|--------|
| No `paths:` | `alwaysApply: true` |
| Has `paths:` | `globs:` + `alwaysApply: false` |
| First `# Heading` | `description:` |
| `.md` | `.mdc` |

**Sync Command:**

```bash
make sync-cursor-rules
```

> **IMPORTANT**: After modifying any files in `.claude/rules/`, you MUST run `make sync-cursor-rules` to synchronize with `.cursor/rules/`. This ensures both Claude Code and Cursor use the same rules.

## Model Context Protocol (MCP) / Skills

MCP, proposed since late 2024, is now widespread as of 2026. AI Agents can use repository scripts and local servers as "Skills".

### Claude Code / Cursor

Place `SKILL.md` in `.claude/skills/` or `.cursor/skills/`.

### Codex (Future Support)

For CLI-based Codex agent, Shell/Python scripts are expected to be bound as Skills and invoked like `@database_tool`.

## Editing Checklist

When editing `.claude/` or `.cursor/`:

- [ ] Make changes to `.claude/` (SSOT)
- [ ] Do NOT edit `.cursor/rules/` directly
- [ ] **Run `make sync-cursor-rules` after modifying `.claude/rules/`** (Required)
- [ ] Check if README.md needs updating

## Related Documents

- @.cursor/rules/README.md - Cursor rules specification
- @.cursor/commands/README.md - Commands description
- @AGENTS.md - Project-wide guidelines
