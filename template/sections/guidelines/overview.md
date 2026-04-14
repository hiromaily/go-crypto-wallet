## Project Guidelines

This directory contains project guidelines, standards, and the Single Source of Truth (SSOT) documents.
All AI agent configurations and developer documentation should reference these guidelines.

### Quick Reference

| Document | Description |
|----------|-------------|
| [Coding Conventions](../../../docs/guidelines/coding-conventions.md) | Code style, naming, formatting, language-specific rules |
| [Security](../../../docs/guidelines/security.md) | Security requirements and practices |
| [Testing](../../../docs/guidelines/testing.md) | Testing strategy, layers, mocking, and requirements |
| [Workflow](../../../docs/guidelines/workflow.md) | Git operations, branch naming, commit messages, verification |
| [Task Classification](../../../docs/guidelines/task-classification.md) | Labels, task types, and skill mappings (SSOT) |
| [Architecture](../../../docs/guidelines/architecture.md) | Clean Architecture principles and layer structure |
| [Core Principles](../../../docs/guidelines/core.md) | Error handling, panic usage, logging, patterns |
| [Code Generation](../../../docs/guidelines/code-generation.md) | Atlas, SQLC, Mockery, and other code generation tools |
| [Multi-Chain](../../../docs/guidelines/multi-chain.md) | Multi-chain cryptocurrency support and architecture |
| [Release](../../../docs/guidelines/release.md) | Release process (GoReleaser, versioning, workflow) |
| [Requirements](../../../docs/guidelines/requirements.md) | Tool versions and installation instructions |
| [Claude-Mem](../../../docs/guidelines/claude-mem.md) | Persistent memory plugin for Claude Code sessions |

### How to Use

#### For AI Agents

Reference these guidelines in agent configuration files:

```markdown
# Example: .claude/rules/coding.md
See docs/guidelines/coding-conventions.md for coding standards.
```

#### For Developers

These guidelines apply to all contributors. Read the relevant document before:

- Writing new code -> [Coding Conventions](../../../docs/guidelines/coding-conventions.md)
- Handling sensitive data -> [Security](../../../docs/guidelines/security.md)
- Writing tests -> [Testing](../../../docs/guidelines/testing.md)
- Creating PRs -> [Workflow](../../../docs/guidelines/workflow.md)
