# Fix Issue #{issue_number} [or multiple: #{issue_number1},#{issue_number2},...]

## Repository

Repo: hiromaily/go-crypto-wallet

## Issue Number Format

- **Single issue**: `#123` or `123`
- **Multiple issues**: `#123,#124,#125` or `123,124,125` (comma-separated)
- **Multiple issues with spaces**: `#123 #124 #125` or `123 124 125` (space-separated)

When multiple issues are provided:

- Issues are processed **in the order specified**
- Each issue gets its own **separate commit** in the same branch
- All commits are included in **one PR** that closes all issues
- Branch name is based on the first issue number

## Common Workflow Steps

This command follows the common workflow steps defined in [Workflow Guidelines](../../agents/workflow.md):

- **Required Tools and Versions**: See [Required Tools and Versions](../../agents/requirements.md)
- **Pre-Flight Checks**: See [Pre-Flight Checks](../../agents/workflow.md#pre-flight-checks)
- **Safety Rules**: See [Safety Rules](../../agents/workflow.md#safety-rules)
- **Verification Steps**: See [Verification Steps](../../agents/workflow.md#verification-steps)
- **Special Considerations**: See [Special Considerations](../../agents/workflow.md#special-considerations)

## Process

### Pre-Flight Checks

Follow the [Pre-Flight Checks](../../agents/workflow.md#pre-flight-checks) from Workflow Guidelines, with these additional issue-specific checks:

1. **Parse Issue Numbers:**
   - Extract issue numbers from input (handle comma or space-separated format)
   - Remove `#` prefix if present
   - Store issue numbers in an ordered list: `[issue_number1, issue_number2, ...]`

2. **Fetch All Issues:**
   - For each issue number, use `gh issue view {issue_number}` to fetch complete issue content
   - Verify each issue exists and is not already closed/assigned
   - Review issue descriptions, comments, and labels
   - Store issue information for later reference
   - **If any issue is invalid or already closed**: Stop and report which issue(s) have problems

3. **Create Feature Branch:**
   - **Single issue**: Format: `feature/issue-{issue_number}-{brief-description}`
   - **Multiple issues**: Format: `feature/issue-{first_issue_number}-{brief-description}` (use first issue's description)
   - Example (single): `feature/issue-123-fix-logger-global-issue`
   - Example (multiple): `feature/issue-123-fix-multiple-issues`
   - Keep description concise and descriptive
   - Create and checkout branch: `git checkout -b feature/issue-{issue_number}-{description}`

### Resolve Systematically

**For Multiple Issues**: Process each issue **sequentially in the specified order**. Each issue gets its own commit, but all commits go into the same branch and PR.

#### Process Each Issue (Loop for Multiple Issues)

For each issue in the ordered list:

1. **Analyze:**
   - Read issue description, comments, and related discussions
   - Understand problem, root cause, and requirements
   - Check if issue relates to security-sensitive areas (private keys, wallet operations)
   - Review `AGENTS.md` for project-specific guidelines
   - Identify affected files and components
   - Consider impact on offline wallet operations (keygen, sign)
   - **For multiple issues**: Check if current issue conflicts with previous issues' changes

2. **Plan:**
   - Break down solution into steps
   - Identify test cases needed
   - Check for auto-generated files - see [Auto-Generated Files](../../agents/workflow.md#auto-generated-files) in Workflow Guidelines
   - Consider backward compatibility
   - Plan rollback strategy if breaking changes
   - **If the issue is too large or complex**: Stop processing and propose creating sub-issues
     - Assess if the issue requires multiple PRs or spans multiple components
     - If too large, interrupt processing and suggest breaking it down into smaller sub-issues
     - Use `gh issue create` to create sub-issues with appropriate scope
     - Link sub-issues to the parent issue
     - Proceed with fixing only after the issue is broken down into manageable pieces

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
   - **For multiple issues**: Ensure changes don't conflict with previous issues' implementations

4. **Self-Review:**
   Follow the [Self-Review](../../agents/workflow.md#self-review) checklist from Workflow Guidelines.

5. **Test:**
   - Run existing tests: `make gotest`
   - Create new test cases for the fix
   - Run integration tests if applicable: `make gotest-integration`
   - Verify test coverage for new code
   - Test edge cases and error scenarios
   - **For multiple issues**: Ensure tests still pass after each issue's changes

6. **Document:**
   - Update relevant documentation (README, API docs, etc.)
   - Add/update code comments as needed

7. **Verify:**
   Follow the [Verification Steps](../../agents/workflow.md#verification-steps) from Workflow Guidelines.
   - For security-related changes, also run: `make check-vuln`
   - **For multiple issues**: Run verification after each issue's changes

8. **Commit (Per Issue):**
   - Stage changes for this specific issue: `git add <files>`
   - Create commit with descriptive message for this issue:

     ```text
     fix: resolve issue #{issue_number} - {brief description}
     
     - {detail 1}
     - {detail 2}
     
     Closes #{issue_number}
     ```

   - Follow conventional commit format when appropriate
   - **CRITICAL**: Each issue must have its own separate commit
   - **CRITICAL**: Commit messages must reference the specific issue number being fixed

9. **Continue to Next Issue (if multiple):**
   - If there are more issues, proceed to the next issue in the list
   - Repeat steps 1-8 for the next issue
   - Each issue gets its own commit on the same branch
   - **Do not create PR until all issues are processed**

10. **PR Draft (After All Issues):**

- Push branch: `git push origin feature/issue-{first_issue_number}-{description}`
- Create PR using `gh pr create`:
  - **Single issue**: Title: `Fix: {issue title} (Closes #{issue_number})`
  - **Multiple issues**: Title: `Fix: {brief description} (Closes #{issue_number1}, #{issue_number2})`
  - Description template:

       ```markdown
       ## Description
       {Brief description of the fixes}
       
       This PR resolves multiple related issues:
       - Issue #{issue_number1}: {brief description}
       - Issue #{issue_number2}: {brief description}
       - ...
       
       ## Changes
       - {Change 1 from issue 1}
       - {Change 2 from issue 2}
       - ...
       
       ## Testing
       - [ ] Unit tests added/updated
       - [ ] Integration tests pass
       - [ ] Manual testing completed
       
       ## Verification
       - [ ] `make lint-fix` passes
       - [ ] `make tidy` passes
       - [ ] `make check-build` passes
       - [ ] `make gotest` passes
       - [ ] Security scan completed (if applicable)
       
       ## Commits
       - Each issue has been resolved in a separate commit
       - Commits are in the order issues were specified
       
       Closes #{issue_number1}
       Closes #{issue_number2}
       ...
       ```

  - Command example:

       ```bash
       gh pr create --title "Fix: {description} (Closes #{issue_number1}, #{issue_number2})" --body-file - <<EOF
       {paste description template here}
       EOF
       ```

  - Or use interactive mode: `gh pr create` (will prompt for title and body)
  - Link related issues/PRs if any

1. **Review Request:**

   After creating PR, ask Claude to review the following:

- Code quality and correctness
- Adherence to project standards (`AGENTS.md`)
- Security implications (especially for wallet/key operations)
- Test coverage adequacy
- Documentation completeness
- **For multiple issues**: Verify that each issue's changes are properly separated in commits

### Sub-Issue Resolution Workflow

**Note**: This workflow is for when sub-issues need to be resolved **separately** (each in its own PR).
If you want to resolve multiple sub-issues in **one PR** with separate commits, use the main workflow above with multiple issue numbers.

When a parent issue has multiple sub-issues that need separate PRs, resolve them sequentially following this workflow:

1. **First Sub-Issue:**
   - Create feature branch: `git checkout -b feature/issue-{sub_issue_number}-{description}`
   - Follow steps 3-10 from "Resolve Systematically" section:
     - Implement the fix
     - Self-review
     - Test
     - Document
     - Verify (run `make lint-fix`, `make tidy`, `make check-build`, `make gotest`)
     - Commit changes
     - Create PR with description referencing the sub-issue
     - Request review
   - **Stop here and wait for user confirmation and merge**

2. **Wait for User Action:**
   - User will review the PR
   - User will merge if approved
   - **Do not proceed to next sub-issue until user confirms merge is complete**

3. **Next Sub-Issue:**
   - Ensure working directory is clean (`git status`)
   - Ensure you're on `main` or `master` branch (or switch to it: `git checkout main`)
   - Pull latest changes: `git pull origin main` (or `master`)
   - Create **new** feature branch for the next sub-issue: `git checkout -b feature/issue-{next_sub_issue_number}-{description}`
   - Follow steps 3-10 again for the next sub-issue
   - **Repeat steps 2-3 for each remaining sub-issue**

**Important Rules for Sub-Issue Workflow:**

- **CRITICAL**: Always create a new branch for each sub-issue (do not reuse branches)
- **CRITICAL**: Always start from clean state (clean working directory, latest main/master)
- **CRITICAL**: Wait for user confirmation before proceeding to next sub-issue
- Each sub-issue should be independent and mergeable separately
- PR title should reference the sub-issue: `Fix: {sub-issue title} (Closes #{sub_issue_number})`
- Link sub-issue PR to parent issue in PR description

**Alternative: Multiple Sub-Issues in One PR**

If sub-issues are related and should be resolved together:

- Use the main workflow with multiple issue numbers: `#123,#124,#125`
- Each sub-issue will get its own commit in the same PR
- All sub-issues will be closed when the PR is merged

### Safety Rules

Follow the [Safety Rules](../../agents/workflow.md#safety-rules) from Workflow Guidelines.

### Special Considerations

- **Security-Sensitive Changes:**
  See [Security-Sensitive Changes](../../agents/workflow.md#security-sensitive-changes) in Workflow Guidelines.

- **Breaking Changes:**
  See [Breaking Changes](../../agents/workflow.md#breaking-changes) in Workflow Guidelines.

<!-- - **Multi-Chain Support:**
  - Verify changes work for all supported chains (BTC, BCH, ETH, XRP)
  - Test ERC-20 token operations if ETH-related -->

## Completion Checklist

After completing all steps, report the completion status to the user using the following checklist format:

### Issue Resolution Status

**For Each Issue (if multiple):**

- [ ] **Issue #{issue_number} - Step 1 - Analyze**: Issue content analyzed, problem understood, affected files identified
- [ ] **Issue #{issue_number} - Step 2 - Plan**: Solution broken down into steps, test cases identified
- [ ] **Issue #{issue_number} - Step 3 - Implement**: Code changes implemented following Clean Architecture principles
- [ ] **Issue #{issue_number} - Step 4 - Self-Review**: Code reviewed for quality, architecture compliance, and security
- [ ] **Issue #{issue_number} - Step 5 - Test**: Tests run and passed, new test cases created
- [ ] **Issue #{issue_number} - Step 6 - Document**: Documentation updated as needed
- [ ] **Issue #{issue_number} - Step 7 - Verify**: All verification commands passed (`make lint-fix`, `make tidy`, `make check-build`, `make gotest`)
- [ ] **Issue #{issue_number} - Step 8 - Commit**: Changes committed with appropriate commit message

**After All Issues:**

- [ ] **Step 9 - PR Draft**: Pull request created with complete description (includes all issues)
- [ ] **Step 10 - Review Request**: PR ready for review

### Summary

Provide a brief summary including:

- Issue number(s) and title(s)
- Branch name created
- PR number (if created)
- Key changes made for each issue
- Commit structure (if multiple issues)
- Any special considerations or notes

**Example completion message (single issue):**

```text
✅ Issue #123 has been resolved

All steps completed:
✓ Analyzed issue and identified affected components
✓ Planned solution and test cases
✓ Implemented fixes following Clean Architecture
✓ Self-reviewed code for quality and security
✓ Tests created and passing
✓ Documentation updated
✓ All verification commands passed
✓ Changes committed
✓ Pull request #456 created
✓ Ready for review

Branch: feature/issue-123-fix-logger-global-issue
PR: #456
Key changes: Fixed logger initialization issue in domain layer
```

**Example completion message (multiple issues):**

```text
✅ Issues #123, #124, #125 have been resolved

All steps completed:
✓ Analyzed all issues and identified affected components
✓ Planned solutions and test cases for each issue
✓ Implemented fixes following Clean Architecture
✓ Self-reviewed code for quality and security
✓ Tests created and passing
✓ Documentation updated
✓ All verification commands passed
✓ Changes committed (3 separate commits, one per issue)
✓ Pull request #456 created
✓ Ready for review

Branch: feature/issue-123-fix-multiple-issues
PR: #456

Commits:
1. fix: resolve issue #123 - Fix logger initialization
2. fix: resolve issue #124 - Update error handling
3. fix: resolve issue #125 - Add domain entity conversion

Key changes:
- Issue #123: Fixed logger initialization issue in domain layer
- Issue #124: Updated error handling in repository layer
- Issue #125: Added domain entity conversion in infrastructure layer
```
