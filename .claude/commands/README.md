# Custom Slash Commands

This directory contains definition files for custom slash commands that can be used in Claude Desktop.

## Architecture

These commands follow a modular structure where:
- **Common workflow steps** are defined in [Workflow Guidelines](../../agents/workflow.md)
- **Command-specific logic** is defined in each command file
- This reduces duplication and makes maintenance easier

## Available Commands

### `create-github-issue`

Creates a GitHub issue.

- Creates a well-structured issue with understanding of project context
- Suggests appropriate labels
- Includes detailed template with security and architecture considerations
- See [create-github-issue.md](create-github-issue.md) for details

### `fix-issue`

Resolves a specified GitHub issue.

- Fetches and analyzes issue content
- Creates a feature branch
- Implements code fixes
- Tests and verifies changes
- Commits and creates a pull request
- See [fix-issue.md](fix-issue.md) for details

### `fix-linter`

Fixes linter errors.

- Analyzes errors detected by `make lint-fix`
- Prioritizes and fixes errors
- Uses a step-by-step approach to fixes
- See [fix-linter.md](fix-linter.md) for details

### `fix-pr-review`

Addresses pull request review comments.

- Fetches PR information and review comments
- Categorizes and prioritizes comments
- Implements fixes and tests
- Commits and pushes changes
- See [fix-pr-review.md](fix-pr-review.md) for details

## Usage

In Claude Desktop, you can use these commands as slash commands. For example:

- `/create-github-issue` - Create a new GitHub issue
- `/fix-issue #123` - Fix issue #123
- `/fix-linter` - Fix linter errors
- `/fix-pr-review #456` - Address review comments for PR #456

## Common Workflow Steps

All commands follow common workflow steps defined in [Workflow Guidelines](../../agents/workflow.md):

- **Required Tools and Versions**: See [Required Tools and Versions](../../agents/requirements.md) - Essential tools (Git, GitHub CLI, Go 1.25.5) and development tools (Atlas v1.0.0, golangci-lint v2.7.2) with version requirements
- **Pre-Flight Checks**: Git status, branch verification, project guidelines review
- **Safety Rules**: Critical rules for security and code quality
- **Verification Steps**: Commands to run before committing (`make lint-fix`, `make tidy`, `make check-build`, `make gotest`)
- **Special Considerations**: Security-sensitive changes, auto-generated files, breaking changes

Each command file links to these common steps to avoid duplication and ensure consistency.

**Important**: Always verify tool versions before starting work. Using incorrect versions (especially Atlas v1.0.0) may cause compatibility issues.

## Project Guidelines

Each command operates in accordance with the project guidelines ([AGENTS.md](../../AGENTS.md)), following:
- Clean Architecture principles
- Security best practices
- Code quality standards
- Testing requirements
