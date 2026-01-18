# Required Tools and Versions

This document lists all required tools and their versions for the go-crypto-wallet project.
Always verify tool versions before starting development work.

## Overview

Before starting any development task, ensure all required tools are installed with the correct versions.
Using incorrect versions may cause compatibility issues or unexpected behavior.

## Essential Tools

### Git

- **Required**: Yes
- **Check version**: `git --version`
- **Installation**: <https://git-scm.com/>
- **Note**: Any recent version should work

### GitHub CLI (gh)

- **Required**: Yes (for custom slash commands: `fix-issue`, `fix-pr-review`, `create-github-issue`)
- **Check version**: `gh --version`
- **Installation**: <https://cli.github.com/>
- **Authentication**: Verify with `gh auth status`, login with `gh auth login` if needed
- **Note**: Required for creating issues, viewing PRs, and managing GitHub resources

### Go

- **Required Version**: Go 1.25.5 (specified in `go.mod`)
- **Check version**: `go version`
- **Installation**: <https://go.dev/dl/>
- **Note**: Must match the version specified in `go.mod` (currently `go 1.25.5`)

## Development Tools

### golangci-lint

- **Required Version**: v2.8.0 (specified in `make/vars.mk`)
- **Check version**: `golangci-lint --version`
- **Installation**: See [Coding Standards](coding-standards.md)
- **Note**: Used for code linting and formatting

### Atlas

- **Required Version**: v1.0.0 (specified in `tools/atlas/atlas.hcl`)
- **Check version**: `atlas version`
- **Installation**:

  ```bash
  brew install arigaio/tap/atlas
  # or
  go install ariga.io/atlas/cmd/atlas@v1.0.0
  ```

- **Note**: **CRITICAL** - This project uses Atlas v1.0.0.
  Using older versions may not support required features.
  Always verify version before running Atlas commands.
- **Usage**: Database schema management and migrations (see [Database Management](database.md))

### markdownlint-cli

- **Required**: Optional (for Markdown file linting)
- **Check version**: `markdownlint --version` or `npx markdownlint-cli --version`
- **Installation**:

  ```bash
  npm install -g markdownlint-cli
  # or use npx (no installation needed)
  npx markdownlint-cli
  ```

- **Note**: Used for linting Markdown files (`.md`)

### Docker

- **Required**: Yes (for database and blockchain node containers)
- **Check version**: `docker --version`
- **Installation**: <https://docs.docker.com/get-docker/>
- **Note**: Required for running MySQL, Bitcoin Core, Ethereum nodes, etc.

### Docker Compose

- **Required**: Yes (Docker Compose v2 or later)
- **Check version**: `docker compose version`
- **Installation**: Included with Docker Desktop, or install separately from <https://docs.docker.com/compose/install/>
- **Note**: **CRITICAL** - This project uses `docker compose wait` command which requires Docker Compose v2 or later.
  The old `docker-compose` (v1, hyphenated) does not support the `wait` command.
  Ensure you're using `docker compose` (space, v2) not `docker-compose` (hyphen, v1).
  Requires Docker Compose v2.21 or later (see `make/db.mk` for details).
- **Usage**: Database migrations and service orchestration (see `make/db.mk`)

## Version Verification

Before starting work, verify all required tools are installed and at the correct versions:

```bash
# Essential tools
git --version
gh --version
go version  # Should show go1.25.5

# Development tools
golangci-lint --version  # Should show v2.8.0
atlas version             # Should show v1.0.0

# Optional tools
markdownlint --version    # If using markdownlint
docker --version
docker compose version    # Should show v2.x.x or later
```

**If any required tool is missing or at an incorrect version:**

- Stop immediately
- Install or update the tool to the correct version
- Verify installation before proceeding
- Do not proceed with the workflow until all requirements are met

## See Also

- [Workflow Guidelines](workflow.md) - Development workflow and common steps
- [Database Management](database.md) - Atlas usage for database migrations
- [Coding Standards](coding-standards.md) - golangci-lint configuration and usage
