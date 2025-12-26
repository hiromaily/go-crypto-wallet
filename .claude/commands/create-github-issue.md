# Create GitHub Issue

You are tasked with creating a GitHub issue using the `gh` command-line tool. This command helps create well-structured GitHub issues that are suitable for tracking features, bugs, refactoring tasks, and technical debt.

## Repository

Repo: hiromaily/go-crypto-wallet

## Prerequisites Check

Before starting, verify all required tools are installed:

1. **GitHub CLI (gh)**: Check with `gh --version`
   - Required for creating issues via CLI
   - If missing: Install from <https://cli.github.com/>
   - Verify authentication: `gh auth status`
   - If not authenticated: Run `gh auth login`

2. **Git**: Check with `git --version`
   - Required for repository operations
   - If missing: Install from <https://git-scm.com/>

**If any required tool is missing, stop and display an error message with installation instructions. Do not proceed with the workflow.**

## Context Understanding

Before creating an issue, understand the project context by reviewing relevant documentation:

1. **Project Overview** (`AGENTS.md`, `README.md`):
   - Cryptocurrency wallet implementation in Go
   - Supports BTC, BCH, ETH, XRP, and ERC-20 tokens
   - Currently under refactoring based on Clean Architecture principles
   - Security is of utmost importance (private key management, offline wallets)
   - Three wallet types: watch (online), keygen (offline), sign (offline)

2. **Architecture** (`AGENTS.md`):
   - Clean Architecture with clear layer separation
   - Domain layer: Pure business logic (ZERO infrastructure dependencies)
   - Application layer: Use case implementations (`internal/application/usecase/`)
   - Infrastructure layer: External dependencies (`internal/infrastructure/`)
   - Interface adapters layer: CLI commands and wallet adapters (`internal/interface-adapters/`)

## Issue Creation Process

### 1. Gather Information

Before creating the issue, gather the following information from the user or conversation context:

- **Issue Type**: Feature request, bug report, refactoring task, documentation, security, technical debt
- **Title**: Clear, concise description (50-72 characters recommended)
- **Description**: Detailed explanation of what needs to be implemented or fixed
- **Priority**: Critical, High, Medium, Low (based on impact and urgency)
- **Affected Components**: Which layers/components are affected (domain, application, infrastructure, interface-adapters)
- **Related Context**: Related issues, PRs, or documentation references

### 2. Determine Issue Category

Based on the issue type, determine appropriate labels and structure:

- **Feature Request**: New functionality or enhancement
- **Bug Report**: Something that's broken or not working as expected
- **Refactoring**: Code improvement without changing functionality
- **Documentation**: Documentation updates or improvements
- **Security**: Security-related issues or improvements
- **Technical Debt**: Code quality improvements, cleanup tasks

### 3. Structure the Issue Body

Create a well-structured issue body following this template:

```markdown
## Description

[Clear, detailed description of what needs to be implemented or fixed]

## Context

[Why this issue exists, what problem it solves, or what improvement it brings]

## Acceptance Criteria

- [ ] Criterion 1 (specific, testable condition)
- [ ] Criterion 2
- [ ] Criterion 3

## Technical Requirements

### Architecture
- [Layer(s) affected: domain, application, infrastructure, interface-adapters]
- [Clean Architecture principles to follow]
- [Dependency direction considerations]

### Implementation Details
- **Files to modify/create**: [List of files with paths]
- **Related code**: [Code references with line numbers if applicable]
- **Dependencies**: [Any new dependencies needed]

### Constraints
- [Security requirements if applicable]
- [Performance considerations]
- [Backward compatibility requirements]
- [Impact on offline wallet operations (keygen, sign) if applicable]

## Testing Requirements

- **Unit tests**: [What to test at unit level]
- **Integration tests**: [What to test at integration level, if applicable]
- **Manual testing**: [Steps to verify manually]

## Related Context

- **Related issues**: [Issue numbers, if any]
- **Related PRs**: [PR numbers, if any]
- **Documentation**: [Links to relevant docs (AGENTS.md, README.md, etc.)]
- **Architecture guidelines**: [References to specific sections in AGENTS.md]

## Additional Notes

[Any additional context, considerations, or constraints]
```

### 4. Determine Labels

Based on the issue type and content, suggest appropriate labels:

- **Type labels**: `bug`, `enhancement`, `refactoring`, `documentation`, `security`, `technical-debt`
- **Priority labels**: `priority:critical`, `priority:high`, `priority:medium`, `priority:low`
- **Component labels**: `domain`, `application`, `infrastructure`, `interface-adapters`, `btc`, `eth`, `xrp`, `bch`
- **Status labels**: `help-wanted`, `good-first-issue`, `blocked`, `in-progress`

**Note**: Labels may vary by repository. Use `gh label list` to see available labels before creating the issue.

### 5. Create the Issue

Use the `gh issue create` command with appropriate flags:

```bash
gh issue create \
  --title "Issue Title" \
  --body-file issue_body.md \
  --label "label1,label2,label3"
```

Or use interactive mode:

```bash
gh issue create
```

**Interactive mode prompts for:**

- Title
- Body (can paste markdown or use editor)
- Labels (comma-separated)
- Assignees (optional)
- Projects (optional)
- Milestone (optional)

### 6. Alternative: Create Issue from File

If the issue body is long, create a temporary markdown file:

```bash
# Create issue body file
cat > /tmp/issue_body.md << 'EOF'
[Issue body markdown content]
EOF

# Create issue using the file
gh issue create \
  --title "Issue Title" \
  --body-file /tmp/issue_body.md \
  --label "label1,label2"

# Clean up
rm /tmp/issue_body.md
```

