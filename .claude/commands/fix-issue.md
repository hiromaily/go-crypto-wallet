# Fix Issue #{issue_number}

## Repository

Repo: hiromaily/go-crypto-wallet

## Process

### Pre-Flight Checks

1. **Check Git Status:**
   - Verify working directory is clean (`git status`)
   - Stop immediately if there are uncommitted changes
   - Check current branch (`git branch --show-current`)
   - Never proceed on `main` or `master` branch

2. **Fetch Issue:**
   - Use `gh issue view {issue_number}` to fetch complete issue content
   - Verify issue exists and is not already closed/assigned
   - Review issue description, comments, and labels

3. **Create Feature Branch:**
   - Format: `feature/issue-{issue_number}-{brief-description}`
   - Example: `feature/issue-123-fix-logger-global-issue`
   - Keep description concise and descriptive
   - Create and checkout branch: `git checkout -b feature/issue-{issue_number}-{description}`

### Resolve Systematically

1. **Analyze:**
   - Read issue description, comments, and related discussions
   - Understand problem, root cause, and requirements
   - Check if issue relates to security-sensitive areas (private keys, wallet operations)
   - Review `AGENTS.md` for project-specific guidelines
   - Identify affected files and components
   - Consider impact on offline wallet operations (keygen, sign)

2. **Plan:**
   - Break down solution into steps
   - Identify test cases needed
   - Check for auto-generated files (sqlc, protoc, go generate) - **DO NOT EDIT** these
   - Consider backward compatibility
   - Plan rollback strategy if breaking changes

3. **Implement:**
   - Follow Clean Architecture principles
   - Use dependency injection and interfaces
   - Follow coding standards from `AGENTS.md`:
     - Use `fmt.Errorf` + `%w` for error wrapping
     - Add `context.Context` to API calls
     - Never log private keys or sensitive information
     - Use structured logging
   - Add godoc comments to exported functions/methods
   - Ensure import order: standard → third-party → local

4. **Test:**
   - Run existing tests: `make gotest`
   - Create new test cases for the fix
   - Run integration tests if applicable: `make gotest-integration`
   - Verify test coverage for new code
   - Test edge cases and error scenarios

5. **Document:**
   - Update relevant documentation (README, API docs, etc.)
   - Add/update code comments as needed

6. **Verify:**
   Before committing, if Go files were changed, run these commands in order and ensure:
   - No errors occur
   - No files are modified (all changes should be committed)
   - All commands pass successfully:

     ```bash
     make lint-fix      # Fix linting issues (not 'fix-lint')
     make tidy          # Organize dependencies
     make check-build   # Verify builds successfully
     make gotest        # Run all tests
     make check-vuln    # Security vulnerability scan (if security-related)
     ```

7. **Commit:**
   - Stage changes: `git add <files>`
   - Create commit with descriptive message:

     ```text
     fix: resolve issue #{issue_number} - {brief description}
     
     - {detail 1}
     - {detail 2}
     
     Closes #{issue_number}
     ```

   - Follow conventional commit format when appropriate

8. **PR Draft:**
   - Push branch: `git push origin feature/issue-{issue_number}-{description}`
   - Create PR using `gh pr create`:
     - Title: `Fix: {issue title} (Closes #{issue_number})`
     - Description template:

       ```markdown
       ## Description
       {Brief description of the fix}
       
       ## Changes
       - {Change 1}
       - {Change 2}
       
       ## Testing
       - [ ] Unit tests added/updated
       - [ ] Integration tests pass
       - [ ] Manual testing completed
       
       ## Verification
       - [ ] `make lint-fix` passes
       - [ ] `make check-build` passes
       - [ ] `make gotest` passes
       - [ ] Security scan completed (if applicable)
       
       Closes #{issue_number}
       ```

     - Command example:

       ```bash
       gh pr create --title "Fix: {issue title} (Closes #{issue_number})" --body-file - <<EOF
       {paste description template here}
       EOF
       ```

     - Or use interactive mode: `gh pr create` (will prompt for title and body)
     - Link related issues/PRs if any

9. **Review Request:**
   - After creating PR, ask Claude to review:
     - Code quality and correctness
     - Adherence to project standards (`AGENTS.md`)
     - Security implications (especially for wallet/key operations)
     - Test coverage adequacy
     - Documentation completeness

### Safety Rules

- **CRITICAL**: Stop immediately if working directory is not clean
- **CRITICAL**: Never proceed on `main`/`master` branch without creating feature branch
- **CRITICAL**: Always verify branch and status before implementing fixes
- **CRITICAL**: Never edit files with `DO NOT EDIT` comments (auto-generated files)
- **CRITICAL**: Never log private keys or sensitive information
- **CRITICAL**: For security-related changes, run `make check-vuln` and conduct security review
- Never use `git merge` operations
- Never commit/push directly to `main`/`master` branches

### Special Considerations

- **Security-Sensitive Changes:**
  - Extra caution for private key management, wallet operations
  - Run security scan: `make check-vuln`
  - Consider impact on offline wallets (keygen, sign)
  - Review encryption/decryption logic carefully

- **Breaking Changes:**
  - Document breaking changes clearly
  - Consider migration path
  - Update version numbers if applicable

<!-- - **Multi-Chain Support:**
  - Verify changes work for all supported chains (BTC, BCH, ETH, XRP)
  - Test ERC-20 token operations if ETH-related -->
