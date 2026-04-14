# Agent Skills Documentation

## Overview

Agent Skills are reusable, discoverable capabilities that AI agents can use to perform specialized tasks. This repository includes Agent Skills in the `.claude/skills/` directory that are automatically discovered by compatible AI coding assistants.

## What are Agent Skills?

Agent Skills are a structured format for defining AI agent capabilities:

- **Discoverable**: AI assistants automatically find skills in `.claude/skills/`
- **Team-wide**: Skills committed to the repository are available to all team members
- **Structured**: YAML frontmatter provides metadata for skill discovery
- **Progressive**: Only relevant content is loaded when needed

### Skills vs Custom Commands

| Feature | Agent Skills (`.claude/skills/`) | Custom Commands (`.claude/commands/`) |
|---------|----------------------------------|---------------------------------------|
| **Discovery** | Automatic by AI assistant | Manual invocation only |
| **Format** | YAML frontmatter + Markdown | Plain Markdown |
| **Sharing** | Repository-wide, team sharing | Local to user |
| **Loading** | Progressive (on-demand) | Full content loaded |
| **Organization** | One skill per directory | One file per command |

## Available Skills

### btc-terminology

**Location**: `.claude/skills/btc-terminology/SKILL.md`

**Description**: Critical Bitcoin terminology rules to prevent confusion between bech32m (encoding) and taproot (address type). Use when working on BTC-related code, config files, or shell scripts.

**Key Concept**:

- **bech32m** = Encoding format (HOW address is serialized) → Use in Bitcoin Core RPC/CLI
- **taproot** = Address type (WHAT the address represents) → Use in config files

**When to use**:

- Working on Taproot-related code (Pattern 9, 10, 11)
- Modifying BTC config files or shell scripts
- Understanding why `bech32m` vs `taproot` is used in different contexts

### github-issue-creation

**Location**: `.claude/skills/github-issue-creation/SKILL.md`

**Description**: Create well-structured GitHub issues using the gh CLI tool. Use when you want to create a GitHub issue to track features, bugs, refactoring tasks, or technical debt.

**Capabilities**:

- Analyzes user requests and creates comprehensive issue proposals
- Suggests appropriate labels from repository configuration
- Follows project-specific guidelines (Clean Architecture, security requirements)
- Provides two-step workflow: propose → review → submit
- Includes detailed templates for different issue types (feature, bug, refactoring, security)

**When to use**:

- Creating feature requests
- Reporting bugs
- Proposing refactoring tasks
- Documenting technical debt
- Creating security-related issues

