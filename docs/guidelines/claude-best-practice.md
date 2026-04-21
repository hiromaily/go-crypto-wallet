<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/docs/guidelines/claude-best-practice.tpl.md · Run `make docs` to regenerate.
-->

# Claude Code Best Practices Applied to go-crypto-wallet

> Based on [Claude Code Best Practices](https://code.claude.com/docs/en/best-practices)
> Applied to the go-crypto-wallet project structure and conventions.

## Table of Contents

- [Overview](#overview)
- [1. Give Claude a Way to Verify Its Work](#1-give-claude-a-way-to-verify-its-work)
- [2. Explore First, Then Plan, Then Code](#2-explore-first-then-plan-then-code)
- [3. Provide Specific Context in Prompts](#3-provide-specific-context-in-prompts)
- [4. Configure Your Environment](#4-configure-your-environment)
- [5. Communicate Effectively](#5-communicate-effectively)
- [6. Manage Your Session](#6-manage-your-session)
- [7. Automate and Scale](#7-automate-and-scale)
- [8. Avoid Common Failure Patterns](#8-avoid-common-failure-patterns)
- [Current State Assessment](#current-state-assessment)
- [Gap Analysis and Recommendations](#gap-analysis-and-recommendations)

---

## Overview

The core principle of Claude Code best practices is **context window management**. Performance degrades as context fills. Every decision below maps back to this constraint.

This document evaluates the go-crypto-wallet project against each best practice, identifies what is already implemented, and recommends improvements.

---

## 1. Give Claude a Way to Verify Its Work

> The single highest-leverage thing you can do.

### What the Best Practice Says

- Provide tests, screenshots, or expected outputs so Claude can check itself
- Address root causes, not symptoms
- Invest in making verification rock-solid

### Current Implementation

| Mechanism | Status | Details |
|-----------|--------|---------|
| Makefile verification targets | Implemented | `make go-lint`, `make check-build`, `make go-test` |
| File-type verification matrix | Implemented | `.claude/rules/task-context-loading.md` maps file extensions to verification commands |
| go-development skill | Implemented | Enforces `make go-lint && make tidy && make check-build && make go-test` |
| Pre-commit hooks (lefthook) | Implemented | Conventional commit enforcement |
| Security vulnerability scan | Implemented | `make go-check-vuln` |
| docs/ai/task-contexts/verification.md | Implemented | Central verification command reference |

### Recommendations

| Priority | Recommendation | Rationale |
|----------|---------------|-----------|
| High | Add a hook to auto-run `make go-lint` after Go file edits | Hooks guarantee execution (unlike advisory CLAUDE.md rules). Best practice: "Use hooks for actions that must happen every time with zero exceptions." |
| Medium | Add integration test targets per chain | Currently `make go-test` runs unit tests only; chain-specific integration tests would strengthen verification |
| Low | Consider screenshot verification for any future web UI components | Not currently applicable (CLI-only project) |

---

## 2. Explore First, Then Plan, Then Code

> Separate research and planning from implementation to avoid solving the wrong problem.

### What the Best Practice Says

Four-phase workflow: Explore -> Plan -> Implement -> Commit

### Current Implementation

| Mechanism | Status | Details |
|-----------|--------|---------|
| Kiro spec-driven workflow | Implemented | Requirements -> Design -> Tasks -> Implementation via `/kiro:spec-*` commands |
| Gap analysis | Implemented | `/kiro:validate-gap` before design |
| Design validation | Implemented | `/kiro:validate-design` after design |
| Implementation validation | Implemented | `/kiro:validate-impl` after implementation |
| OpenSpec exploration | Implemented | `/opsx:explore` for idea exploration before proposing changes |
| Plan Mode support | Available | Built into Claude Code (Ctrl+G) |

### Assessment

This project **exceeds** the best practice with the Kiro spec-driven development workflow. The 3-phase approval process (Requirements -> Design -> Tasks) is more rigorous than the simple Plan Mode approach described in the best practice.

### Recommendations

| Priority | Recommendation | Rationale |
|----------|---------------|-----------|
| Medium | Use Plan Mode for small tasks that don't warrant a full Kiro spec | Kiro specs are heavyweight; simple bug fixes or small refactors benefit from lightweight Plan Mode |
| Low | Document when to use Kiro specs vs Plan Mode vs direct implementation | Help users choose the right approach based on task complexity |

---

## 3. Provide Specific Context in Prompts

> The more precise your instructions, the fewer corrections you'll need.

### What the Best Practice Says

- Reference specific files
- Point to existing patterns
- Describe symptoms with location and success criteria
- Use rich content: `@` references, images, URLs, piped data

### Current Implementation

| Mechanism | Status | Details |
|-----------|--------|---------|
| Task context loading | Implemented | `.claude/rules/task-context-loading.md` auto-loads relevant context docs based on task type |
| Label-context mapping | Implemented | `.claude/skills/label-context-mapping/` routes GitHub labels to skills and context files |
| Chain-specific contexts | Implemented | `docs/ai/task-contexts/chain-specific.md` + `docs/chains/{chain}/README.md` |
| Task-type contexts | Implemented | `docs/ai/task-contexts/{bug-fix,feature-add,refactoring,...}.md` |
| Architecture docs | Implemented | `ARCHITECTURE.md`, `docs/guidelines/architecture.md` |
| `@` file references in CLAUDE.md | Implemented | CLAUDE.md references `@AGENTS.md`, `@docs/guidelines/`, `@ARCHITECTURE.md` |

### Assessment

This project has a **sophisticated context loading system** that goes beyond the best practice. The label-driven routing pattern automatically provides Claude with the right context based on GitHub issue labels.

### Recommendations

| Priority | Recommendation | Rationale |
|----------|---------------|-----------|
| Medium | Add example prompts to task context docs | Show users how to write effective prompts for each task type (e.g., "fix BCH address validation bug in watch wallet") |
| Low | Add `@` references in frequently-used commands | Some slash commands could benefit from explicit file references to reduce Claude's exploration |

---

## 4. Configure Your Environment

### 4.1 CLAUDE.md

> Run `/init` to generate a starter CLAUDE.md, then refine over time.

#### What the Best Practice Says

- Keep it short and human-readable
- Include only things Claude can't infer from code
- Prune regularly; bloated CLAUDE.md causes instruction loss
- Check it into git for team sharing

#### Current Implementation

| Criterion | Status | Details |
|-----------|--------|---------|
| Checked into git | Yes | `CLAUDE.md` is tracked |
| Concise | Yes | Focused on Kiro workflow, paths, and development rules |
| Bash commands Claude can't guess | Partial | Via `@AGENTS.md` and `@docs/guidelines/` references |
| Code style that differs from defaults | Delegated | To `.claude/rules/` and `docs/guidelines/coding-conventions.md` |
| Architectural decisions | Delegated | To `ARCHITECTURE.md` |

#### Assessment

The CLAUDE.md follows a **layered approach**: it is concise itself and delegates detailed rules to `.claude/rules/` (auto-loaded) and `docs/guidelines/` (loaded on demand via `@` references). This aligns well with the best practice of keeping CLAUDE.md short.

#### Recommendations

| Priority | Recommendation | Rationale |
|----------|---------------|-----------|
| Medium | Audit CLAUDE.md + rules for redundancy | Best practice: "If Claude already does something correctly without the instruction, delete it." Some rules may overlap between CLAUDE.md, AGENTS.md, and `.claude/rules/` |
| Low | Add compaction preservation instruction | e.g., "When compacting, always preserve the full list of modified files and any test commands" |

### 4.2 Permissions

> Use `/permissions` to allowlist safe commands.

#### Current Implementation

| Mechanism | Status | Details |
|-----------|--------|---------|
| `.claude/settings.local.json` allow list | Implemented | Whitelists `make *`, `go *`, `git *`, `gh *`, `docker`, `sqlc`, shell tools, `bun`, etc. |
| Granular git permissions | Implemented | Specific git subcommands allowed (add, commit, push, branch, etc.) |

#### Assessment

Well-configured. The allow list covers standard development commands while maintaining security boundaries.

### 4.3 CLI Tools

> Tell Claude Code to use CLI tools like `gh`, `aws`, `gcloud`.

#### Current Implementation

| Tool | Status | Usage |
|------|--------|-------|
| `gh` (GitHub CLI) | Configured | Issue management, PR creation, release management |
| `docker` | Configured | Database containers, build verification |
| `atlas` | Configured | Database schema management |
| `sqlc` | Configured | SQL code generation |
| `mockery` | Configured | Mock generation |
| `golangci-lint` | Configured | Via `make go-lint` |
| `govulncheck` | Configured | Via `make go-check-vuln` |

### 4.4 MCP Servers

> Run `claude mcp add` to connect external tools.

#### Current Implementation

| MCP | Status | Details |
|-----|--------|---------|
| claude-mem plugin | Installed | Persistent cross-session memory via `@thedotmack/claude-mem` |

#### Recommendations

| Priority | Recommendation | Rationale |
|----------|---------------|-----------|
| Low | Consider adding a code intelligence MCP/plugin | Best practice recommends code intelligence plugins for typed languages; could provide precise symbol navigation for Go |

### 4.5 Hooks

> Use hooks for actions that must happen every time with zero exceptions.

#### Current Implementation

| Hook | Status | Details |
|------|--------|---------|
| Stop notification sound | Implemented | Plays `Funk.aiff` when Claude stops |
| Pre-commit (lefthook) | Implemented | Conventional commit enforcement |

#### Recommendations

| Priority | Recommendation | Rationale |
|----------|---------------|-----------|
| High | Add post-edit hook for Go files: `make go-lint` | Guarantees linting after every Go file edit, not just advisory |
| Medium | Add hook to block writes to auto-generated files | Files with `DO NOT EDIT` headers should be protected by hooks, not just rules |
| Medium | Add hook to block writes to `main` branch | Currently advisory in AGENTS.md; a hook would enforce it deterministically |

### 4.6 Skills

> Create `SKILL.md` files in `.claude/skills/` for domain knowledge and reusable workflows.

#### Current Implementation

The project has **extensive skills** covering:

- Language-specific development (Go, TypeScript, Solidity, Shell)
- Task-type workflows (git-workflow, fix-pr-review, fix-issue)
- Chain-specific knowledge (bch-development, btc-terminology)
- Infrastructure (db-migration, devops, makefile-update)
- AI development lifecycle (openspec-*, knowledge-gap)

#### Assessment

**Exceeds** the best practice. The skill system with label-context-mapping creates an automated routing layer that the best practice doesn't cover.

### 4.7 Custom Subagents

> Define specialized assistants in `.claude/agents/` for isolated tasks.

#### Current Implementation

No custom subagents defined in `.claude/agents/`.

#### Recommendations

| Priority | Recommendation | Rationale |
|----------|---------------|-----------|
| Medium | Create a `security-reviewer` subagent | This project handles private keys and crypto transactions; a dedicated security review agent would catch vulnerabilities in a separate context |
| Medium | Create a `architecture-reviewer` subagent | Verify Clean Architecture layer separation in a fresh context window |
| Low | Create a `chain-compatibility` subagent | Review changes for cross-chain compatibility issues |

---

## 5. Communicate Effectively

### 5.1 Codebase Questions

> Ask Claude questions you'd ask a senior engineer.

No project-specific configuration needed. Claude can already explore the codebase.

### 5.2 Let Claude Interview You

> For larger features, have Claude interview you first.

#### Current Implementation

The Kiro spec workflow (`/kiro:spec-init` -> `/kiro:spec-requirements`) already captures requirements through a structured process. For features outside the Kiro workflow, the interview pattern is available but not documented.

#### Recommendations

| Priority | Recommendation | Rationale |
|----------|---------------|-----------|
| Low | Add an interview-style command for quick features | For features too small for Kiro specs but complex enough to need requirements gathering |

---

## 6. Manage Your Session

### What the Best Practice Says

- Course-correct early with `Esc`, `/rewind`, `/clear`
- Manage context aggressively: `/clear` between tasks, `/compact` with focus instructions
- Use subagents for investigation to keep main context clean
- Checkpoints for safe experimentation
- Resume conversations with `--continue` / `--resume`

### Current Implementation

| Mechanism | Status | Details |
|-----------|--------|---------|
| claude-mem plugin | Implemented | Persistent memory across sessions |
| Session context hooks | Implemented | Startup hook loads recent context |
| Kiro spec artifacts | Implemented | Specs persist across sessions as files |

### Recommendations

| Priority | Recommendation | Rationale |
|----------|---------------|-----------|
| Medium | Add compaction instructions to CLAUDE.md | e.g., "When compacting, preserve: modified file list, current spec phase, chain context, and verification results" |
| Medium | Document subagent usage patterns in guidelines | Best practice: "Subagents are one of the most powerful tools available" - document when to use them for this project |
| Low | Add `/rename` convention for spec-related sessions | e.g., `"spec-postgresql-integration"`, `"fix-bch-e2e-p1"` for easy resumption |

---

## 7. Automate and Scale

### What the Best Practice Says

- `claude -p "prompt"` for headless mode in CI/scripts
- Multiple parallel sessions (Writer/Reviewer pattern)
- Fan out across files for large migrations
- `--dangerously-skip-permissions` only in sandboxed containers

### Current Implementation

| Mechanism | Status | Details |
|-----------|--------|---------|
| Headless mode | Not used | CI uses standard GitHub Actions, not Claude |
| Parallel sessions | Not configured | Single-session workflow |
| Fan-out scripts | Not used | No batch migration tooling |

### Recommendations

| Priority | Recommendation | Rationale |
|----------|---------------|-----------|
| Medium | Use Writer/Reviewer pattern for security-critical changes | Have one session implement, another review for security issues (critical for a crypto wallet) |
| Low | Consider headless Claude for automated code review in CI | `claude -p "Review this PR for security issues" --output-format json` |
| Low | Fan-out pattern for future large migrations | Useful when migrating across chains or refactoring infrastructure layer |

---

## 8. Avoid Common Failure Patterns

### What the Best Practice Says

| Anti-Pattern | Fix |
|-------------|-----|
| Kitchen sink session | `/clear` between unrelated tasks |
| Correcting over and over | After two failed corrections, `/clear` and write a better prompt |
| Over-specified CLAUDE.md | Prune ruthlessly; use hooks instead of advisory rules |
| Trust-then-verify gap | Always provide verification |
| Infinite exploration | Scope investigations or use subagents |

### Project-Specific Risk Assessment

| Anti-Pattern | Risk Level | Mitigation in Place |
|-------------|-----------|-------------------|
| Kitchen sink session | Medium | No explicit mitigation; relies on user discipline |
| Correcting over and over | Medium | Kiro spec workflow reduces this by front-loading requirements |
| Over-specified CLAUDE.md | Low | Layered architecture delegates detail to rules/skills |
| Trust-then-verify gap | Low | Makefile verification matrix and go-development skill enforce checks |
| Infinite exploration | Medium | Task-context-loading scopes context; but no subagent enforcement |

### Recommendations

| Priority | Recommendation | Rationale |
|----------|---------------|-----------|
| Medium | Add guideline: "Use subagents for codebase exploration" | Prevent context pollution during investigation |
| Low | Add `/clear` reminder in long-running skills | Some E2E fix commands may run long; periodic context reset improves quality |

---

## Current State Assessment

### Scorecard

| Best Practice | Score | Notes |
|---------------|-------|-------|
| Verification | 4/5 | Strong Makefile targets; missing post-edit hooks |
| Planning workflow | 5/5 | Kiro spec-driven development exceeds the recommendation |
| Context in prompts | 5/5 | Automated label-driven context loading |
| CLAUDE.md | 4/5 | Well-structured layered approach; minor redundancy |
| Permissions | 4/5 | Comprehensive allow list |
| CLI tools | 5/5 | Full coverage of needed tools |
| MCP servers | 3/5 | claude-mem only; could add code intelligence |
| Hooks | 2/5 | Only notification sound; missing enforcement hooks |
| Skills | 5/5 | Extensive, well-organized skill system |
| Subagents | 1/5 | No custom subagents defined |
| Session management | 3/5 | claude-mem helps; missing compaction and subagent guidelines |
| Automation | 2/5 | No headless/parallel/fan-out usage |

### Overall: 43/60 (72%)

---

## Gap Analysis and Recommendations

### High Priority

1. **Add enforcement hooks**
   - Post-edit hook: Run `make go-lint` after Go file modifications
   - Block hook: Prevent writes to auto-generated files (`DO NOT EDIT`)
   - Block hook: Prevent commits on `main` branch
   - Reference: [Hooks Guide](https://code.claude.com/docs/en/hooks-guide)

### Medium Priority

1. **Create custom subagents** (`.claude/agents/`)
   - `security-reviewer.md`: Review code for private key leaks, injection, OWASP issues
   - `architecture-reviewer.md`: Verify Clean Architecture layer separation
   - Reference: [Subagents Guide](https://code.claude.com/docs/en/sub-agents)

2. **Add compaction instructions to CLAUDE.md**

   ```markdown
   # Context Management
   When compacting, preserve: modified file list, current spec phase,
   chain context, and all verification results.
   ```

3. **Document subagent usage patterns**
   - Add to `docs/guidelines/workflow.md`: when to use subagents for exploration vs. main context

4. **Adopt Writer/Reviewer pattern for security-critical changes**
   - One session implements; another reviews with fresh context
   - Essential for crypto wallet security

5. **Audit rules for redundancy**
   - Some rules appear in both CLAUDE.md -> AGENTS.md and `.claude/rules/`
   - Consolidate to reduce context consumption

### Low Priority

1. **Explore code intelligence plugin** for Go symbol navigation
2. **Consider headless Claude in CI** for automated security review
3. **Add session naming conventions** for spec-related work
4. **Document Plan Mode vs Kiro specs decision criteria**

---

## Related Documents

- `CLAUDE.md` - Project root instructions
- `AGENTS.md` - AI agent guidelines
- [AI Agents Instruction Architecture](../ai/design.md) - System design
- [Workflow Guidelines](./workflow.md) - Development workflow
- [Security Guidelines](./security.md) - Security requirements
- [Verification Matrix](../ai/task-contexts/verification.md) - File-type verification commands
