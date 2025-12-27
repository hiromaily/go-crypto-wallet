# Custom Slash Commands

This directory contains definition files for custom slash commands that can be used in Claude Desktop.

## Available Commands

### `commit-push-pr`

Commits changes, pushes them, and creates a pull request.

- Creates a feature branch (from main/master to a working branch)
- Analyzes changes and generates commit messages
- Pushes to remote repository
- Creates a pull request using GitHub CLI

### `create-github-issue`

Creates a GitHub issue.

- Creates a well-structured issue with understanding of project context
- Suggests appropriate labels
- Includes detailed template with security and architecture considerations

### `fix-issue`

Resolves a specified GitHub issue.

- Fetches and analyzes issue content
- Creates a feature branch
- Implements code fixes
- Tests and verifies changes
- Commits and creates a pull request

### `fix-linter`

Fixes linter errors.

- Analyzes errors detected by `make lint-fix`
- Prioritizes and fixes errors
- Uses a step-by-step approach to fixes

### `fix-pr-review`

Addresses pull request review comments.

- Fetches PR information and review comments
- Categorizes and prioritizes comments
- Implements fixes and tests
- Commits and pushes changes

## Usage

In Claude Desktop, you can use these commands as slash commands. For example:

- `/commit-push-pr` - Commit changes and create a pull request
- `/create-github-issue` - Create a new GitHub issue
- `/fix-issue #123` - Fix issue #123
- `/fix-linter` - Fix linter errors
- `/fix-pr-review #456` - Address review comments for PR #456

Each command operates in accordance with the project guidelines (`AGENTS.md`), following Clean Architecture
principles and security best practices.
