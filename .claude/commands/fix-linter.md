# Fix Linter

Fix linter errors by selecting appropriate skill based on language.

## Task Classification

| Linter Command | Language | Skill |
|----------------|----------|-------|
| `make go-lint` | Go | `go-development` |
| `npm run lint` (in apps/) | TypeScript/JS | `typescript-development` |
| `npm run lint` (contracts) | Solidity | `solidity-development` |

## Process

1. **Identify language** from linter command or error messages
2. **Load Skill**: Use appropriate `{lang}-development` Skill
3. **Prioritize**: syntax > security > type > style
4. **Fix**: Address errors by priority
5. **Verify**: Run Skill-specific verification commands

## Quick Reference

### Go

```bash
make go-lint
# Fix errors
make go-lint && make tidy && make check-build && make go-test
```

### TypeScript

```bash
cd apps/xrpl-grpc-server
bun run lint
# Fix errors
bun run lint && bun run build && bun test
```

### Solidity

```bash
cd apps/erc20-token
npm run lint
# Fix errors
truffle compile && truffle test
```
