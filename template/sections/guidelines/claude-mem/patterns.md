### Effective Usage Patterns

#### 1. Session Naming

Start sessions with a clear description of your intent. This helps claude-mem generate better summaries and improves future retrieval:

```
# Good - clear intent
"Fix the BCH CashAddr validation bug in address_info.go"

# Less effective - vague
"Fix a bug"
```

#### 2. Leverage Past Context for Recurring Tasks

When working on similar tasks across sessions, claude-mem automatically injects relevant past context. This is especially useful for:

- **Chain-specific patterns**: Past BCH/BTC/ETH/XRP work informs future chain tasks
- **Architecture decisions**: Previous design discussions are recalled
- **Debugging history**: Past error resolutions are available

#### 3. Cross-Session Knowledge Transfer

Claude-mem bridges the gap between sessions. Example workflow:

```
Session 1: Research and design PostgreSQL migration approach
  → claude-mem captures research findings, design decisions

Session 2: Implement the migration
  → claude-mem injects relevant context from Session 1
  → No need to re-explain the approach
```

#### 4. Query Patterns for This Project

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

#### 5. Combine with Existing Workflow

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
