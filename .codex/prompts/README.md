# Codex Prompts

> **TODO**: Codex CLI has unique characteristics different from Cursor and Claude Code.
> Prompts configuration is pending until the optimal approach is determined.

## Differences from Cursor / Claude Code

| Feature | Codex CLI | Cursor | Claude Code |
|---------|-----------|--------|-------------|
| Config format | `config.toml` | `.mdc` files | `.md` files |
| Prompts | CLI args / interactive | `/commands/` | `/commands/` |
| Skills/Rules | `.codex/rules/` | `.cursor/rules/` | `.claude/skills/` |
| Execution modes | Suggest/Auto Edit/Full Auto | N/A | N/A |
| Sandbox | Network-disabled sandbox | IDE sandbox | Terminal sandbox |

## Codex-Specific Features

### Approval Modes

| Mode | Description |
|------|-------------|
| Suggest (default) | Proposes changes, requires approval |
| Auto Edit | Auto file edits, approval for commands |
| Full Auto | Autonomous in sandboxed environment |

### Usage

```bash
# Suggest mode (default)
codex "describe this function"

# Auto edit mode
codex --auto-edit "refactor this module"

# Full auto mode
codex --full-auto "implement feature"
```

### Configuration

Codex uses `config.toml` instead of markdown prompts:

```toml
# .codex/config.toml
model = "gpt-5-codex"
approval_mode = "suggest"
```

## Current Status

For now, Codex should use:

- `.codex/rules/general.md` - General guidelines (references docs/guidelines/)

Prompt-based workflows (like `/fix-issue`) are **not yet implemented** for Codex.

## References

- [Codex CLI Documentation](https://developers.openai.com/codex/)
- [Codex Configuration](https://developers.openai.com/codex/config-advanced/)
