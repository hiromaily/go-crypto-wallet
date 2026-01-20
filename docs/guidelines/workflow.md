# Workflow Guidelines

This document describes development workflow, dependency management, and git operations for the go-crypto-wallet project.

## Refactoring Status

- Make changes incrementally without breaking existing functionality
- Run tests before and after refactoring
- Follow the phased approach outlined in refactoring documents

**Important**: This is an ongoing refactoring project moving toward Clean Architecture. When working on code:

1. Check if the code is part of a planned refactoring
2. Follow the refactoring plan if it exists
3. Don't introduce new patterns that conflict with the target architecture
4. When in doubt, ask about refactoring priorities

## Dependency Management

**Go Modules:**

- Use `go mod tidy` to organize dependencies
- Run security scans (`govulncheck`)
- Keep dependencies up-to-date while maintaining stability

**Commands:**

- `make tidy`: Organize dependencies and clean up `go.mod`
- `make go-check-vuln`: Run security vulnerability scan (govulncheck)

**Best Practices:**

- Review dependency changes carefully
- Test thoroughly after dependency updates
- Document breaking changes in dependencies
- Consider security implications of new dependencies

## Git Operations

**Allowed Operations:**

- `git add`: Stage changes for commit
- `git commit`: Commit staged changes
- `git push`: Push commits to GitHub (remote repository)

**NOT Allowed Operations:**

- ❌ `git merge` operations
- ❌ `gh` command merge operations (e.g., `gh pr merge`)
- ❌ `git commit` and `git push` to `main` branch

**Workflow:**

1. Create a feature branch for your work
2. Make changes and commit to the feature branch
3. Push the feature branch to GitHub
4. Create a pull request
5. Wait for review and approval
6. Let maintainers merge the pull request

**Commit Messages:**

Follow conventional commit format when appropriate:

- `fix:` for bug fixes
- `feat:` for new features
- `refactor:` for refactoring changes
- `docs:` for documentation changes
- `test:` for test additions/changes
- `chore:` for maintenance tasks

Example:

```text
fix: resolve logger initialization issue in domain layer

- Move logger initialization to dependency injection
- Remove global logger usage from domain package
- Add logger interface for testability

Closes #123
```

## Verification Commands

After making code changes, always run these commands in order:

1. `make go-lint` - Fix linting issues automatically
2. `make tidy` - Organize dependencies and clean up `go.mod`
3. `make check-build` - Verify that the code builds successfully
4. `make gotest` - Run Go tests to verify functionality

**Optional but Recommended:**

- `make go-check-vuln` - Run security vulnerability scan (for security-related changes)
- `make gotest-integration` - Run integration tests (if applicable)

**Important**:

- Ensure no errors occur
- Ensure no files are modified (all changes should be committed)
- Ensure all commands pass successfully

## Pull Request Workflow

**Before Creating a PR:**

1. Run all verification commands
2. Ensure all tests pass
3. Review your own changes for:
   - Code quality and correctness
   - Adherence to Clean Architecture principles
   - Compliance with coding standards
   - Proper error handling
   - Security considerations

**PR Description Should Include:**

- Brief description of the changes
- Related issue number(s)
- Testing performed
- Any breaking changes
- Screenshots (if UI changes)

**PR Template:**

```markdown
## Description
[Brief description of the fix/feature]

## Changes
- [Change 1]
- [Change 2]

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests pass
- [ ] Manual testing completed

## Verification
- [ ] `make go-lint` passes
- [ ] `make check-build` passes
- [ ] `make gotest` passes
- [ ] Security scan completed (if applicable)

Closes #[issue_number]
```

## Branching Strategy

**Branch Naming:**

Format: `{type}/{brief-description}-{issue-number}`

- `feature/{brief-description}-{issue-number}` - For new features
- `fix/{brief-description}-{issue-number}` - For bug fixes
- `refactor/{brief-description}-{issue-number}` - For refactoring

See [docs/standards/workflow.md](../standards/workflow.md) for the complete branch naming convention.

**Branch Lifecycle:**

1. Create branch from `main`
2. Make changes and commit
3. Push to GitHub
4. Create pull request
5. Address review comments
6. Wait for merge by maintainers
7. Delete branch after merge

## Common Workflow Steps

This section describes common workflow steps that are shared across multiple development tasks.
These steps should be followed when implementing fixes, addressing PR reviews, or making code changes.

### Required Tools and Versions

See [Required Tools and Versions](requirements.md) for complete information about:

- Essential tools (Git, GitHub CLI, Go)
- Development tools (golangci-lint, Atlas, markdownlint-cli, Docker)
- Version requirements and verification commands
- Installation instructions

**Important**: Always verify tool versions before starting work.
Using incorrect versions (especially Atlas v1.0.0) may cause compatibility issues.

### Pre-Flight Checks

Before starting any development task, perform these checks:

1. **Check Git Status:**
   - Verify working directory is clean (`git status`)
   - Stop immediately if there are uncommitted changes
   - Check current branch (`git branch --show-current`)
   - Never proceed on `main` branch without creating a feature branch

