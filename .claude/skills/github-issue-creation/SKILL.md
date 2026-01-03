---
name: github-issue-creation
description: Create well-structured GitHub issues using the gh CLI tool. Use when the user wants to create a GitHub issue, track features, bugs, refactoring tasks, or technical debt for the go-crypto-wallet repository.
---

# Create GitHub Issue

You are tasked with creating a GitHub issue using the `gh` command-line tool.
This command helps create well-structured GitHub issues for tracking features, bugs, refactoring tasks, and technical debt.

## Repository

Repo: hiromaily/go-crypto-wallet

## Prerequisites Check

Before starting, verify all required tools are installed with the correct versions.

See [Required Tools and Versions](../../../agents/requirements.md) for:

- Complete list of required tools
- Version requirements (Go 1.25.5, Atlas v1.0.0, golangci-lint v2.7.2, etc.)
- Installation instructions
- Version verification commands

**For this command specifically:**

- **GitHub CLI (gh)**: Required - Check with `gh --version`
  - Verify authentication: `gh auth status`
  - If not authenticated: Run `gh auth login`
- **Git**: Required - Check with `git --version`

**If any required tool is missing or at an incorrect version, stop and display an error message with installation instructions.
Do not proceed with the workflow.**

### Label Management

Before creating issues, ensure all required labels exist in the repository.

**Check existing labels:**

```bash
gh label list
```

**Sync labels from `.github/labels.yml`:**

```bash
# Parse labels.yml and create each label
# Note: This requires yq (YAML parser) or manual parsing
while IFS= read -r line; do
  if [[ $line =~ ^-\ name:\ (.+)$ ]]; then
    name="${BASH_REMATCH[1]}"
  elif [[ $line =~ ^\ \ color:\ (.+)$ ]]; then
    color="${BASH_REMATCH[1]}"
  elif [[ $line =~ ^\ \ description:\ (.+)$ ]]; then
    desc="${BASH_REMATCH[1]}"
    # Create or update label
    gh label create "$name" --color "$color" --description "$desc" --force 2>/dev/null || echo "Label $name already exists or error occurred"
  fi
done < .github/labels.yml
```

**Alternative using yq (if installed):**

```bash
# Install yq if needed: brew install yq (macOS) or see https://github.com/mikefarah/yq
yq eval '.[] | "gh label create \"" + .name + "\" --color " + .color + " --description \"" + .description + "\" --force"' .github/labels.yml | sh
```

**Note**: If labels are missing, run one of the sync commands above to create them from `.github/labels.yml`.
The repository uses `.github/labels.yml` as the source of truth for label definitions.
The `--force` flag updates existing labels with new colors/descriptions.

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

### Step 1: Analyze and Propose Issue

Before creating the issue, analyze the user's request and create a comprehensive issue proposal.

#### 1.1 Gather Information

From the user's request or conversation context, gather:

- **Issue Type**: Feature request, bug report, refactoring task, documentation, security, technical debt
- **Title**: Clear, concise description (50-72 characters recommended)
- **Description**: Detailed explanation of what needs to be implemented or fixed
- **Priority**: Critical, High, Medium, Low (based on impact and urgency)
- **Affected Components**: Which layers/components are affected (domain, application, infrastructure, interface-adapters)
- **Related Context**: Related issues, PRs, or documentation references

#### 1.2 Determine Issue Category

Based on the issue type, determine appropriate labels and structure:

- **Feature Request**: New functionality or enhancement
- **Bug Report**: Something that's broken or not working as expected
- **Refactoring**: Code improvement without changing functionality
- **Documentation**: Documentation updates or improvements
- **Security**: Security-related issues or improvements
- **Technical Debt**: Code quality improvements, cleanup tasks

#### 1.3 Create Issue Proposal

Create a well-structured issue proposal following this template:

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

#### 1.4 Determine Labels

Based on the issue type and content, suggest appropriate labels from `.github/labels.yml`:

**Available labels (from `.github/labels.yml`):**

- **Type labels**: `bug`, `enhancement`, `refactoring`, `documentation`, `security`, `technical-debt`
- **Status labels**: `duplicate`, `invalid`, `question`, `wontfix`, `help wanted`, `good first issue`
- **Dependency labels**: `dependencies`

**Additional labels that may exist:**

- **Priority labels**: `priority:critical`, `priority:high`, `priority:medium`, `priority:low` (if defined)
- **Component labels**: `domain`, `application`, `infrastructure`, `interface-adapters`, `btc`, `eth`, `xrp`, `bch`
  (if defined)

#### 1.5 Present Issue Proposal

Present the complete issue proposal to the user in the following format:

```markdown
## Proposed Issue

**Title**: [Issue title]

**Labels**: [Comma-separated list of labels]

**Issue Body**:
[Complete issue body following the template below]

**Affected Components**: [List of affected components/layers]

**Priority**: [Priority level]

**Related Context**: [Related issues, PRs, or documentation]
```

**Wait for user approval before proceeding to Step 2.**

### Step 2: Submit Issue (After User Approval)

After the user approves the issue proposal, proceed with creating the issue.

