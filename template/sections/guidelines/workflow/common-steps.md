### Common Workflow Steps

This section describes common workflow steps that are shared across multiple development tasks.
These steps should be followed when implementing fixes, addressing PR reviews, or making code changes.

#### Required Tools and Versions

See [Required Tools and Versions](../../../docs/guidelines/requirements.md) for complete information about:

- Essential tools (Git, GitHub CLI, Go)
- Development tools (golangci-lint, Atlas, markdownlint-cli, Docker)
- Version requirements and verification commands
- Installation instructions

**Important**: Always verify tool versions before starting work.
Using incorrect versions (especially Atlas v1.0.0) may cause compatibility issues.

#### Pre-Flight Checks

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

#### Safety Rules

**CRITICAL Rules (Must Always Follow):**

- **CRITICAL**: Stop immediately if working directory is not clean
- **CRITICAL**: Never proceed on `main` branch without creating feature branch
- **CRITICAL**: Always verify branch and status before implementing fixes
- **CRITICAL**: Never edit files with `DO NOT EDIT` comments (auto-generated files)
- **CRITICAL**: Never log private keys or sensitive information
- **CRITICAL**: For security-related changes, run `make go-check-vuln` and conduct security review
- Never use `git merge` operations
- Never commit/push directly to `main` branch

#### Self-Review

After implementing code changes and before running tests, perform a self-review of your implementation:

**Review Checklist:**

- **Code quality and correctness**: Ensure the code works as intended
- **Adherence to Clean Architecture principles**: Verify layer separation and dependency direction
- **Compliance with coding standards**: Follow guidelines from `AGENTS.md` and [Coding Standards](../../../docs/guidelines/coding-conventions.md)
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

#### Verification Steps

Before committing changes, determine what type of files were changed and run the appropriate verification commands.

##### For Go Code Changes

If Go files (`.go`) were changed, run these commands in order and ensure:

- No errors occur
- No files are modified (all changes should be committed)
- All commands pass successfully:

```bash
make go-lint       # Fix linting issues (not 'fix-lint')
make tidy          # Organize dependencies
make check-build   # Verify builds successfully
make go-test        # Run all tests
```

**Optional but Recommended:**

- `make go-check-vuln` - Run security vulnerability scan (for security-related changes)
- `make go-test-integration` - Run integration tests (if applicable)

##### For Markdown File Changes Only

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

##### For Mixed Changes

If both Go files and Markdown files were changed:

1. Run Go verification commands (as above)
2. Run markdown linting (as above)

#### Special Considerations

##### Security-Sensitive Changes

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

##### Auto-Generated Files

**CRITICAL**: Never edit files with `DO NOT EDIT` comments:

- SQLC generated files (`internal/infrastructure/database/sqlc/`)
- Protocol buffer generated files
- Files generated by `go generate`

See [Code Generation Guidelines](../../../docs/guidelines/code-generation.md) for details.

##### Breaking Changes

- Document breaking changes clearly
- Consider migration path
- Update version numbers if applicable
- Implement incrementally with rollback plans