2. **Verify Branch:**
   - Ensure you're on the correct branch for your task
   - If working on a PR, verify branch matches PR source branch
   - If working on an issue, create a feature branch first

3. **Review Project Guidelines:**
   - Review `AGENTS.md` for project-specific guidelines
   - Check relevant architecture and coding standards documents
   - Understand security requirements for the task

### Safety Rules

**CRITICAL Rules (Must Always Follow):**

- **CRITICAL**: Stop immediately if working directory is not clean
- **CRITICAL**: Never proceed on `main` branch without creating feature branch
- **CRITICAL**: Always verify branch and status before implementing fixes
- **CRITICAL**: Never edit files with `DO NOT EDIT` comments (auto-generated files)
- **CRITICAL**: Never log private keys or sensitive information
- **CRITICAL**: For security-related changes, run `make go-check-vuln` and conduct security review
- Never use `git merge` operations
- Never commit/push directly to `main` branch

### Self-Review

After implementing code changes and before running tests, perform a self-review of your implementation:

**Review Checklist:**

- **Code quality and correctness**: Ensure the code works as intended
- **Adherence to Clean Architecture principles**: Verify layer separation and dependency direction
- **Compliance with coding standards**: Follow guidelines from `AGENTS.md` and [Coding Standards](coding-standards.md)
- **Proper error handling**: Use `fmt.Errorf` + `%w` for error wrapping, add `context.Context` to API calls
- **Security considerations**: Especially important for wallet/key operations
  - Never log private keys or sensitive information
- **Code organization**:
  - Import order: standard → third-party → local
  - Remove unused code, variables, or functions
  - Add godoc comments to exported functions/methods
- **Architecture compliance**:
  - Use dependency injection and interfaces
  - Domain layer has ZERO infrastructure dependencies
  - Follow interface segregation principles

**After Self-Review:**

- Fix any issues found during self-review
- Ensure all changes align with project guidelines
- Proceed to testing and verification steps

### Verification Steps

Before committing changes, determine what type of files were changed and run the appropriate verification commands.

#### For Go Code Changes

If Go files (`.go`) were changed, run these commands in order and ensure:

- No errors occur
- No files are modified (all changes should be committed)
- All commands pass successfully:

```bash
make go-lint       # Fix linting issues (not 'fix-lint')
make tidy          # Organize dependencies
make check-build   # Verify builds successfully
make gotest        # Run all tests
```

**Optional but Recommended:**

- `make go-check-vuln` - Run security vulnerability scan (for security-related changes)
- `make gotest-integration` - Run integration tests (if applicable)

#### For Markdown File Changes Only

If **only** Markdown files (`.md`) were changed (no Go code changes), run markdown linting:

```bash
# Using markdownlint-cli (if installed via npm/npx)
npx markdownlint-cli "**/*.md" --config .markdownlint.json

# Or using markdownlint command (if installed globally)
markdownlint "**/*.md" --config .markdownlint.json
```

**Note:**

- If markdownlint is not installed, install it first: `npm install -g markdownlint-cli` or use `npx markdownlint-cli`
- The project uses `.markdownlint.json` for configuration
- Only markdown files need to be linted; Go verification commands are not required for markdown-only changes

#### For Mixed Changes

If both Go files and Markdown files were changed:

1. Run Go verification commands (as above)
2. Run markdown linting (as above)

### Special Considerations

#### Security-Sensitive Changes

For issues involving:

- Private key management
- Wallet operations
- Authentication/authorization
- Encryption/decryption

**Additional requirements:**

- Extra caution for private key management, wallet operations
- Run security scan: `make go-check-vuln`
- Consider impact on offline wallets (keygen, sign)
- Review encryption/decryption logic carefully
- Never include sensitive information in commits or PR descriptions

#### Auto-Generated Files

**CRITICAL**: Never edit files with `DO NOT EDIT` comments:

- SQLC generated files (`internal/infrastructure/database/sqlc/`)
- Protocol buffer generated files
- Files generated by `go generate`

See [Code Generation Guidelines](code-generation.md) for details.

#### Breaking Changes

- Document breaking changes clearly
- Consider migration path
- Update version numbers if applicable
- Implement incrementally with rollback plans

## Important Notes

- This is a financial-related project; make changes carefully
- Implement breaking changes incrementally with rollback plans
- Security-related changes must be reviewed
- Always verify that changes don't break existing functionality
- Consider the impact on offline wallet operations (keygen, sign)
- **DO NOT** edit files that contain `DO NOT EDIT` comments (auto-generated files)

## See Also

- [Required Tools and Versions](requirements.md) - Tool requirements and version information
- [Core Principles](core.md) - Security and important notes
- [Code Generation Guidelines](code-generation.md) - Working with auto-generated files
- [Coding Standards](coding-standards.md) - Linting and formatting commands
- [Testing Guidelines](testing.md) - Testing requirements and strategies
