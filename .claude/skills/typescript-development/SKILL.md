---
name: typescript-development
description: TypeScript/JavaScript development workflow for apps/ directory. Use when modifying TypeScript code in ripple-lib-server or JavaScript in erc20-token.
---

# TypeScript/JavaScript Development Workflow

Workflow for TypeScript/JavaScript code changes in `apps/` directory.

## Prerequisites

- **Use `git-workflow` Skill** for branch management, commit conventions, and PR creation.
- **Refer to `.claude/rules/typescript.md`** for detailed verification commands (SSOT).

## Applicable Directories

| App | Language | Path |
|-----|----------|------|
| ripple-lib-server | TypeScript | `apps/ripple-lib-server/` |
| erc20-token | JavaScript | `apps/erc20-token/` |

## Workflow

### 1. Make Changes

Edit TypeScript/JavaScript files following the rules in `.claude/rules/typescript.md`.

### 2. Verify (from rules/typescript.md)

```bash
# ripple-lib-server
cd apps/ripple-lib-server && yarn lint && yarn test

# erc20-token
cd apps/erc20-token && npm run lint-js && npm run fmt
```

### 3. Self-Review Checklist

- [ ] No TypeScript errors
- [ ] No `any` types (unless documented reason)
- [ ] Async errors properly handled
- [ ] Auto-generated files not edited

## Related

- `.claude/rules/typescript.md` - TypeScript rules (SSOT)
- `git-workflow` - Branch, commit, PR workflow
- `solidity-development` - For Solidity contracts in erc20-token
