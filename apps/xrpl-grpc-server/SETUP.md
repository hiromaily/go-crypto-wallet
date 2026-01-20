# Setup Instructions for xrpl-grpc-server

## Prerequisites

1. **Bun** >= 1.3.6
   ```bash
   # Install Bun
   curl -fsSL https://bun.sh/install | bash

   # Verify version
   bun --version
   ```

2. **protoc** >= 33.0 (for Edition 2024 support)
   ```bash
   # Check protoc version
   protoc --version

   # Install protoc if needed
   # See https://grpc.io/docs/protoc-installation/
   ```

3. **Buf CLI** >= 1.64.0 (optional, for future use)

   ⚠️ **Note**: Buf CLI v1.63.0 does NOT support Protobuf Edition 2024.
   Currently, `protoc` is used for code generation via `make proto-ts`.

## Installation Steps

1. **Install dependencies**
   ```bash
   cd apps/xrpl-grpc-server
   bun install
   ```

2. **Generate TypeScript proto files**
   ```bash
   bun run proto
   # Or from repository root:
   make proto-ts
   ```

   This will:
   - Use protoc (>= 33.0) to read proto files from `proto/rippleapi/`
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
| `../../make/codegen.mk` | Makefile targets for proto generation (`proto-ts`) |
| `../../buf.yaml` | Buf module configuration (for future buf CLI use) |
| `../../buf.gen.yaml` | Buf code generation config (for future buf CLI use) |
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

### Protobuf Edition 2024 - Using protoc Instead of Buf

**The proto files in this project use `edition = "2024"`** (upgraded in PR #469, issue #409).

**Current Tool Support Status:**
- ✅ **protoc >= 33.0**: Full Edition 2024 support (used for code generation via `make proto` and `make proto-ts`)
- ❌ **buf CLI v1.63.0**: Does NOT support Edition 2024 (support expected in v1.64.0+)

**Current Approach:**
- TypeScript code generation uses `protoc` via `make proto-ts`
- Go code generation uses `protoc` via `make proto`
- Buf CLI configuration files (`buf.yaml`, `buf.gen.yaml`) are kept for future use

**Buf CLI Commands (will fail until buf >= 1.64.0):**
```bash
buf generate     # WILL FAIL with: "edition 2024 not yet fully supported"
buf lint         # WILL FAIL
buf format       # WILL FAIL
```

**For More Information:**
- See `../../docs/proto.md` - Comprehensive Edition 2024 documentation
- See `../../make/codegen.mk` - Proto generation targets and version checks
- See PR #469 - Edition 2024 migration implementation

## Troubleshooting

### protoc version too old

If you see "Edition 2024 requires protoc >= 33.0":
```bash
# Check your protoc version
protoc --version

# Upgrade protoc - see https://grpc.io/docs/protoc-installation/
```

### Missing protoc-gen-* executables

The `make proto-ts` target requires local plugins. Install them via bun:
```bash
cd apps/xrpl-grpc-server
bun install
```

This installs:
- `@bufbuild/protoc-gen-es` → creates `node_modules/.bin/protoc-gen-es`
- `@connectrpc/protoc-gen-connect-es` → creates `node_modules/.bin/protoc-gen-connect-es`

### buf generate fails with edition error

This is expected. Use `make proto-ts` instead of `buf generate` until buf CLI >= 1.64.0 is released.

**Do not**:
- Downgrade proto files to edition 2023 (breaks compatibility with Go code)
- Modify feature flags in proto files

## Next Steps

After setup is complete:
1. Implement xrpl.js client wrapper (Issue #473)
2. Implement RippleAddressAPI service (Issue #474)
3. Implement RippleAccountAPI service (Issue #475)
4. Implement RippleTransactionAPI service (Issue #476)
5. Set up ConnectRPC server (Issue #477)
