# Claude Code Commands

Minimal commands that invoke Skills for common tasks.

## Available Commands

| Command | Description |
|---------|-------------|
| `/fix-issue #123` | Fix a GitHub issue |
| `/fix-linter` | Fix linter errors |
| `/fix-pr-review #123` | Address PR review comments |

## Skills (Referenced by Commands)

Commands use these Skills for detailed workflows:

| Skill | Description |
|-------|-------------|
| `go-development` | Branch, verify, review, commit workflow |
| `github-issue-creation` | Create GitHub issues |

## Usage

1. Invoke command (e.g., `/fix-issue #123`)
2. Command loads the `go-development` Skill automatically
3. Follow the workflow defined in the Skill
