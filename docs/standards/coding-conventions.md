# Coding Conventions

Project coding standards for go-crypto-wallet.

## Language-Specific Rules (SSOT)

For detailed rules, format commands, and verification commands, see the corresponding rule files in `.claude/rules/`:

| Language | Rule File | Key Commands |
|----------|-----------|--------------|
| Go | [.claude/rules/go.md](../../.claude/rules/go.md) | `make go-lint`, `make check-build` |
| TypeScript/JS | [.claude/rules/typescript.md](../../.claude/rules/typescript.md) | `yarn lint`, `npm run lint` |
| Shell | [.claude/rules/shell-script.md](../../.claude/rules/shell-script.md) | `make shfmt`, `shellcheck` |
| SQL | [.claude/rules/sql.md](../../.claude/rules/sql.md) | `make sqlc-validate`, `make sqlc` |
| HCL | [.claude/rules/hcl.md](../../.claude/rules/hcl.md) | `make atlas-fmt`, `make atlas-lint` |
| Proto | [.claude/rules/proto.md](../../.claude/rules/proto.md) | `make proto-fmt`, `make proto` |
| YAML | [.claude/rules/yaml.md](../../.claude/rules/yaml.md) | `make yaml-lint` |

## Quick Verification Reference

```bash
# Go files
make go-lint && make tidy && make check-build && make gotest

# Database schema (HCL)
make atlas-fmt && make atlas-lint

# SQL queries
make sqlc-validate && make sqlc
```

## Detailed Guidelines

See [docs/guidelines/coding-standards.md](../guidelines/coding-standards.md) for additional guidelines.
