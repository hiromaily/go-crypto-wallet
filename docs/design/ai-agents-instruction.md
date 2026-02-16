# AI Agent Instruction System Design

This document describes the design philosophy and structure of the AI agent instruction system in this project.

## Overview

This project uses a layered instruction system to guide AI agents (Claude, Cursor) in performing tasks consistently and correctly. The system is designed around the principles of:

1. **Single Source of Truth (SSOT)** - Each piece of information is defined in one authoritative location
2. **Separation of Concerns** - Commands, Skills, Rules, and Contexts have distinct responsibilities
3. **Composability** - Components can be combined to handle various task types

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         User Request                                │
└─────────────────────────────────┬───────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                      Slash Commands                                 │
│                   (.claude/commands/)                               │
│                                                                     │
│  Entry points for specific workflows                                │
│  Examples: /fix-issue, /fix-linter, /fix-pr-review                  │
└─────────────────────────────────┬───────────────────────────────────┘
                                  │ invokes
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         Skills                                      │
│                    (.claude/skills/)                                │
│                                                                     │
│  Reusable expertise for specific domains                            │
│  Examples: go-development, git-workflow, label-context-mapping      │
└─────────────────────────────────┬───────────────────────────────────┘
                                  │ references
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                         Rules                                       │
│                    (.claude/rules/)                                 │
│                                                                     │
│  Always-applied constraints and behaviors                           │
│  Examples: security.md, task-context-loading.md                     │
└─────────────────────────────────┬───────────────────────────────────┘
                                  │ loads
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    Context Documents                                │
│                   (docs/task-contexts/)                             │
│                                                                     │
│  Domain-specific knowledge and procedures                           │
│  Examples: bug-fix.md, chains/btc.md                                │
└─────────────────────────────────────────────────────────────────────┘
```

## Component Definitions

### Slash Commands

**Location**: `.claude/commands/`

**Purpose**: Entry points that trigger specific workflows. Commands are user-facing and provide a simple interface to complex operations.

**Characteristics**:

- Triggered explicitly by user (e.g., `/fix-issue #123`)
- Orchestrate the loading of appropriate skills and contexts
- Should be minimal and delegate to skills

**Examples**:

| Command | Purpose |
|---------|---------|
| `/fix-issue` | Work on a GitHub issue |
| `/fix-linter` | Fix linter errors |
| `/fix-pr-review` | Address PR review comments |

### Skills

**Location**: `.claude/skills/`

**Purpose**: Encapsulate expertise for specific domains or workflows. Skills are reusable building blocks that can be composed.

**Characteristics**:

- Loaded on-demand by commands or other skills
- Contain domain-specific knowledge and procedures
- Can reference other skills (composition)
- Define verification commands and checklists

**Categories**:

| Category | Skills | Description |
|----------|--------|-------------|
| **Development** | `go-development`, `typescript-development`, `solidity-development` | Language-specific workflows |
| **Infrastructure** | `db-migration`, `devops`, `shell-scripts`, `makefile-update` | Non-code development |
| **Documentation** | `docs-update` | Documentation workflows |
| **Workflow** | `git-workflow`, `fix-pr-review` | Process-oriented skills |
| **Mapping** | `label-context-mapping`, `github-issue-creation` | Classification and routing |
| **Domain** | `bch-development`, `btc-terminology`, `wallet-cli` | Cryptocurrency-specific |

### Rules

**Location**: `.claude/rules/`

**Purpose**: Define constraints and behaviors that are always applied or conditionally applied based on file patterns.

**Characteristics**:

- Applied automatically (not explicitly loaded)
- Can be global (always applied) or scoped (applied to specific files)
- Define constraints, conventions, and guard rails

**Examples**:

| Rule | Type | Purpose |
|------|------|---------|
| `general.md` | Global | Project-wide conventions |
| `security.md` | Global | Security constraints |
| `task-context-loading.md` | Global | Context loading logic |
| `go/conventions.md` | Scoped | Go file conventions |

### Context Documents

**Location**: `docs/task-contexts/`

**Purpose**: Provide domain-specific knowledge and procedures for different task types.

**Characteristics**:

- Loaded based on task classification
- Contain checklists, procedures, and domain knowledge
- Organized by task type and chain

**Structure**:

```
docs/task-contexts/
├── bug-fix.md           # Bug fixing procedures
├── feature-add.md       # Feature implementation
├── refactoring.md       # Refactoring guidelines
├── documentation.md     # Documentation updates
├── security.md          # Security considerations
├── test.md              # Testing procedures
├── db-change.md         # Database changes
├── chain-specific.md    # Multi-chain considerations
└── chains/
    ├── btc.md           # Bitcoin-specific
    ├── bch.md           # Bitcoin Cash-specific
    ├── eth.md           # Ethereum-specific
    └── xrp.md           # XRP-specific
```

