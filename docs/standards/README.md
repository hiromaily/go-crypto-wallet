# Project Standards (SSOT)

This directory contains the **Single Source of Truth (SSOT)** for project standards.
All AI agent configurations and developer documentation should reference these standards.

## Quick Reference

| Standard | Description |
|----------|-------------|
| [Coding Conventions](coding-conventions.md) | Code style, naming, formatting |
| [Security](security.md) | Security requirements and practices |
| [Testing](testing.md) | Testing strategy and requirements |
| [Workflow](workflow.md) | Git operations and development workflow |
| [Task Classification](task-classification.md) | Labels, task types, and skill mappings (SSOT) |

## Detailed Guidelines

For comprehensive documentation (including AI-specific context), see:

- [Detailed AI Agent Guidelines](../ai-agents/guidelines/) - Full guidelines with examples

### Database Management

- [Database Schema Changes](../guidelines/database-schema-changes.md) - **Complete workflow** for schema modifications across multiple databases (MySQL, SQLite, PostgreSQL)
- [Database Quick Reference](../guidelines/database-quick-reference.md) - **Cheat sheet** for common database operations and commands
- [Database Management](../guidelines/database.md) - Overview of database architecture and tools
- [Database Architecture](../development/database.md) - Detailed database setup and operations

## How to Use

### For AI Agents

Reference these standards in agent configuration files:

```markdown
# Example: .cursor/rules/coding.mdc
See docs/standards/coding-conventions.md for coding standards.
```

### For Developers

These standards apply to all contributors. Read the relevant document before:

- Writing new code → [Coding Conventions](coding-conventions.md)
- Handling sensitive data → [Security](security.md)
- Writing tests → [Testing](testing.md)
- Creating PRs → [Workflow](workflow.md)
