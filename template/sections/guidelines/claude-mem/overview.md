## Claude-Mem: Persistent Memory for Claude Code

### Overview

[claude-mem](https://github.com/thedotmack/claude-mem) is a plugin that provides persistent memory across Claude Code sessions. It automatically captures what Claude does during coding sessions, compresses it with AI, and injects relevant context back into future sessions.

This complements our existing `.claude/rules/` and `.kiro/steering/` systems by adding **automatic, searchable session history**.

#### How It Fits Into Our Memory Stack

| Layer | Source | Scope | Persistence |
|-------|--------|-------|-------------|
| **CLAUDE.md / Rules** | `.claude/rules/`, `CLAUDE.md` | Project conventions & architecture | Git-tracked, shared |
| **Kiro Steering** | `.kiro/steering/` | Product vision, tech stack, structure | Git-tracked, shared |
| **Kiro Specs** | `.kiro/specs/` | Feature-level requirements & design | Git-tracked, shared |
| **Claude-Mem** | `~/.claude-mem/` | Session history & observations | Local per-developer |

**Key distinction**: Rules/Steering/Specs are **declarative knowledge** (what to do). Claude-Mem is **experiential knowledge** (what was done, what worked, what failed).
