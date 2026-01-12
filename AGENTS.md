# Claude Code Configuration

This document contains **Claude Code specific settings** for the go-crypto-wallet project.
For general agent guidelines, see [AGENTS.md](AGENTS.md).

## General Guidelines

Refer to [AGENTS.md](AGENTS.md) for:

- Core values and priorities
- Expected behavior (do/don't/ask)
- Documentation map
- Verification commands

## Claude Code Preferences

### Tool Usage

- Prefer built-in file operations over shell commands for reading/writing files
- Use semantic search (`SemanticSearch`) for exploring unfamiliar code
- Use `Grep` for exact text/symbol searches
- Batch parallel tool calls when operations are independent

### Code Changes

- Always read files before proposing edits
- Provide context with code changes (before/after)
- Run verification commands after changes:
  - Go: `make go-lint && make check-build`
  - TypeScript: `npm run lint && npm run build`

### Response Style

- Be concise but thorough
- Show relevant code snippets with file paths
- Explain "why" not just "what"
- Use Japanese when the user writes in Japanese

## Project-Specific Notes

### Security-Critical Areas

These areas require extra caution:

- `internal/infrastructure/wallet/key/` - Key generation
- `internal/domain/key/` - Key value objects
- Any code handling private keys or seeds

### Auto-Generated Files (DO NOT EDIT)

- `internal/infrastructure/database/sqlc/*.go`
- `internal/infrastructure/api/ripple/xrp/*.pb.go`
- `internal/infrastructure/contract/token-abi.go`

### Build Tags

- Integration tests use `//go:build integration`
- Run with: `go test -tags=integration ./...`

## See Also

- [AGENTS.md](AGENTS.md) - General agent guidelines
- [llms.txt](llms.txt) - Project sitemap
- [ARCHITECTURE.md](ARCHITECTURE.md) - System architecture
