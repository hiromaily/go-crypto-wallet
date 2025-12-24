# Cursor Development Guide

This document describes useful features and best practices for using Cursor AI in this project.

---

## Table of Contents

- [Basic Usage](#basic-usage)
- [Custom Commands](#custom-commands)
- [Code References](#code-references)
- [Best Practices](#best-practices)
- [Keyboard Shortcuts](#keyboard-shortcuts)

---

## Basic Usage

### Opening Chat Interface

- **Mac**: `Cmd+L` or `Cmd+K`
- **Windows/Linux**: `Ctrl+L` or `Ctrl+K`

### Chat Features

1. **Ask Questions**: Ask about code, architecture, or implementation details
2. **Code Generation**: Request code implementations following project patterns
3. **Code Review**: Ask for code review and suggestions
4. **Refactoring**: Request refactoring while maintaining architecture principles

### Composer (Multi-file Editing)

- **Mac**: `Cmd+I`
- **Windows/Linux**: `Ctrl+I`

Use Composer for:

- Large refactorings across multiple files
- Implementing features that span multiple layers
- Following Clean Architecture patterns

---

## Custom Commands

This project includes custom commands in `.cursor/commands/` that provide structured workflows for common tasks.

### Create GitHub Issue

The `create-github-issue` command helps generate well-structured GitHub issues suitable for AI implementation.

#### Usage

1. Open Cursor chat (`Cmd+L` or `Ctrl+L`)
2. Type `@create-github-issue` or reference the command file
3. Describe what you want to implement or fix
4. The AI will generate a complete GitHub issue with all required sections

#### Example

```
@create-github-issue I need to add a new validator function for account types in the domain layer. 
It should validate that account types follow Clean Architecture principles and be placed in 
pkg/domain/account/validator.go. Reference the existing validators in pkg/domain/transaction/ 
for the pattern.
```

The generated issue will include:

- Clear description
- Acceptance criteria
- Technical requirements
- Implementation details
- Testing requirements
- Related context

See `.cursor/commands/create-github-issue.md` for the full template and guidelines.

---

## Code References

Cursor supports referencing existing code using a special syntax that helps AI understand context.

### Syntax

```startLine:endLine:filepath
// code content here
```

### Example

```12:14:pkg/domain/account/types.go
type AccountType string

const (
 AccountTypeStandard AccountType = "standard"
 AccountTypeMultisig AccountType = "multisig"
)
```

### When to Use Code References

- **Explaining Context**: When asking about specific code sections
- **Requesting Changes**: When modifying existing code
- **Following Patterns**: When implementing similar functionality
- **Architecture Guidance**: When discussing layer separation

---

## Best Practices

### 1. Reference Project Guidelines

Always mention relevant documentation when asking questions:

```
@AGENTS.md How should I implement a new validator in the domain layer?
```

### 2. Follow Clean Architecture

When requesting implementations, specify the layer:

- **Domain Layer**: Pure business logic, no infrastructure dependencies
- **Application Layer**: Business logic orchestration (`pkg/wallet/service/`)
- **Infrastructure Layer**: External dependencies (`pkg/infrastructure/`)

### 3. Use Code References

Instead of describing code locations, use code references:

```
Please refactor this function following the pattern in:
```12:25:pkg/domain/transaction/validator.go
// existing validator code
```

```

### 4. Specify Constraints

Always mention project-specific constraints:

- Security requirements (no private key logging)
- Architecture patterns (Clean Architecture, dependency injection)
- Error handling (use `fmt.Errorf` + `%w`)
- Testing requirements (unit tests, integration tests)

### 5. Request Incremental Changes

For large refactorings, break them into smaller steps:

```

Step 1: Create the domain validator interface
Step 2: Implement the validator in domain layer
Step 3: Update the application layer to use it

```

---

## Keyboard Shortcuts

### General

| Action | Mac | Windows/Linux |
|--------|-----|---------------|
| Open Chat | `Cmd+L` | `Ctrl+L` |
| Open Composer | `Cmd+I` | `Ctrl+I` |
| Accept Suggestion | `Tab` | `Tab` |
| Reject Suggestion | `Esc` | `Esc` |
| Show All Shortcuts | `Cmd+Shift+P` | `Ctrl+Shift+P` |

### Code Actions

| Action | Mac | Windows/Linux |
|--------|-----|---------------|
| Explain Code | `Cmd+K Cmd+E` | `Ctrl+K Ctrl+E` |
| Generate Code | `Cmd+K` | `Ctrl+K` |
| Edit Selection | `Cmd+K Cmd+I` | `Ctrl+K Ctrl+I` |
