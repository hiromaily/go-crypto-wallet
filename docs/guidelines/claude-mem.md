# Claude-Mem: Persistent Memory for Claude Code

## Overview

[claude-mem](https://github.com/thedotmack/claude-mem) is a plugin that provides persistent memory across Claude Code sessions. It automatically captures what Claude does during coding sessions, compresses it with AI, and injects relevant context back into future sessions.

This complements our existing `.claude/rules/` and `.kiro/steering/` systems by adding **automatic, searchable session history**.

### How It Fits Into Our Memory Stack

| Layer | Source | Scope | Persistence |
|-------|--------|-------|-------------|
| **CLAUDE.md / Rules** | `.claude/rules/`, `CLAUDE.md` | Project conventions & architecture | Git-tracked, shared |
| **Kiro Steering** | `.kiro/steering/` | Product vision, tech stack, structure | Git-tracked, shared |
| **Kiro Specs** | `.kiro/specs/` | Feature-level requirements & design | Git-tracked, shared |
| **Claude-Mem** | `~/.claude-mem/` | Session history & observations | Local per-developer |

**Key distinction**: Rules/Steering/Specs are **declarative knowledge** (what to do). Claude-Mem is **experiential knowledge** (what was done, what worked, what failed).

## Installation

### Prerequisites

- Node.js >= 18.0.0
- Claude Code with plugin support
- Bun and uv (auto-installed if missing)

### Steps

```bash
# 1. In a Claude Code session, run:
/plugin marketplace add thedotmack/claude-mem

# 2. Install the plugin:
/plugin install claude-mem

# 3. Restart Claude Code
```

### Verify Installation

After restarting, claude-mem should automatically:
- Start the worker service on port 37777
- Begin capturing session events
- Provide the `mem-search` skill

You can verify by visiting `http://localhost:37777` in your browser to see the web viewer UI.

## Architecture

```
Claude Code Session
    │
    ├─ SessionStart hook ──────► Worker (port 37777)
    ├─ UserPromptSubmit hook ──►     │
    ├─ PostToolUse hook ───────►     ├─ SQLite DB (sessions, observations)
    ├─ Stop hook ──────────────►     ├─ Chroma Vector DB (semantic search)
    └─ SessionEnd hook ────────►     └─ AI Summarizer
```

### 5 Lifecycle Hooks

| Hook | When | What It Captures |
|------|------|-----------------|
| `SessionStart` | Session begins | Injects relevant past context |
| `UserPromptSubmit` | User sends prompt | Records user intent |
| `PostToolUse` | After each tool call | Records tool usage and results |
| `Stop` | Claude stops responding | Captures completion state |
| `SessionEnd` | Session ends | Generates session summary |

## Usage

### Automatic Operation

Claude-mem works automatically once installed. No manual intervention needed for:
- Capturing observations during sessions
- Generating session summaries
- Injecting relevant context at session start

### Manual Search (mem-search Skill)

Use natural language queries to search past sessions:

```
/mem-search "How did we implement the PostgreSQL repository layer?"
/mem-search "What approach was used for BCH address validation?"
/mem-search "What errors occurred during the XRP integration?"
```

### 3-Layer Progressive Disclosure

Claude-mem uses a token-efficient 3-layer approach:

1. **`search`** - Returns compact index with observation IDs (~50-100 tokens per result)
2. **`timeline`** - Shows chronological context around specific results
3. **`get_observations`** - Fetches full details for filtered IDs (~500-1,000 tokens per result)

This yields ~10x token savings compared to fetching all details upfront.

### Privacy Controls

Wrap sensitive content in `<private>` tags to exclude it from memory:

```
<private>
This content will NOT be stored in claude-mem.
Use for private keys, credentials, or sensitive discussion.
</private>
```

**Important for this project**: Since we handle cryptocurrency private keys, always use `<private>` tags when discussing or working with:
- Private keys and seeds
- Wallet passwords
- API credentials
- Any content that would violate our [security guidelines](./security.md)

## Effective Usage Patterns

### 1. Session Naming

Start sessions with a clear description of your intent. This helps claude-mem generate better summaries and improves future retrieval:

```
# Good - clear intent
"Fix the BCH CashAddr validation bug in address_info.go"

# Less effective - vague
"Fix a bug"
```

### 2. Leverage Past Context for Recurring Tasks

When working on similar tasks across sessions, claude-mem automatically injects relevant past context. This is especially useful for:

- **Chain-specific patterns**: Past BCH/BTC/ETH/XRP work informs future chain tasks
- **Architecture decisions**: Previous design discussions are recalled
- **Debugging history**: Past error resolutions are available

### 3. Cross-Session Knowledge Transfer

Claude-mem bridges the gap between sessions. Example workflow:

```
Session 1: Research and design PostgreSQL migration approach
  → claude-mem captures research findings, design decisions

Session 2: Implement the migration
  → claude-mem injects relevant context from Session 1
  → No need to re-explain the approach
```

### 4. Query Patterns for This Project

Effective search queries for our codebase:

```
# Architecture decisions
/mem-search "clean architecture layer separation decision"

# Chain-specific history
/mem-search "BCH override pattern for Bitcoin struct"
/mem-search "BTC taproot descriptor implementation"

# Build/CI issues
/mem-search "goreleaser configuration changes"
/mem-search "GitHub Actions workflow fixes"

# Database work
/mem-search "SQLC query generation issues"
/mem-search "Atlas migration schema changes"
```

### 5. Combine with Existing Workflow

Claude-mem works alongside our spec-driven development:

```
# Phase 1: Spec creation (captured by claude-mem)
/kiro:spec-init "new-feature"
/kiro:spec-requirements new-feature
/kiro:spec-design new-feature

# Phase 2: Implementation (claude-mem recalls spec discussions)
/kiro:spec-impl new-feature
# claude-mem automatically provides context from Phase 1
```

## Configuration

Settings are stored in `~/.claude-mem/settings.json` (auto-created on first run).

Key settings to consider:

```json
{
  "worker": {
    "port": 37777
  },
  "context": {
    "maxTokens": 2000
  }
}
```

Adjust `maxTokens` based on your needs:
- **Lower (1000-1500)**: Faster sessions, less past context
- **Higher (2500-4000)**: More past context, slightly slower start

## Data Management

### Storage Location

All data is stored locally at `~/.claude-mem/`:
- `sessions.db` - SQLite database
- `chroma/` - Vector database for semantic search

### Viewing Past Sessions

Visit `http://localhost:37777` for the web viewer UI showing:
- Session timeline
- Observations per session
- AI-generated summaries

### Clearing Data

```bash
# Remove all claude-mem data (irreversible)
rm -rf ~/.claude-mem/
```

## Troubleshooting

### Worker Not Starting

```bash
# Check if port 37777 is in use
lsof -i :37777

# Manually start worker (if needed)
claude-mem worker start
```

### No Context Being Injected

1. Verify the plugin is installed: check for claude-mem in `/plugin list`
2. Ensure the worker is running: visit `http://localhost:37777`
3. Check that previous sessions exist in the database

### High Token Usage

If sessions start slowly due to too much injected context:
- Reduce `maxTokens` in settings
- Use `<private>` tags more liberally for non-essential discussions

## Related Documentation

- [Workflow Guidelines](./workflow.md) - Development workflow
- [Security Guidelines](./security.md) - Security requirements (relevant for privacy controls)
- `CLAUDE.md` - Project-level Claude Code configuration
