---
paths: ["apps/xrpl-grpc-server/**"]
---

# XRP Server Directory Rules

> **⚠️ DEPRECATED: This directory is no longer used**
>
> **Status**: `apps/xrpl-grpc-server/` and `proto/xrpapi/` are **DEPRECATED**
>
> **Current Approach**: XRP functionality is now implemented using native Go with xrpl-go libraries. All XRP operations are in `internal/infrastructure/api/xrp/`
>
> **Active Specification**: See `.kiro/specs/xrp-transaction-flow-alignment/` for current implementation details.

## Overview

Rules for working with XRP/Ripple-related server applications in the `apps/` directory.

## Directory Status

| Directory                | Status         | Description                           |
| ------------------------ | -------------- | ------------------------------------- |
| `apps/xrpl-grpc-server/` | **DEPRECATED** | No longer used; replaced by native Go |

## Critical Rules

### `apps/xrpl-grpc-server/` - DEPRECATED

This directory is **NO LONGER USED**.

- **DO NOT** add new features here
- **DO NOT** fix bugs in this codebase
- All XRP functionality now uses native Go implementation in `internal/infrastructure/api/xrp/`

## Migration Context

The project has **completed migration** from gRPC-based architecture to native Go:

- **OLD (Deprecated)**: `xrpl-grpc-server` (TypeScript gRPC server) + `proto/xrpapi/` (Protocol Buffers)
- **NEW (Current)**: Native Go implementation using xrpl-go libraries

See `.kiro/specs/xrp-transaction-flow-alignment/` for current architecture.

## Verification Commands [DEPRECATED]

These commands are **NO LONGER NEEDED**:

```bash
# DEPRECATED - Do not use
cd apps/xrpl-grpc-server
bun install           # No longer needed
bun run lint          # No longer needed
bun run build         # No longer needed
bun run dev           # No longer needed

# DEPRECATED - Do not use
make proto            # Protobuf generation no longer needed for XRP
make proto-ts         # TypeScript protobuf no longer needed
```

## Quick Reference

| Action          | Correct Target                                   |
| --------------- | ------------------------------------------------ |
| Add XRP feature | `internal/infrastructure/api/xrp/` (Go codebase) |
| Fix XRP bug     | `internal/infrastructure/api/xrp/` (Go codebase) |

## Related Documentation

- @.kiro/specs/xrp-transaction-flow-alignment/ - Current XRP implementation specification
- @docs/chains/xrp/README.md - XRP chain documentation
- @docs/chains/xrp/architecture-xrpl-grpc-server-version.md - Historical gRPC architecture (OBSOLETE)
