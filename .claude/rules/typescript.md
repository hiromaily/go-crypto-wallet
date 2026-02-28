---
paths: ["**/*.ts", "**/*.tsx", "**/*.js", "**/*.jsx"]
---

# TypeScript/JavaScript File Rules

## Overview

Rules for modifying TypeScript (`*.ts`, `*.tsx`) and JavaScript (`*.js`, `*.jsx`) files in go-crypto-wallet.

## Applicable Directories

| App              | Language              | Runtime | Path                     | Status                     |
| ---------------- | --------------------- | ------- | ------------------------ | -------------------------- |
| xrpl-grpc-server | TypeScript            | **Bun** | `apps/xrpl-grpc-server/` | **READ-ONLY** (deprecated) |
| eth-contracts    | JavaScript/TypeScript | Node.js | `apps/eth-contracts/`      | Active                     |

## Verification Commands

**Navigate to app directory first:**

### xrpl-grpc-server (TypeScript + Bun) - Active Development

```bash
cd apps/xrpl-grpc-server
bun install           # Install dependencies (if needed)
bun run lint          # Lint with Biome
bun run format        # Format with Biome
bun run typecheck     # TypeScript type checking
bun run dev           # Run dev server with hot reload
bun run build         # Build for production
bun run proto         # Generate protobuf code
```

> **Important**: See @.claude/rules/apps/bun.md for Bun runtime rules.
> **DO NOT** use `npm`/`npx` commands - always use `bun`/`bunx` instead.

### eth-contracts (JavaScript/Solidity)

```bash
cd apps/eth-contracts
npm install           # Install dependencies (if needed)
npm run lint          # Lint Solidity files
npm run lint-js       # Lint JavaScript/TypeScript files
npm run fmt           # Format all files with Prettier
npm run build         # Compile contracts with Truffle
npm run test-all      # Run all tests
```

## Command Summary

| App              | Lint              | Format           | Build           | Test                |
| ---------------- | ----------------- | ---------------- | --------------- | ------------------- |
| xrpl-grpc-server | `bun run lint`    | `bun run format` | `bun run build` | `bun run typecheck` |
| eth-contracts    | `npm run lint-js` | `npm run fmt`    | `npm run build` | `npm run test-all`  |

## Code Style

### TypeScript Best Practices

```typescript
// Good: Explicit types
function getBalance(address: string): Promise<number> {
  // ...
}

// Good: Async/await with error handling
async function fetchData(): Promise<Data> {
  try {
    const result = await api.call();
    return result;
  } catch (error) {
    throw new Error(`Failed to fetch data: ${error.message}`);
  }
}

// Avoid: any type (unless absolutely necessary)
// Bad: function process(data: any)
// Good: function process(data: TransactionData)
```

### Critical Value Handling

Never use nullish coalescing (`??`) with empty string for critical values like secrets or keys.
Instead, throw an error to fail fast and prevent silent failures.

```typescript
// Bad: Silently returns empty string if seed is undefined
const secret = wallet.seed ?? "";

// Good: Fail fast with explicit error
if (!wallet.seed) {
  throw new Error("Failed to generate a wallet seed.");
}
const secret = wallet.seed;
```

### Module Exports

Prefer named exports over default exports for consistency and better tooling support.

```typescript
// Bad: Mixed exports cause inconsistent import styles
export const myService = { ... };
export default myService;

// Good: Named exports only
export const myService = { ... };
```

### Import Order

1. Node.js built-ins
2. External packages
3. Internal modules

#### xrpl-grpc-server (xrpl.js 4.5.0)

```typescript
import * as path from "path";

import { Client, Wallet, xrpToDrops, dropsToXrp } from "xrpl";
import { createConnectRouter } from "@connectrpc/connect";

import { AccountService } from "./services/account";
```

#### eth-contracts

```typescript
import * as path from "path";

import { ethers } from "ethers";

import { TokenService } from "./services/token";
```

## Auto-Generated Files

**DO NOT EDIT:**

| App              | Generated Files                                                  |
| ---------------- | ---------------------------------------------------------------- |
| xrpl-grpc-server | `apps/xrpl-grpc-server/src/protogen/` (Buf/ConnectRPC generated) |
| eth-contracts    | `apps/eth-contracts/build/` (Truffle build artifacts)              |

## Security

- No hardcoded secrets or API keys
- No sensitive data in logs
- Input validation at boundaries
- Use environment variables for configuration

## Quick Checklist

### xrpl-grpc-server (Active)

- [ ] `bun run lint` passes
- [ ] `bun run format` applied
- [ ] `bun run typecheck` passes
- [ ] `bun run build` passes
- [ ] No `any` types (unless documented reason)
- [ ] Async errors properly handled

### eth-contracts

- [ ] `npm run lint-js` passes
- [ ] `npm run fmt` applied
- [ ] `npm run build` passes
- [ ] `npm run test-all` passes

## Related Documentation

- @apps/xrpl-grpc-server/README.md - xrpl-grpc-server documentation
- @apps/xrpl-grpc-server/docs/MIGRATION-GUIDE.md - Migration guide from ripple-lib
- @apps/xrpl-grpc-server/package.json - xrpl-grpc-server scripts
- @apps/eth-contracts/package.json - eth-contracts scripts

## Related Skills

- `typescript-development` - Full TypeScript workflow
- `solidity-development` - For Solidity contracts in eth-contracts

## Related Rules

- @.claude/rules/apps/bun.md - Bun runtime rules (use bun/bunx instead of npm/npx)
- @.claude/rules/apps/xrp-server.md - XRP server directory rules
