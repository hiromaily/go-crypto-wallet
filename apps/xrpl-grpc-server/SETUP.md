# Setup Instructions for xrpl-grpc-server

## Prerequisites

1. **Bun** >= 1.3.6

   ```bash
   # Install Bun
   curl -fsSL https://bun.sh/install | bash

   # Verify version
   bun --version
   ```

2. **protoc** >= 33.0

   ```bash
   # Check protoc version
   protoc --version

   # Install protoc if needed
   # See https://grpc.io/docs/protoc-installation/
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
   # Or from repository root:
   make proto-ts
   ```

   This will generate TypeScript files in `src/gen/` from proto files in `proto/rippleapi/`.

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
| `../../make/codegen_proto.mk` | Proto generation targets (tool selection via `PROTO_TOOL` flag) |
| `../../buf.yaml` | Buf module configuration |
| `../../buf.gen.yaml` | Buf code generation config |
| `package.json` | Dependencies and scripts |
| `tsconfig.json` | TypeScript compiler options |
| `biome.json` | Linting and formatting rules |

**Note**: The Makefile manages tool selection (protoc vs buf) via `PROTO_TOOL` flag. Currently defaults to `protoc` due to Edition 2024 support. When buf CLI adds support, the default can be changed without modifying this project.

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

## Troubleshooting

### protoc version too old

If you see "Edition 2024 requires protoc >= 33.0":

```bash
# Check your protoc version
protoc --version

# Upgrade protoc - see https://grpc.io/docs/protoc-installation/
```

### Missing protoc-gen-* executables

Install dependencies via bun:

```bash
cd apps/xrpl-grpc-server
bun install
```

## Technical Notes

This project uses Protobuf Edition 2024. The Makefile automatically selects the appropriate tool via `PROTO_TOOL` flag. For detailed information, see `../../docs/proto.md`.

## Next Steps

After setup is complete:

1. Implement xrpl.js client wrapper (Issue #473)
2. Implement RippleAddressAPI service (Issue #474)
3. Implement RippleAccountAPI service (Issue #475)
4. Implement RippleTransactionAPI service (Issue #476)
5. Set up ConnectRPC server (Issue #477)
