---
paths: ["apps/ripple-lib-server/**", "apps/xrpl-grpc-server/**"]
---

# XRP Server Directory Rules

## Overview

Rules for working with XRP/Ripple-related server applications in the `apps/` directory.

## Directory Status

| Directory | Status | Description |
|-----------|--------|-------------|
| `apps/ripple-lib-server/` | **READ-ONLY** | Legacy server (deprecated) |
| `apps/xrpl-grpc-server/` | **Active** | Current work target |

## Critical Rules

### `apps/ripple-lib-server/` - DO NOT MODIFY

This directory is **write-prohibited**. It contains the legacy Ripple library server that is being deprecated.

- **DO NOT** create, edit, or delete any files in this directory
- **DO NOT** add new features or fix bugs here
- **READ ONLY** for reference purposes

If you need to understand the legacy implementation, you may read files for reference only.

### `apps/xrpl-grpc-server/` - Active Development

This is the **active work target** for XRP server functionality.

- All new XRP server features should be implemented here
- Bug fixes and improvements go here
- Uses modern XRPL library and gRPC

## Migration Context

The project is migrating from:

- **Old**: `ripple-lib-server` (legacy ripple-lib)
- **New**: `xrpl-grpc-server` (modern xrpl.js with gRPC)

See `apps/xrpl-grpc-server/docs/MIGRATION-GUIDE.md` for migration details.

## Verification Commands

For `xrpl-grpc-server` only:

```bash
cd apps/xrpl-grpc-server
bun install           # Install dependencies
bun run lint          # Lint with Biome
bun run build         # Build TypeScript
bun run dev           # Run server in dev mode
```

## Quick Reference

| Action | Correct Target |
|--------|----------------|
| Add XRP feature | `apps/xrpl-grpc-server/` |
| Fix XRP bug | `apps/xrpl-grpc-server/` |
| Reference old code | `apps/ripple-lib-server/` (read-only) |

## Related Documentation

- @apps/xrpl-grpc-server/README.md - Server documentation
- @apps/xrpl-grpc-server/docs/MIGRATION-GUIDE.md - Migration guide
- @docs/task-contexts/chains/xrp.md - XRP chain context
