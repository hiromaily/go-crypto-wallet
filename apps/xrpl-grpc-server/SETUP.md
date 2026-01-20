# Setup Instructions for xrpl-grpc-server

## Prerequisites

1. **Bun** >= 1.0.0
   ```bash
   # Install Bun
   curl -fsSL https://bun.sh/install | bash
   ```

2. **Buf CLI** >= 1.63.0
   ```bash
   # Check if buf is installed
   buf --version

   # Install buf if needed
   # See https://buf.build/docs/installation
   ```

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

### Protobuf Edition 2024 Support

⚠️ **Important**: The proto files use `edition = "2024"` which is not yet fully supported by buf CLI v1.63.0.

**Status as of January 2026:**
- ✅ protoc >= 33.0 fully supports Edition 2024 (used for Go code generation)
- ❌ buf CLI v1.63.0 does NOT support Edition 2024 (expected in v1.64.0+)
- ❓ TypeScript plugins (@bufbuild/protoc-gen-es, @connectrpc/protoc-gen-connect-es) support status unknown

**Impact:**
Running `bun run proto` will currently fail with:
```
edition "2024" not yet fully supported; latest supported edition "2023"
```

**Resolution:**
This configuration is future-proof and ready for when buf CLI >= 1.64.0 is released with Edition 2024 support.

For reference, see:
- Issue #409: Protobuf Edition 2024 upgrade
- PR #469: Implementation details
- `docs/proto.md`: Comprehensive protobuf documentation

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