## Key Design Patterns

### 1. Label-Driven Routing

GitHub labels drive the entire workflow:

```
GitHub Issue Labels
        │
        ▼
label-context-mapping skill
        │
        ├─ Type label → Context document (what to do)
        ├─ Lang/Scope label → Development skill (how to do)
        └─ Chain label → Chain context (domain knowledge)
```

### 2. Skill Composition

Skills can load other skills to compose complex workflows:

```
fix-issue command
    │
    ├─ git-workflow skill (always)
    │
    └─ label-context-mapping skill
            │
            └─ go-development skill (from lang:go label)
                    │
                    └─ db-migration skill (if scope:db also present)
```

### 3. SSOT Chain

Information flows from authoritative sources:

```
.github/labels.yml          ← GitHub label definitions
        │
        ▼
docs/guidelines/task-classification.md  ← SSOT for label semantics
        │
        ▼
.claude/skills/label-context-mapping/  ← Operational mappings
        │
        ▼
Commands and Skills         ← Reference, don't duplicate
```

## File Organization

### SSOT Locations

| Category | SSOT Location | Description |
|----------|---------------|-------------|
| Label definitions | `docs/guidelines/task-classification.md` | All label types and meanings |
| Label → Skill mapping | `.claude/skills/label-context-mapping/` | Operational routing |
| Coding conventions | `docs/guidelines/coding-conventions.md` | Code style |
| Security rules | `docs/guidelines/security.md` | Security requirements |
| Testing standards | `docs/guidelines/testing.md` | Test requirements |

### Claude vs Cursor

| Component | Claude (SSOT) | Cursor |
|-----------|---------------|--------|
| Commands | `.claude/commands/` | `.cursor/commands/` (reference only) |
| Skills | `.claude/skills/` | `.cursor/skills/` (symlink) |
| Rules | `.claude/rules/*.md` | `.cursor/rules/*.mdc` (auto-generated) |

**Sync Process**:

- Rules: `make sync-cursor-rules`
- Skills: Symlinks (automatic)
- Commands: Cursor loads from `.claude/commands/`

## Workflow Examples

### Working on a GitHub Issue

```
1. User: /fix-issue #123

2. fix-issue command:
   - Fetch issue: gh issue view 123
   - Load label-context-mapping skill

3. label-context-mapping skill:
   - Parse labels: bug, lang:go, chain:btc
   - Determine:
     - Context: bug-fix.md
     - Skill: go-development
     - Chain: chains/btc.md

4. Load skills:
   - git-workflow (always)
   - go-development (from lang:go)

5. Execute workflow:
   - Create branch (git-workflow)
   - Load bug-fix.md context
   - Load chains/btc.md context
   - Implement fix
   - Run verification (go-development)
   - Commit and PR (git-workflow)
```

### Creating a GitHub Issue

```
1. User: Create issue for fixing BCH address validation

2. github-issue-creation skill:
   - Load label-context-mapping skill
   - Classify task:
     - Type: bug
     - Language: lang:go
     - Chain: chain:bch

3. Create proposal:
   - Title: Fix BCH address validation
   - Labels: bug, lang:go, chain:bch
   - Skills that will be used: git-workflow, go-development

4. After approval:
   - gh issue create --title "..." --label "bug,lang:go,chain:bch"
```

## Design Principles

### 1. Don't Repeat Yourself (DRY)

Information should be defined once and referenced elsewhere:

```
❌ Bad: Duplicate label→skill mapping in multiple files
✅ Good: Define in label-context-mapping, reference from others
```

### 2. Explicit Over Implicit

Relationships should be explicitly documented:

```
❌ Bad: Skills silently loaded based on file patterns
✅ Good: Commands explicitly invoke skills with clear reasoning
```

### 3. Composable Components

Skills should be small and composable:

```
❌ Bad: One mega-skill that handles everything
✅ Good: git-workflow + go-development + db-migration composed as needed
```

### 4. Progressive Disclosure

Start simple, add complexity as needed:

```
fix-issue command (simple interface)
    └─ label-context-mapping (routing logic)
        └─ Development skills (implementation details)
            └─ Context documents (domain knowledge)
```

## Related Documents

- [AGENTS.md](../../../AGENTS.md) - Entry point for all agents
- [Task Classification](../guidelines/task-classification.md) - SSOT for labels
- [label-context-mapping](../../.claude/skills/label-context-mapping/SKILL.md) - Mapping skill
- [Commands Documentation](../commands/README.md) - Command reference
- [Skills Documentation](../agent-skills.md) - Skills reference
