## Coding Conventions

This document describes the coding standards and conventions for the go-crypto-wallet project.

### Language-Specific Rules

For detailed rules, format commands, and verification commands, see the corresponding rule files in `.claude/rules/`:

| Language | Rule File | Key Commands |
|----------|-----------|--------------|
| Go | `.claude/rules/go/` | `make go-lint`, `make check-build` |
| TypeScript/JS | `.claude/rules/typescript.md` | `yarn lint`, `npm run lint` |
| Shell | `.claude/rules/shell-script.md` | `make shfmt`, `shellcheck` |
| SQL | `.claude/rules/sql.md` | `make sqlc-validate`, `make sqlc` |
| HCL | `.claude/rules/hcl.md` | `make atlas-fmt`, `make atlas-lint` |
| Proto | `.claude/rules/proto.md` | `make proto-fmt`, `make proto` |
| YAML | `.claude/rules/yaml.md` | `make yaml-lint` |
