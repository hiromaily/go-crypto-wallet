---
paths: ["apps/xrpl-grpc-server/**"]
---

# Bun Runtime Rules

## Overview

Rules for projects using [Bun](https://bun.sh/) runtime instead of Node.js.

## Critical Rule: Use Bun Commands Only

**DO NOT use npm/npx commands in Bun projects.**

| Wrong (npm/npx) | Correct (bun/bunx) |
|-----------------|-------------------|
| `npm install` | `bun install` |
| `npm run <script>` | `bun run <script>` |
| `npm test` | `bun test` |
| `npm ci` | `bun install --frozen-lockfile` |
| `npx <package>` | `bunx <package>` |
| `npm init` | `bun init` |
| `npm add <pkg>` | `bun add <pkg>` |
| `npm remove <pkg>` | `bun remove <pkg>` |

## Why This Matters

1. **Lock file mismatch**: npm creates `package-lock.json`, Bun uses `bun.lock`
2. **Dependency resolution**: Different algorithms may cause inconsistencies
3. **Performance**: Bun is significantly faster
4. **Project consistency**: Keep tooling consistent across the team

## Command Reference

### Package Management

```bash
# Install all dependencies
bun install

# Add a dependency
bun add <package>
bun add <package> -D  # dev dependency

# Remove a dependency
bun remove <package>

# Update dependencies
bun update
```

### Running Scripts

```bash
# Run a script from package.json
bun run <script>

# Run TypeScript directly
bun run src/index.ts

# Run with watch mode
bun --watch src/index.ts
```

### Executing Packages

```bash
# Execute a package binary (like npx)
bunx <package>

# Examples
bunx tsc --noEmit
bunx biome check .
```

## Applicable Directories

This rule applies to:

- `apps/xrpl-grpc-server/` - Uses Bun runtime

## Verification

Check if a project uses Bun:

1. Look for `bun.lock` file (not `package-lock.json`)
2. Check `engines` field in `package.json`:

   ```json
   "engines": {
     "bun": ">=1.3.6"
   }
   ```

## Common Mistakes to Avoid

```bash
# WRONG - Do not use these
npm install
npm run dev
npx biome check .

# CORRECT - Use these instead
bun install
bun run dev
bunx biome check .
```
