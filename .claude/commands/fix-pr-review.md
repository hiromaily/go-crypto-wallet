# Fix PR Review Comments for PR #{pr_number}

## Repository

Repo: hiromaily/go-crypto-wallet

## Common Workflow Steps

This command follows the common workflow steps defined in [Workflow Guidelines](../../agents/workflow.md):
- **Required Tools and Versions**: See [Required Tools and Versions](../../agents/requirements.md)
- **Pre-Flight Checks**: See [Pre-Flight Checks](../../agents/workflow.md#pre-flight-checks)
- **Safety Rules**: See [Safety Rules](../../agents/workflow.md#safety-rules)
- **Verification Steps**: See [Verification Steps](../../agents/workflow.md#verification-steps)
- **Special Considerations**: See [Special Considerations](../../agents/workflow.md#special-considerations)

## Process

### Pre-Flight Checks

Follow the [Pre-Flight Checks](../../agents/workflow.md#pre-flight-checks) from Workflow Guidelines, with these additional PR-specific checks:

1. **Fetch PR Information:**
   - URL: <https://github.com/hiromaily/go-crypto-wallet/pull/{pr_number}>
   - Use `web_search` to fetch complete PR content including:
     - PR title and description
     - All review comments (both general comments and line-specific comments)
     - Review status and reviewer feedback
     - Discussion threads
   - Verify PR exists and is not already merged/closed
   - Identify the PR's source branch name

2. **Verify Branch:**
   - Check if current branch matches PR source branch
   - If not on the PR branch, checkout the correct branch: `git checkout {pr_branch_name}`
   - If PR branch doesn't exist locally, fetch and checkout: `git fetch origin {pr_branch_name} && git checkout {pr_branch_name}`

### Resolve Review Comments Systematically

1. **Analyze Review Comments:**
   - Categorize comments by type:
     - **Code quality**: Style, naming, structure
     - **Functionality**: Bugs, logic errors, edge cases
     - **Security**: Security concerns, sensitive data handling
     - **Testing**: Missing tests, test improvements
     - **Documentation**: Missing comments, unclear code
     - **Architecture**: Design patterns, Clean Architecture violations
   - Prioritize comments:
     - Security issues (highest priority)
     - Functionality bugs
     - Code quality and architecture
     - Documentation and style (lowest priority)
   - Identify affected files and line numbers
   - Check if comments relate to security-sensitive areas (private keys, wallet operations)
   - Review `AGENTS.md` for project-specific guidelines

2. **Plan Fixes:**
   - Group related comments together
   - Break down fixes into logical steps
   - Identify test cases needed for functionality fixes
   - Check for auto-generated files - see [Auto-Generated Files](../../agents/workflow.md#auto-generated-files) in Workflow Guidelines
   - Consider impact on offline wallet operations (keygen, sign)
   - Plan rollback strategy if breaking changes

3. **Implement Fixes:**
   - Address each review comment systematically
   - Follow Clean Architecture principles
   - Use dependency injection and interfaces
   - Follow coding standards from `AGENTS.md`:
   - Ensure import order: standard → third-party → local
   - For each fix, explain what was changed and why

4. **Self-Review:**
   Follow the [Self-Review](../../agents/workflow.md#self-review) checklist from Workflow Guidelines.

5. **Test:**
   - Run existing tests: `make gotest`
   - Create new test cases for functionality fixes
   - Run integration tests if applicable: `make gotest-integration`
   - Verify test coverage for new code
   - Test edge cases and error scenarios mentioned in reviews

6. **Document:**
   - Update relevant documentation if requested in reviews
   - Add/update code comments as needed
   - Ensure all exported functions have godoc comments

7. **Verify:**
   Follow the [Verification Steps](../../agents/workflow.md#verification-steps) from Workflow Guidelines.

8. **Commit:**
   - Stage changes: `git add <files>`
   - Create commit with descriptive message:

     ```text
     fix(pr): address review comments for PR #{pr_number}
     
     - {fix 1: brief description}
     - {fix 2: brief description}
     - {fix 3: brief description}
     
     Addresses review feedback on PR #{pr_number}
     ```

   - Follow conventional commit format
   - If fixes are extensive, consider multiple atomic commits grouped by category

9. **Push Changes:**
   - Push to PR branch: `git push origin {pr_branch_name}`
   - This will automatically update the existing PR

10. **Update PR (Optional):**
   - Add a comment to the PR summarizing the fixes made
   - Reference specific review comments that were addressed
   - Mark resolved comments if using GitHub's "Resolve conversation" feature

### Safety Rules

Follow the [Safety Rules](../../agents/workflow.md#safety-rules) from Workflow Guidelines.

### Special Considerations

- **Security-Sensitive Changes:**
  See [Security-Sensitive Changes](../../agents/workflow.md#security-sensitive-changes) in Workflow Guidelines.

- **Breaking Changes:**
  See [Breaking Changes](../../agents/workflow.md#breaking-changes) in Workflow Guidelines.
  - If review suggests breaking changes, discuss with reviewer first

- **Conflicting Comments:**
  - If reviewers have conflicting opinions, prioritize security and functionality concerns
  - Consider asking for clarification in PR comments if needed

### Output Format

For each fix iteration:

1. **Show the review comment** being addressed (quote the comment)
2. **Describe the fix** being implemented
3. **Show the specific changes** made (code diff or explanation)
4. **Explain the reasoning** behind the fix
5. **Note any remaining comments** to address in subsequent steps

### Example Workflow

1. Fetch PR #123: `web_search` for "github.com/hiromaily/go-crypto-wallet/pull/123"
2. Identify 5 review comments across 3 files
3. Categorize: 2 security, 2 code quality, 1 documentation
4. Fix security issues first
5. Fix code quality issues
6. Add documentation
7. Run verification commands
8. Commit and push

Remember: Quality over speed. Address each review comment thoroughly and ensure all fixes maintain code quality and project standards.
