# Claude Development Guide

This document describes useful features and best practices for using Claude AI in this project.

## Table of Contents

- [Basic Usage](#basic-usage)
- [Custom Commands](#custom-commands)
- [Best Practices](#best-practices)
- [Workflow Examples](#workflow-examples)

## Basic Usage

### Claude Chat Interface

Claude can be accessed through:

- **Claude Desktop App**: Official desktop application
- **Claude Web**: Browser-based interface at claude.ai
- **VS Code Extension**: If available

### Key Features

1. **Code Analysis**: Understand code structure, architecture, and patterns
2. **Implementation**: Generate code following project guidelines
3. **Refactoring**: Refactor code while maintaining architecture principles
4. **Documentation**: Generate and update documentation
5. **Issue Resolution**: Fix bugs and implement features from GitHub issues

## Custom Commands

This project includes custom commands in `.claude/commands/` that provide structured workflows for common development tasks.

### Fix Issue

The `fix-issue` command provides a systematic workflow for resolving GitHub issues.

#### Usage

Reference the command file when asking Claude to fix an issue:

```
@.claude/commands/fix-issue.md Fix issue #123
```

Or describe the issue number directly:

```
Fix issue #123 following the fix-issue workflow
```

#### Workflow Steps

1. **Pre-Flight Checks**
   - Verify clean working directory
   - Check current branch (never work on main/master)
   - Fetch issue details using `gh issue view`

2. **Create Feature Branch**
   - Format: `feature/issue-{number}-{description}`
   - Example: `feature/issue-123-fix-logger-global-issue`

3. **Resolve Systematically**
   - Analyze the issue
   - Plan the solution
   - Implement following Clean Architecture
   - Test thoroughly
   - Document changes
   - Verify with `make lint-fix`, `make check-build`, `make gotest`
   - Commit and create PR

#### Safety Rules

- **CRITICAL**: Never work on `main`/`master` branch
- **CRITICAL**: Never edit auto-generated files (sqlc, protoc, etc.)
- **CRITICAL**: Never log private keys or sensitive information
- Always verify branch and git status before starting

See `.claude/commands/fix-issue.md` for the complete workflow.

### Fix PR Review Comments

The `fix-pr-review` command helps address review comments on pull requests.

#### Usage

```
@.claude/commands/fix-pr-review.md Fix review comments for PR #456
```

#### Process

1. **Pre-Flight Checks**
   - Verify clean working directory
   - Checkout the PR branch
   - Fetch PR information and review comments

2. **Resolve Comments Systematically**
   - Categorize comments (security, functionality, code quality, etc.)
   - Prioritize fixes (security first)
   - Implement fixes
   - Test changes
   - Verify with lint and build commands
   - Commit and push to PR branch

See `.claude/commands/fix-pr-review.md` for details.

### Fix Linter

The `fix-linter` command helps fix linting errors systematically.

#### Usage

```
@.claude/commands/fix-linter.md Fix the linting errors from make lint-fix
```

#### Guidelines

1. Analyze error types and severity
2. Prioritize critical errors first
3. Fix errors incrementally
4. Batch similar errors together
5. Preserve functionality

#### Process

- Show error summary
- Propose fix order
- Fix critical errors first
- Verify fixes don't introduce new issues

See `.claude/commands/fix-linter.md` for details.

### Commit, Push, and Create PR

The `commit-push-pr` command automates the process of committing changes, pushing to remote, and creating a pull request.

#### Usage

```
@.claude/commands/commit-push-pr.md Create a commit, push changes, and create a PR
```

#### Process

1. **Prerequisites Check**
   - Verify Git and GitHub CLI are installed
   - Check authentication

2. **Branch Management**
   - Check current branch
   - Create feature branch if on main/master
   - Use naming: `feature/<descriptive-name>`

3. **Commit and Push**
   - Analyze changes with `git status` and `git diff`
   - Create conventional commit message
   - Push to remote

4. **Create Pull Request**
   - Use `gh pr create` or provide manual instructions
   - Include comprehensive PR description

#### Commit Message Format

Follow conventional commits:

```
<type>(<scope>): <subject>

<body>

<footer>
```

Types: `feat`, `fix`, `refactor`, `style`, `docs`, `test`, `chore`, `perf`

See `.claude/commands/commit-push-pr.md` for the complete process.

## Best Practices

### 1. Reference Project Guidelines

Always mention relevant documentation:

```
@AGENTS.md How should I implement a new validator in the domain layer?
```

### 2. Follow Clean Architecture

Specify the layer when requesting implementations:

- **Domain Layer**: Pure business logic, no infrastructure dependencies
- **Application Layer**: Business logic orchestration (`pkg/wallet/service/`)
- **Infrastructure Layer**: External dependencies (`pkg/infrastructure/`)

### 3. Use Custom Commands

Leverage structured workflows for common tasks:

```
@.claude/commands/fix-issue.md Fix issue #123
```

### 4. Specify Constraints

Always mention project-specific constraints:

- Security requirements (no private key logging)
- Architecture patterns (Clean Architecture, dependency injection)
- Error handling (use `fmt.Errorf` + `%w`)
- Testing requirements

### 5. Request Incremental Changes

Break large tasks into smaller steps:

```
Step 1: Create the domain validator interface
Step 2: Implement the validator in domain layer
Step 3: Update the application layer to use it
```

### 6. Verify Before Committing

Always request verification:

```
After implementing, please run:
- make lint-fix
- make check-build
- make gotest
```

## Workflow Examples

### Example 1: Fixing a GitHub Issue

```
@.claude/commands/fix-issue.md Fix issue #123

The issue is about adding validation for account types in the domain layer.
Reference the existing validators in pkg/domain/transaction/ for the pattern.
```

Claude will:

1. Check git status and create feature branch
2. Fetch issue details
3. Analyze requirements
4. Implement following Clean Architecture
5. Add tests
6. Verify with lint/build/test
7. Commit and create PR

### Example 2: Addressing PR Review Comments

```
@.claude/commands/fix-pr-review.md Fix review comments for PR #456
```

Claude will:

1. Fetch PR and review comments
2. Categorize and prioritize comments
3. Implement fixes systematically
4. Test and verify
5. Commit and push to PR branch

### Example 3: Fixing Linter Errors

```
@.claude/commands/fix-linter.md Fix the linting errors

I ran make lint-fix and got these errors:
[list of errors]
```

Claude will:

1. Analyze error types
2. Prioritize fixes
3. Fix incrementally
4. Verify no new issues

### Example 4: Creating a Feature

```
I need to add a new use case for account validation. 
Follow Clean Architecture principles and reference:
- @AGENTS.md for project guidelines
- pkg/domain/account/ for domain patterns
- pkg/wallet/service/ for application layer patterns

After implementation, run make lint-fix, make check-build, and make gotest.
```

## Safety Rules

### Critical Rules

- **Never work on `main`/`master` branch** - Always create a feature branch
- **Never edit auto-generated files** - Files with "DO NOT EDIT" comments (sqlc, protoc, etc.)
- **Never log private keys** - Security-sensitive information must never be logged
- **Always verify git status** - Check working directory is clean before starting
- **Never use git merge** - Use PR workflow instead

### Security Considerations

- Extra caution for private key management and wallet operations
- Run `make check-vuln` for security-related changes
- Consider impact on offline wallets (keygen, sign)
- Review encryption/decryption logic carefully

## Troubleshooting

### Claude Doesn't Follow Architecture

**Solution**: Explicitly reference `AGENTS.md`:

```
@AGENTS.md Following Clean Architecture, implement this in the domain layer with no infrastructure dependencies.
```

### Claude Modifies Auto-Generated Files

**Solution**: Explicitly mention:

```
DO NOT edit files with "DO NOT EDIT" comments (auto-generated files from sqlc, protoc, etc.)
```

### Claude Uses Panic in Business Logic

**Solution**: Remind about panic usage:

```
Remember: panic is only allowed in main.go and pkg/di. Use proper error handling with error returns.
```

### Verification Commands Fail

**Solution**: Fix issues incrementally:

1. Run `make lint-fix` to fix linting issues
2. Run `make tidy` to organize dependencies
3. Run `make check-build` to verify compilation
4. Run `make gotest` to run tests
5. Address any failures before proceeding

## Additional Resources

- [AGENTS.md](../../AGENTS.md) - Project guidelines and architecture principles
- [REFACTORING_CHECKLIST.md](../../REFACTORING_CHECKLIST.md) - Refactoring guidelines
- [USECASE_LAYER_IMPLEMENTATION_GUIDE.md](../../USECASE_LAYER_IMPLEMENTATION_GUIDE.md) - Use case layer implementation guide
- [Cursor Development Guide](./cursor.md) - Cursor-specific development guide

## Getting Help

If you encounter issues:

1. Check `AGENTS.md` for project-specific guidelines
2. Reference existing code patterns in the codebase
3. Use custom commands for structured workflows
4. Review related documentation files
5. Ask Claude with specific context and code references