See [Usage Examples](#usage-examples) below for detailed examples.

## Installation

### For Claude Desktop / Claude Code

Agent Skills are **automatically discovered** by Claude Desktop and Claude Code - no installation needed!

#### Requirements

1. **Claude Desktop** (3.x or later) or **Claude Code** CLI
2. **Repository cloned locally** with `.claude/skills/` directory

#### Verification

To verify skills are discovered:

1. Open your repository in Claude Desktop or Claude Code
2. Skills in `.claude/skills/` are automatically available
3. Claude will use skills when appropriate based on your requests

**Note**: No manual configuration or installation is required. Skills are loaded automatically from the repository.

### For Cursor

Cursor supports Claude skills through repository context:

#### Setup Steps

1. **Open repository** in Cursor
2. **Ensure `.claude/skills/` directory** is part of the workspace
3. **Skills are available** through Cursor's AI features

**Note**: Cursor reads the `.claude/` directory structure and makes skills available to its AI models.

### For Windsurf (Codeium)

Windsurf supports project-level AI configurations:

#### Setup Steps

1. **Clone repository** with `.claude/skills/` directory
2. **Open in Windsurf**
3. Skills should be **automatically discovered** via the `.claude/` directory

**Note**: Windsurf's AI agent reads repository structure including `.claude/` configurations.

### For Other AI Assistants

For AI coding assistants that support custom prompts or context:

1. **Manual Reference**: Copy skill content from `.claude/skills/{skill-name}/SKILL.md`
2. **Custom Prompts**: Use skill content as custom prompt templates
3. **Repository Context**: Ensure the AI assistant can read the `.claude/` directory

## How to Use Agent Skills

### Using with Claude Code (CLI)

When using Claude Code CLI, skills are invoked automatically based on your natural language requests.

#### Example: Creating a GitHub Issue

```bash
# Start Claude Code session
claude

# In the conversation, request issue creation
You: "I want to create a GitHub issue for adding XRP transaction batching support"

# Claude will:
# 1. Automatically detect this matches the github-issue-creation skill
# 2. Use the skill to analyze your request
# 3. Propose a comprehensive issue with proper structure
# 4. Wait for your approval before creating the issue
```

### Using with Claude Desktop

In Claude Desktop conversations:

1. **Natural Language Request**: Simply describe what you want to do
2. **Skill Auto-Discovery**: Claude automatically uses appropriate skills
3. **Interactive Workflow**: Follow the skill's prompts and workflows

#### Example Conversation

```
You: "Create an issue for implementing native SegWit support"

Claude: [Uses github-issue-creation skill]
## Proposed Issue

**Title**: [BTC] Add native SegWit-Bech32 address support

**Labels**: enhancement, btc

**Issue Body**:
[... comprehensive issue proposal ...]

**Ready to submit?** Please review and confirm.

You: "Yes, create it"

Claude: [Creates the issue using gh CLI]
✓ Issue created: #123
✓ URL: https://github.com/hiromaily/go-crypto-wallet/issues/123
```

### Explicit Skill Invocation (Advanced)

For explicit control, you can reference skills directly:

```
You: "Use the github-issue-creation skill to create an issue for database migration"

Claude: [Explicitly invokes the specified skill]
```

## Usage Examples

### Example 1: Creating a Feature Request

**User Request**:

```
I want to track a new feature for adding hardware wallet support using Ledger devices
```

**Claude's Response** (using github-issue-creation skill):

```markdown
## Proposed Issue

**Title**: Add Ledger hardware wallet integration for offline signing

**Labels**: enhancement, security

**Issue Body**:
## Description
Integrate Ledger hardware wallet support for the sign wallet to enable secure offline transaction signing using hardware devices.

## Context
Hardware wallets provide an additional layer of security by keeping private keys in secure hardware elements, preventing key exposure even if the host system is compromised.

## Acceptance Criteria
- [ ] Detect connected Ledger devices
- [ ] Sign Bitcoin transactions using Ledger
- [ ] Sign Ethereum transactions using Ledger
- [ ] Support BIP-44/BIP-49/BIP-84 derivation paths
- [ ] Handle device disconnection gracefully

## Technical Requirements

### Architecture
- Layer: Infrastructure (device communication)
- Integration point: Sign wallet signing flow
- Dependencies: ledger-go library

### Implementation Details
- **Files to create**: `internal/infrastructure/hardware/ledger/client.go`
- **Files to modify**: `internal/application/usecase/sign/*/sign.go`
- **Dependencies**: github.com/btcsuite/btcd/btcec, Ledger USB/HID libraries

### Constraints
- Must work with offline sign wallet (no network access)
- Must support multiple address formats
- Must handle device-specific errors clearly

## Testing Requirements
- Unit tests for Ledger client interface
- Integration tests with Ledger simulator
- Manual testing with physical Ledger device

## Related Context
- Related to: #45 (Offline signing improvements)
- Documentation: AGENTS.md (Security guidelines)

**Affected Components**: Infrastructure layer (hardware device communication)

**Priority**: High

---

**Ready to submit?** Please review the proposal above and confirm if you'd like me to create this issue.
```

**User confirms**: "Yes, create it"

**Claude creates the issue** and provides the issue number and URL.

### Example 2: Reporting a Bug

**User Request**:

```
I found a bug where the fee calculation is overpaying on Bitcoin testnet transactions
```

**Claude's Response**:

```markdown
## Proposed Issue

**Title**: Fix fee calculation overpayment issue on Bitcoin testnet

**Labels**: bug, btc

**Issue Body**:
## Description
Fee calculation is resulting in overpayment for Bitcoin testnet transactions.

## Context
Users are observing transaction fees significantly higher than expected when sending testnet BTC.

## Steps to Reproduce
1. Create a transaction on Bitcoin testnet
2. Calculate fees using current fee rate
3. Compare actual fee paid vs expected fee

## Expected Behavior
Transaction fees should match the calculated fee based on transaction size and fee rate.

## Actual Behavior
Transaction fees are 2-3x higher than expected calculation.

## Affected Components
- Watch wallet (fee calculation)
- Bitcoin infrastructure layer

## Technical Requirements

### Architecture
- Layer: Application (use case) and Infrastructure (Bitcoin client)
- Files to investigate: `internal/application/usecase/watch/btc/estimate_fee.go`

### Root Cause Investigation Needed
- Check fee rate source (Bitcoin Core RPC)
- Verify transaction size calculation
- Review UTXO selection logic

## Testing Requirements
- Unit tests for fee calculation
- Integration tests on testnet
- Manual verification with different transaction sizes

## Additional Notes
- May be related to Taproot transaction size estimation
- Check if issue exists on mainnet

**Affected Components**:
- Application layer (use case)
- Infrastructure layer (Bitcoin client)

**Priority**: High

---

**Ready to submit?** Please review the proposal above and confirm if you'd like me to create this issue.
```

### Example 3: Refactoring Task

**User Request**:

```
We need to refactor the XRP wallet to use the new DTO pattern like we did for Bitcoin
```

**Claude's Response**:

```markdown
## Proposed Issue

**Title**: [Refactor] Migrate XRP wallet to use application-layer DTOs

**Labels**: refactoring, xrp

**Issue Body**:
## Description
Refactor XRP wallet implementation to use application-layer Data Transfer Objects (DTOs) instead of infrastructure-layer types, following the pattern established for Bitcoin in #242.

## Context
To maintain Clean Architecture principles, wallet types should be defined at the application layer, with infrastructure types converted to application DTOs at the boundary.

## Acceptance Criteria
- [ ] Create XRP DTOs in `internal/application/usecase/watch/xrp/`
- [ ] Update XRP wallet adapters to convert infrastructure types to DTOs
- [ ] Remove direct usage of infrastructure types in use cases
- [ ] Update tests to use application-layer types
- [ ] Verify no infrastructure dependencies in domain/application layers

## Technical Requirements

### Architecture
- Layer: Application and Infrastructure boundary
- Pattern: Infrastructure → Application DTO conversion at adapter layer
- Reference: Bitcoin implementation (#242)

### Implementation Details
- **Files to create**: `internal/application/usecase/watch/xrp/dto.go`
- **Files to modify**:
  - `internal/application/usecase/watch/xrp/*.go`
  - `internal/interface-adapters/wallet/xrp/adapter.go`
- **Reference implementation**: `internal/application/usecase/watch/btc/dto.go`

### Constraints
- Must maintain existing functionality (no behavior changes)
- Must follow Clean Architecture dependency rules
- Must be backward compatible with existing tests

## Testing Requirements
- Unit tests for DTO conversion functions
- Integration tests to verify functionality unchanged
- Verify no infrastructure imports in application layer

## Related Context
- Related to: #242 (Bitcoin DTO migration)
- Related to: #224 (Infrastructure refactoring)
- Documentation: docs/guidelines/architecture.md (Clean Architecture guidelines)

**Affected Components**:
- Application layer (use cases)
- Infrastructure layer (adapter conversion)

**Priority**: Medium

---

**Ready to submit?** Please review the proposal above and confirm if you'd like me to create this issue.
```

## Troubleshooting

### Skill Not Found / Not Discovered

**Problem**: AI assistant doesn't seem to recognize skills

**Solutions**:

1. **Verify file location**: Skills must be in `.claude/skills/{skill-name}/SKILL.md`
2. **Check YAML frontmatter**: Ensure proper format with `name` and `description` fields
3. **Restart AI assistant**: Close and reopen Claude Desktop/Code
4. **Verify repository context**: Ensure `.claude/` directory is in the workspace

### Skill Not Triggered Automatically

**Problem**: Claude doesn't use the skill when expected

**Solutions**:

1. **Be explicit**: Mention the skill name or capability directly
2. **Check description**: Skill description determines when it's triggered
3. **Manual invocation**: Explicitly ask Claude to use a specific skill

### gh CLI Commands Fail

**Problem**: GitHub CLI commands in skills fail

**Solutions**:

1. **Install gh**: `brew install gh` (macOS) or see [GitHub CLI installation](https://cli.github.com/)
2. **Authenticate**: Run `gh auth login`
3. **Verify access**: Run `gh auth status`
4. **Check repository**: Ensure you're in the correct repository directory

### Label Sync Issues

**Problem**: Label creation/sync commands fail

**Solutions**:

1. **Check permissions**: Ensure you have write access to the repository
2. **Verify labels.yml**: Confirm `.github/labels.yml` exists and is valid YAML
3. **Use correct method**:
   - Method 1 (bash script): Works with standard bash
   - Method 2 (yq): Requires `yq` to be installed (`brew install yq`)
4. **Manual creation**: Create labels manually using `gh label create` if sync fails

### Permission Errors

**Problem**: Cannot create issues or labels

**Solutions**:

1. **Check GitHub permissions**: Verify you have write access to the repository
2. **Re-authenticate gh**: Run `gh auth logout` then `gh auth login`
3. **Check repository**: Ensure you're in the correct repository: `git remote -v`

## Best Practices

### For Users

1. **Be Clear**: Provide clear, specific requests for what you want
2. **Review Proposals**: Always review AI-generated issue proposals before submission
3. **Provide Context**: Include relevant details about features, bugs, or refactoring needs
4. **Follow Workflow**: Let the skill guide you through its workflow

### For Skill Authors

1. **Clear Descriptions**: Write descriptive YAML frontmatter (include "what" and "when")
2. **Proper Structure**: Follow the `.claude/skills/{skill-name}/SKILL.md` pattern
3. **Comprehensive Content**: Include prerequisites, examples, and error handling
4. **Test Commands**: Verify all shell commands work before committing
5. **Documentation**: Keep skills documented and up-to-date

## Additional Resources

- [Agent Skills Overview (Anthropic)](https://platform.claude.com/docs/en/agents-and-tools/agent-skills/overview)
- [GitHub CLI Documentation](https://cli.github.com/manual/)
- `AGENTS.md` - Project-specific AI agent guidelines
- [Workflow Guidelines](./guidelines/workflow.md) - Common workflow steps

## Contributing

To add new Agent Skills to this repository:

1. **Create skill directory**: `.claude/skills/{skill-name}/`
2. **Create SKILL.md**: With proper YAML frontmatter
3. **Follow structure**: Use existing skills as templates
4. **Test thoroughly**: Verify all commands work
5. **Document**: Add entry to this document's "Available Skills" section
6. **Submit PR**: Create pull request for review

## Support

If you encounter issues with Agent Skills:

1. Check [Troubleshooting](#troubleshooting) section above
2. Review skill-specific documentation in `.claude/skills/{skill-name}/SKILL.md`
3. Check project guidelines in `AGENTS.md`
4. Create a GitHub issue for repository-specific problems
