# Setup Instructions for xrpl-grpc-server

## Prerequisites

1. **Bun** >= 1.3.6
   ```bash
   # Install Bun
   curl -fsSL https://bun.sh/install | bash

   # Verify version
   bun --version
   ```

2. **Buf CLI** >= 1.64.0 (NOT YET AVAILABLE)
   ```bash
   # Check if buf is installed
   buf --version

   # Install buf if needed
   # See https://buf.build/docs/installation
   ```

   ⚠️ **Important**: Buf CLI v1.63.0 does NOT support Protobuf Edition 2024.
   Code generation will NOT work until buf >= 1.64.0 is released.

## Installation Steps

1. **Install dependencies**
   ```bash
   cd apps/xrpl-grpc-server
   bun install
   ```

2. **Generate TypeScript proto files**
   ```bash
   bun run proto
   ```

   This will:
   - Use buf to read proto files from `../../proto/rippleapi/`
   - Generate TypeScript files in `src/gen/`
   - Use @bufbuild/protoc-gen-es for protobuf messages
   - Use @connectrpc/protoc-gen-connect-es for gRPC services

3. **Verify the setup**
   ```bash
   # Type checking
   bun run typecheck

   # Linting
   bun run lint
   ```

## Configuration Files

| File | Purpose |
|------|---------|
| `buf.yaml` | Buf module configuration with googleapis dependency |
| `buf.gen.yaml` | Code generation configuration for TypeScript |
| `package.json` | Dependencies and scripts |
| `tsconfig.json` | TypeScript compiler options |
| `biome.json` | Linting and formatting rules |

## Generated Files

After running `bun run proto`, the following structure will be created:

```
src/gen/
├── rippleapi/
│   ├── account_pb.ts           # Account message types
│   ├── account_connect.ts      # Account service client
│   ├── address_pb.ts           # Address message types
│   ├── address_connect.ts      # Address service client
│   ├── transaction_pb.ts       # Transaction message types
│   └── transaction_connect.ts  # Transaction service client
```

**Note**: Generated files are ignored by git (see `.gitignore`) and should not be manually edited.

## Known Issues

### Protobuf Edition 2024 - Buf CLI Not Usable

⚠️ **CRITICAL**: Buf CLI **CANNOT** currently be used for code generation due to Protobuf Edition 2024.

**The proto files in this project use `edition = "2024"`** (upgraded in PR #469, issue #409).

**Current Tool Support Status:**
- ✅ **protoc >= 33.0**: Full Edition 2024 support (used for Go code generation via `make proto`)
- ❌ **buf CLI v1.63.0**: Does NOT support Edition 2024 (support expected in v1.64.0+)
- ❓ **TypeScript plugins**: Unknown - will be tested once buf CLI adds support

**Impact:**
```bash
bun run proto    # WILL FAIL with: "edition 2024 not yet fully supported"
buf generate     # WILL FAIL with: "edition 2024 not yet fully supported"
buf lint         # WILL FAIL
buf format       # WILL FAIL
```

**What This Means:**
- The configuration in this directory (`buf.yaml`, `buf.gen.yaml`) is **correct and future-proof**
- Code generation **will work** once buf CLI >= 1.64.0 is released
- This is a **temporary limitation** of the buf CLI tool, not our configuration

**For More Information:**
- See `../../docs/proto.md` - Comprehensive Edition 2024 documentation
- See `../../make/codegen.mk` - Proto generation targets and version checks
- See PR #469 - Edition 2024 migration implementation

## Troubleshooting

### buf generate fails with edition error

**Expected behavior**: This is a known limitation. The configuration is correct and will work once:
1. buf CLI >= 1.64.0 is released with Edition 2024 support
2. Dependencies are updated with `bun install`

**Current workaround**: N/A - wait for buf CLI update

**Do not**:
- Downgrade proto files to edition 2023 (breaks compatibility with Go code)
- Modify feature flags in proto files

### Missing protoc-gen-* executables

The buf.gen.yaml uses local plugins (`local: protoc-gen-es`). These must be installed via npm/bun:
```bash
bun install
```

This installs:
- `@bufbuild/protoc-gen-es` → creates `node_modules/.bin/protoc-gen-es`
- `@connectrpc/protoc-gen-connect-es` → creates `node_modules/.bin/protoc-gen-connect-es`

Buf will look for these in your PATH or node_modules/.bin/.

## Next Steps

After setup is complete:
1. Implement xrpl.js client wrapper (Issue #473)
2. Implement RippleAddressAPI service (Issue #474)
3. Implement RippleAccountAPI service (Issue #475)
4. Implement RippleTransactionAPI service (Issue #476)
5. Set up ConnectRPC server (Issue #477)
