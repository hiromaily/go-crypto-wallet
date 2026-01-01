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
- `make check-vuln`: Run security vulnerability scan (govulncheck)

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
- ❌ `git commit` and `git push` to `main` or `master` branches

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

1. `make lint-fix` - Fix linting issues automatically
2. `make tidy` - Organize dependencies and clean up `go.mod`
3. `make check-build` - Verify that the code builds successfully
4. `make gotest` - Run Go tests to verify functionality

**Optional but Recommended:**

- `make check-vuln` - Run security vulnerability scan (for security-related changes)
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
- [ ] `make lint-fix` passes
- [ ] `make check-build` passes
- [ ] `make gotest` passes
- [ ] Security scan completed (if applicable)

Closes #[issue_number]
```

## Branching Strategy

**Branch Naming:**

- `feature/issue-{number}-{brief-description}` - For new features
- `fix/issue-{number}-{brief-description}` - For bug fixes
- `refactor/issue-{number}-{brief-description}` - For refactoring

**Branch Lifecycle:**

1. Create branch from `main`/`master`
2. Make changes and commit
3. Push to GitHub
4. Create pull request
5. Address review comments
6. Wait for merge by maintainers
7. Delete branch after merge

## Important Notes

- This is a financial-related project; make changes carefully
- Implement breaking changes incrementally with rollback plans
- Security-related changes must be reviewed
- Always verify that changes don't break existing functionality
- Consider the impact on offline wallet operations (keygen, sign)
- **DO NOT** edit files that contain `DO NOT EDIT` comments (auto-generated files)

## See Also

- [Core Principles](core.md) - Security and important notes
- [Code Generation Guidelines](code-generation.md) - Working with auto-generated files
- [Coding Standards](coding-standards.md) - Linting and formatting commands
- [Testing Guidelines](testing.md) - Testing requirements and strategies
