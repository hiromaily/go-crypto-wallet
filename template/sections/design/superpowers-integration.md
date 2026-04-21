# Superpowers Integration Design

This document describes the plan to integrate [obra/superpowers](https://github.com/obra/superpowers) into this project's AI agent instruction system.

## Overview

Superpowers is an agentic skills framework that provides composable workflow skills for AI coding agents. It enforces systematic development practices including design-first brainstorming, test-driven development, structured planning, and subagent-driven execution.

### Why Integrate?

This project already has a mature `.claude/` instruction system (rules, skills, commands). Superpowers complements it by adding **process-level skills** that govern _how_ agents work, while the existing system defines _what_ agents know about this codebase.

| Layer | Existing System | Superpowers |
|-------|----------------|-------------|
| Domain knowledge | `.claude/rules/`, `.claude/skills/` | - |
| Task workflows | `.claude/commands/` | - |
| Development process | Partially (git-workflow, go-development) | Fully (TDD, planning, debugging, code review) |
| Spec-driven dev | `.kiro/specs/` | brainstorming, writing-plans |

## Superpowers Skills Inventory

The 14 skills and their relevance to this project:

| Skill | Purpose | Relevance | Conflict Risk |
|-------|---------|-----------|---------------|
| `brainstorming` | Socratic design refinement | High - complements `.kiro/` spec workflow | Low |
| `writing-plans` | Detailed task breakdown (2-5 min tasks) | High - enhances implementation planning | Medium - overlaps with `.kiro/specs/tasks.md` |
| `executing-plans` | Batch execution with checkpoints | High - structured task execution | Low |
| `test-driven-development` | RED-GREEN-REFACTOR enforcement | High - this project needs more TDD discipline | Low |
| `systematic-debugging` | 4-phase root cause analysis | High - valuable for crypto wallet debugging | Low |
| `verification-before-completion` | Validates fixes before marking done | High - aligns with existing verification matrix | Low |
| `subagent-driven-development` | Parallel subagent execution with 2-stage review | Medium - useful for large features | Low |
| `dispatching-parallel-agents` | Concurrent subagent workflows | Medium - useful for multi-file changes | Low |
| `requesting-code-review` | Pre-review checklist | Medium - useful for PR preparation | Low |
| `receiving-code-review` | Feedback response workflow | Medium - useful for PR review responses | Low |
| `using-git-worktrees` | Parallel development branches | Low - project uses standard branching | Medium - may conflict with git-workflow skill |
| `finishing-a-development-branch` | Merge/PR decision workflow | Low - project has git-workflow skill | Medium - overlaps with git-workflow |
| `writing-skills` | Meta: create new skills | Low - rarely needed | Low |
| `using-superpowers` | System introduction | Low - onboarding only | Low |

## Integration Strategy

### Option A: Plugin Install (Recommended)

Use the Claude Code plugin marketplace for a clean, updatable installation.

```bash
# One-time marketplace registration
/plugin marketplace add obra/superpowers-marketplace

# Install superpowers
/plugin install superpowers@superpowers-marketplace
```

**Pros:**

- Clean separation: superpowers skills live in plugin space, not in repo
- Easy updates via `/plugin update superpowers`
- No file conflicts with existing `.claude/` structure
- Zero maintenance burden on this project

**Cons:**

- All 14 skills are installed (no cherry-picking)
- Plugin skills may trigger when existing project skills are more appropriate
- Requires each developer to install individually

### Option B: Selective Vendoring

Copy only high-value skills into `.claude/skills/` with project-specific adaptations.

```
.claude/skills/
├── superpowers/                    # Vendored superpowers skills
│   ├── brainstorming/SKILL.md
│   ├── test-driven-development/SKILL.md
│   ├── systematic-debugging/SKILL.md
│   ├── writing-plans/SKILL.md
│   ├── executing-plans/SKILL.md
│   ├── verification-before-completion/SKILL.md
│   └── subagent-driven-development/SKILL.md
├── go-development/SKILL.md         # Existing (unchanged)
├── git-workflow/SKILL.md            # Existing (unchanged)
└── ...
```

**Pros:**

- Cherry-pick only relevant skills (7 of 14)
- Can adapt skills to project conventions (e.g., reference existing verification matrix)
- Committed to repo - works for all contributors automatically
- No external dependency at runtime

**Cons:**

- Manual updates when upstream changes
- Potential drift from upstream
- More files in the repo

### Option C: Plugin Install + Override Rules (Recommended)

Install via plugin marketplace but add project-specific rules to control behavior.

```
.claude/rules/superpowers.md        # Override/integration rules
```

This rule file would:

1. Define when superpowers skills should defer to existing project skills
2. Add project-specific constraints (e.g., security rules during TDD)
3. Map superpowers workflows to existing verification commands

**Pros:**

- Clean plugin installation with easy updates
- Project-specific behavior via override rules
- Minimal repo changes (one rule file)
- Best of both worlds

**Cons:**

- Requires each developer to install the plugin
- Override rules need maintenance when superpowers updates

## Recommended Approach: Option C

### Rationale

1. **Separation of concerns**: Superpowers handles process, project rules handle domain
2. **Updatability**: Plugin updates don't require repo changes
3. **Control**: Override rules prevent conflicts with existing workflows
4. **Minimal impact**: Only one new file added to the repo

### Implementation Plan

#### Phase 1: Install and Evaluate

1. Register the superpowers marketplace
2. Install the plugin
3. Test with a small feature task to observe skill interactions
4. Identify conflicts with existing skills

#### Phase 2: Create Override Rules

Create `.claude/rules/superpowers.md` with:

```markdown
# Superpowers Integration Rules

## Skill Priority

When superpowers skills conflict with project skills, project skills take precedence:

- **git-workflow** overrides `using-git-worktrees` and `finishing-a-development-branch`
- **go-development** verification commands override superpowers verification steps
- `.kiro/specs/` workflow overrides `writing-plans` for spec-tracked features

## Security Constraints

During TDD and debugging workflows:
- Never log private keys or sensitive data in test fixtures
- Test wallets must use testnet/regtest only
- Follow docs/guidelines/security.md

## Verification Integration

Superpowers `verification-before-completion` must use this project's verification matrix:
- Go files: `make go-lint`, `make check-build`
- SQL/HCL files: `make atlas-fmt`, `make atlas-lint`
- Shell scripts: `make shfmt`
- See docs/ai/task-contexts/verification.md for full matrix

## Planning Integration

When `writing-plans` is triggered for a `.kiro/specs/` tracked feature:
- Reference the existing spec's requirements.md and design.md
- Plans should align with the spec's tasks.md structure
```

#### Phase 3: Update Documentation

1. Update `AGENTS.md` to mention superpowers integration
2. Add setup instructions to contributor docs
3. Update [ai-agents-instruction.md](../../../docs/ai/design.md) architecture diagram

#### Phase 4: Team Onboarding

1. Add superpowers install commands to project setup guide
2. Document which skills are most useful for this project
3. Add a Makefile target for convenience:

```makefile
# make/ai.mk
.PHONY: setup-superpowers
setup-superpowers:
 @echo "Run these commands in Claude Code:"
 @echo "  /plugin marketplace add obra/superpowers-marketplace"
 @echo "  /plugin install superpowers@superpowers-marketplace"
```

## Conflict Resolution Matrix

| Scenario | Superpowers Skill | Project Skill | Resolution |
|----------|-------------------|---------------|------------|
| Creating a feature branch | `using-git-worktrees` | `git-workflow` | Use `git-workflow` |
| Planning implementation | `writing-plans` | `.kiro/spec-tasks` | Use `.kiro/` for spec-tracked features, `writing-plans` for ad-hoc tasks |
| Running tests | `test-driven-development` | `go-development` | Use both: TDD process from superpowers, commands from `go-development` |
| Debugging | `systematic-debugging` | (none) | Use superpowers |
| Code review | `requesting-code-review` | `review-code` command | Use both: checklist from superpowers, domain checks from `review-code` |
| Finishing PR | `finishing-a-development-branch` | `git-workflow` | Use `git-workflow` |
| Verification | `verification-before-completion` | `task-context-loading` verification matrix | Merge: use project's verification commands within superpowers' verification flow |

## Architecture After Integration

```
┌─────────────────────────────────────────────────────────────────────┐
│                         User Request                                │
└─────────────────────────────────┬───────────────────────────────────┘
                                  │
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│                  Superpowers (Process Skills)                       │
│                   (Plugin - auto-triggered)                         │
│                                                                     │
│  brainstorming → writing-plans → TDD → verification                │
│  systematic-debugging, subagent-driven-development                  │
└─────────────────────────────────┬───────────────────────────────────┘
                                  │ delegates to
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│              Project Skills (Domain Knowledge)                      │
│                    (.claude/skills/)                                │
│                                                                     │
│  go-development, git-workflow, db-migration, etc.                   │
│  Override rules in .claude/rules/superpowers.md                     │
└─────────────────────────────────┬───────────────────────────────────┘
                                  │ references
                                  ▼
┌─────────────────────────────────────────────────────────────────────┐
│              Project Rules & Context Documents                      │
│              (.claude/rules/, docs/task-contexts/)                  │
│                                                                     │
│  security.md, task-context-loading.md, verification.md              │
└─────────────────────────────────────────────────────────────────────┘
```

## Risk Assessment

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| Superpowers overrides project conventions | Medium | Medium | Override rules in `.claude/rules/superpowers.md` |
| TDD skill conflicts with existing test patterns | Low | Low | Go-development skill takes precedence for commands |
| Git worktree skill creates unexpected branches | Medium | Low | Override rule: use `git-workflow` instead |
| Plugin updates break integration | Low | Low | Pin version, test before updating |
| Increased token usage from extra skills | Low | Medium | Acceptable tradeoff for better process |

## Success Criteria

1. Superpowers brainstorming and TDD skills activate for new features
2. Existing project skills (go-development, git-workflow) remain authoritative for domain-specific tasks
3. No security regressions - superpowers respects security rules
4. Developer onboarding takes < 5 minutes
5. Override rules prevent workflow conflicts