#### 2.1 Verify Labels

Before creating the issue:

1. Check available labels: `gh label list`
2. If labels are missing, sync from `.github/labels.yml` using one of these methods:

   **Method 1: Using bash script (no additional dependencies)**
   ```bash
   while IFS= read -r line; do
     if [[ $line =~ ^-\ name:\ (.+)$ ]]; then
       name="${BASH_REMATCH[1]}"
     elif [[ $line =~ ^\ \ color:\ (.+)$ ]]; then
       color="${BASH_REMATCH[1]}"
     elif [[ $line =~ ^\ \ description:\ (.+)$ ]]; then
       desc="${BASH_REMATCH[1]}"
       gh label create "$name" --color "$color" --description "$desc" --force 2>/dev/null || echo "Label $name already exists"
     fi
   done < .github/labels.yml
   ```

   **Method 2: Using yq (if installed)**
   ```bash
   yq eval '.[] | "gh label create \"" + .name + "\" --color " + .color + " --description \"" + .description + "\" --force"' .github/labels.yml | sh
   ```

3. Verify the labels you plan to use exist in the repository

#### 2.2 Create the Issue

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

#### 2.3 Alternative: Create Issue from File

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

#### 2.4 Verify Issue Creation

After creating the issue:

1. Verify the issue was created: `gh issue view {issue_number}`
2. Check that all labels are applied correctly
3. Verify the issue body is formatted correctly
4. Share the issue URL with the user

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
- See [Security-Sensitive Changes](../../../agents/workflow.md#security-sensitive-changes) in Workflow Guidelines

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

**CRITICAL**: Never create issues to edit files with `DO NOT EDIT` comments.

See [Auto-Generated Files](../../../agents/workflow.md#auto-generated-files) in Workflow Guidelines for details.

## Output Format

### Step 1 Output: Issue Proposal

When analyzing and proposing an issue, present it in this format:

```markdown
## Proposed Issue

**Title**: [Issue title - 50-72 characters]

**Labels**: [Comma-separated list of labels from .github/labels.yml]

**Issue Body**:
[Complete issue body following the template structure]

**Affected Components**:
- [Component 1]
- [Component 2]

**Priority**: [Critical/High/Medium/Low]

**Related Context**:
- Related issues: #[number]
- Related PRs: #[number]
- Documentation: [links]

---

**Ready to submit?** Please review the proposal above and confirm if you'd like me to create this issue.
```

**Important**: Do NOT proceed to Step 2 (creating the issue) until the user explicitly approves the proposal.

### Step 2 Output: Issue Creation

After user approval, show the creation process:

```markdown
Creating issue...

✓ Labels verified
✓ Issue created: #{issue_number}
✓ URL: https://github.com/hiromaily/go-crypto-wallet/issues/{issue_number}
```

## Example Workflow

### Step 1: Analyze and Propose

```bash
# User provides context about the issue they want to create
# AI analyzes and creates a proposal
```

**AI presents proposal:**

```markdown
## Proposed Issue

**Title**: [BTC] Add native SegWit-Bech32 address support

**Labels**: enhancement,btc

**Issue Body**:
## Description
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

**Affected Components**:
- Infrastructure layer (Bitcoin API client)

**Priority**: Medium

---

**Ready to submit?** Please review the proposal above and confirm if you'd like me to create this issue.
```

### Step 2: Submit (After User Approval)

```bash
# 1. Verify labels exist
gh label list

# 2. Sync labels if needed (use bash parsing method)
while IFS= read -r line; do
  if [[ $line =~ ^-\ name:\ (.+)$ ]]; then
    name="${BASH_REMATCH[1]}"
  elif [[ $line =~ ^\ \ color:\ (.+)$ ]]; then
    color="${BASH_REMATCH[1]}"
  elif [[ $line =~ ^\ \ description:\ (.+)$ ]]; then
    desc="${BASH_REMATCH[1]}"
    gh label create "$name" --color "$color" --description "$desc" --force 2>/dev/null
  fi
done < .github/labels.yml

# 3. Create the issue
gh issue create \
  --title "[BTC] Add native SegWit-Bech32 address support" \
  --body-file /tmp/issue_body.md \
  --label "enhancement,btc"

# 4. Verify creation
gh issue view {issue_number}
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

See also [Safety Rules](../../../agents/workflow.md#safety-rules) in Workflow Guidelines for general safety rules.

## Process Summary

1. **Step 1: Analyze and Propose**
   - Gather information from user request
   - Analyze issue type and requirements
   - Create comprehensive issue proposal
   - Present proposal to user for review
   - **Wait for user approval**

2. **Step 2: Submit Issue (After Approval)**
   - Verify labels exist in repository
   - Sync labels from `.github/labels.yml` if needed
   - Create issue using `gh issue create`
   - Verify issue was created correctly
   - Share issue URL with user

## Notes

- Issues can be edited after creation using `gh issue edit {issue_number}`
- Use `gh issue list` to view all issues
- Use `gh issue view {issue_number}` to view a specific issue
- Use `gh issue close {issue_number}` to close an issue (not recommended during creation)
