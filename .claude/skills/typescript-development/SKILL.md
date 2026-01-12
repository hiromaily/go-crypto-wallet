---
name: typescript-development
description: Workflow for TypeScript/JavaScript development in apps/ directory. Use when modifying TypeScript code in ripple-lib-server or JavaScript in erc20-token.
---

# TypeScript/JavaScript Development Workflow

Standard workflow for TypeScript/JavaScript code changes in `apps/` directory.

## Applicable Directories

| App | Language | Path |
|-----|----------|------|
| ripple-lib-server | TypeScript | `apps/ripple-lib-server/` |
| erc20-token | JavaScript | `apps/erc20-token/` |

## Branch Management

Same as Go development:

```bash
git fetch origin
git checkout main
git reset --hard origin/main
git checkout -b {branch-type}/issue-{number}-{brief-description}
```

## Verification Commands

**Navigate to the app directory first:**

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
- [ ] Proper imports (no circular dependencies)

### Security

- [ ] No hardcoded secrets or API keys
- [ ] No sensitive data in logs
- [ ] Input validation at boundaries

### Auto-Generated Files

- [ ] NOT editing files in:
  - `apps/ripple-lib-server/src/pb/` (Protocol Buffer generated)
  - `apps/erc20-token/build/` (Truffle build artifacts)

## Commit Message Format

Same as Go development:

```
{type}: {brief description}

- {detail 1}
- {detail 2}

Closes #{issue_number}
```

## Related Chain Context

| App | Related Chain |
|-----|---------------|
| ripple-lib-server | XRP |
| erc20-token | ETH, ERC20 |
