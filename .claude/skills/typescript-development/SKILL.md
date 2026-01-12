---
name: typescript-development
description: TypeScript/JavaScript development workflow for apps/ directory. Use when modifying TypeScript code in ripple-lib-server or JavaScript in erc20-token.
---

# TypeScript/JavaScript Development Workflow

Workflow for TypeScript/JavaScript code changes in `apps/` directory.

## Prerequisites

**Use `git-workflow` Skill** for branch management, commit conventions, and PR creation.

## Applicable Directories

| App | Language | Path |
|-----|----------|------|
| ripple-lib-server | TypeScript | `apps/ripple-lib-server/` |
| erc20-token | JavaScript | `apps/erc20-token/` |

## Verification Commands

**Navigate to app directory first:**

```bash
cd apps/{app-name}
npm install       # Install dependencies (if needed)
npm run lint      # Lint check
npm run format    # Format code
npm run build     # Build
npm test          # Run tests
```

### Quick Reference

| App | Commands |
|-----|----------|
| ripple-lib-server | `cd apps/ripple-lib-server && npm run lint && npm run build` |
| erc20-token | `cd apps/erc20-token && npm run lint && npm run build` |

## Self-Review Checklist

### Code Quality

- [ ] Follows project ESLint configuration
- [ ] Proper TypeScript types (no `any` unless necessary)
- [ ] Async/await error handling
- [ ] No circular dependencies

### Security

- [ ] No hardcoded secrets or API keys
- [ ] No sensitive data in logs
- [ ] Input validation at boundaries

### Auto-Generated Files

**DO NOT EDIT:**

- `apps/ripple-lib-server/src/pb/` (Protocol Buffer generated)
- `apps/erc20-token/build/` (Truffle build artifacts)

## Related Chain Context

| App | Chain |
|-----|-------|
| ripple-lib-server | XRP |
| erc20-token | ETH, ERC20 |

## Related Skills

- `git-workflow` - Branch, commit, PR workflow
- `github-issue-creation` - Task classification