## Issue Title Guidelines

- **Format**: Use imperative mood ("Add feature" not "Added feature")
- **Length**: 50-72 characters (GitHub UI optimized)
- **Clarity**: Be specific and descriptive
- **Prefixes**: Use prefixes when helpful:
  - `[BTC]` for Bitcoin-specific issues
  - `[ETH]` for Ethereum-specific issues
  - `[XRP]` for XRP-specific issues
  - `[Security]` for security-related issues
  - `[Refactor]` for refactoring tasks

**Examples:**

- `[BTC] Add native SegWit-Bech32 address support`
- `[Security] Implement private key encryption in memory`
- `[Refactor] Migrate wallet service to use case layer`
- `Fix fee calculation overpayment issue on Signet`

## Issue Body Guidelines

### For Feature Requests

- Clearly describe the feature and its use case
- Explain how it fits into the existing architecture
- Specify which layer(s) it affects
- Include acceptance criteria
- Consider security implications (especially for wallet operations)

### For Bug Reports

- Describe the bug clearly
- Include steps to reproduce
- Specify expected vs. actual behavior
- Include error messages or logs (sanitized - no private keys!)
- Specify affected components and wallet types

### For Refactoring Tasks

- Explain what needs to be refactored and why
- Reference related documentation (AGENTS.md, REFACTORING_CHECKLIST.md)
- Specify affected files and components
- Consider backward compatibility
- Plan migration strategy if applicable

### For Security Issues

- **CRITICAL**: Never include sensitive information (private keys, passwords, etc.)
- Describe the security concern clearly
- Specify affected components
- Consider impact on offline wallet operations
- Reference security best practices

## Special Considerations

### Security-Sensitive Issues

For issues involving:

- Private key management
- Wallet operations
- Authentication/authorization
- Encryption/decryption

**Additional requirements:**

- Mark as `security` label
- Set priority appropriately (usually `priority:high` or `priority:critical`)
- Never include sensitive information in issue description
- Consider impact on offline wallets (keygen, sign)
- Reference security guidelines in AGENTS.md

### Architecture-Related Issues

For issues affecting architecture:

- Reference Clean Architecture principles
- Specify layer separation requirements
- Consider dependency direction
- Reference AGENTS.md architecture guidelines
- Consider impact on existing code

### Multi-Chain Support

For issues affecting multiple cryptocurrencies:

- Specify which chains are affected (BTC, BCH, ETH, XRP, ERC-20)
- Consider chain-specific differences
- Test across all affected chains

### Auto-Generated Files

**CRITICAL**: Never create issues to edit files with `DO NOT EDIT` comments:

- SQLC generated files (`internal/infrastructure/database/sqlc/`)
- Protocol buffer generated files
- Files generated by `go generate`

## Output Format

Before creating the issue, show the plan:

1. **Proposed issue title**
2. **Proposed labels**
3. **Issue body preview** (first 10-15 lines)
4. **Affected components**
5. **Priority level**

Ask for confirmation before proceeding with the actual `gh issue create` command.

## Example Workflow

```bash
# 1. Check prerequisites
gh --version
gh auth status

# 2. Review existing issues (optional, to avoid duplicates)
gh issue list --limit 10

# 3. Create issue interactively
gh issue create

# Or create with flags
gh issue create \
  --title "[BTC] Add native SegWit-Bech32 address support" \
  --body "## Description
Add support for native SegWit-Bech32 addresses (bc1...) in Bitcoin wallet operations.

## Context
Currently, the wallet only supports legacy addresses. Native SegWit addresses provide lower transaction fees and better scalability.

## Acceptance Criteria
- [ ] Generate native SegWit addresses for new accounts
- [ ] Support receiving funds to native SegWit addresses
- [ ] Support sending from native SegWit addresses
- [ ] Update tests to cover native SegWit functionality

## Technical Requirements
- **Layer**: Infrastructure (`internal/infrastructure/api/bitcoin/`)
- **Files**: `internal/infrastructure/api/bitcoin/client.go`
- **Dependencies**: None (btcsuite/btcd already supports Bech32)

## Testing Requirements
- Unit tests for address generation
- Integration tests with Bitcoin Core
- Manual testing with testnet

## Related Context
- Related to: #123 (Address format improvements)
- Documentation: AGENTS.md (Architecture guidelines)
" \
  --label "enhancement,btc,priority:medium"
```

## Error Handling

If any step fails:

- Explain what went wrong clearly
- Provide remediation steps
- Don't proceed to subsequent steps until the issue is resolved
- If `gh` command fails, provide alternative: manual issue creation instructions

## Safety Rules

- **CRITICAL**: Never include sensitive information (private keys, passwords, API keys) in issue descriptions
- **CRITICAL**: Verify issue doesn't already exist before creating (use `gh issue list` or search)
- Always review the issue body before finalizing
- Use appropriate labels to help with issue triage
- Consider impact on offline wallet operations (keygen, sign)

## Verification

After creating the issue:

1. Verify the issue was created: `gh issue view {issue_number}`
2. Check that all labels are applied correctly
3. Verify the issue body is formatted correctly
4. Share the issue URL with the user

## Notes

- Issues can be edited after creation using `gh issue edit {issue_number}`
- Use `gh issue list` to view all issues
- Use `gh issue view {issue_number}` to view a specific issue
- Use `gh issue close {issue_number}` to close an issue (not recommended during creation)
