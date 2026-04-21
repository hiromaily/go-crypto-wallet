<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/docs/guidelines/overview.tpl.md · Run `make docs` to regenerate.
-->

# Project Guidelines

This directory contains project guidelines, standards, and the Single Source of Truth (SSOT) documents.
All AI agent configurations and developer documentation should reference these guidelines.

## Quick Reference

| Document | Description |
|----------|-------------|
| [Coding Conventions](./coding-conventions.md) | Code style, naming, formatting, language-specific rules |
| [Security](./security.md) | Security requirements and practices |
| [Testing](./testing.md) | Testing strategy, layers, mocking, and requirements |
| [Workflow](./workflow.md) | Git operations, branch naming, commit messages, verification |
| [Task Classification](./task-classification.md) | Labels, task types, and skill mappings (SSOT) |
| [Architecture](./architecture.md) | Clean Architecture principles and layer structure |
| [Core Principles](./core.md) | Error handling, panic usage, logging, patterns |
| [Code Generation](./code-generation.md) | Atlas, SQLC, Mockery, and other code generation tools |
| [Multi-Chain](./multi-chain.md) | Multi-chain cryptocurrency support and architecture |
| [Release](./release.md) | Release process (GoReleaser, versioning, workflow) |
| [Requirements](./requirements.md) | Tool versions and installation instructions |
| [Claude-Mem](./claude-mem.md) | Persistent memory plugin for Claude Code sessions |

## How to Use

### For AI Agents

Reference these guidelines in agent configuration files:

```markdown
# Example: .claude/rules/coding.md
See docs/guidelines/coding-conventions.md for coding standards.
```

### For Developers

These guidelines apply to all contributors. Read the relevant document before:

- Writing new code -> [Coding Conventions](./coding-conventions.md)
- Handling sensitive data -> [Security](./security.md)
- Writing tests -> [Testing](./testing.md)
- Creating PRs -> [Workflow](./workflow.md)
