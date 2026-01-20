# Protobuf Edition 2024 Support Status

This document describes the current state of Protobuf Edition 2024 support in the xrpl-grpc-server and the planned migration path.

## Current Situation

### Proto Files

The proto files in `proto/rippleapi/` use **Protobuf Edition 2024**:

```protobuf
edition = "2024";

package rippleapi.account;

option features.field_presence = IMPLICIT;
option features.enforce_naming_style = STYLE_LEGACY;
```

### Code Generation Status

| Tool | Version | Edition 2024 Support | Status |
|------|---------|---------------------|--------|
| `protoc` | 33.4+ | **Supported** | Used |
| `@bufbuild/protoc-gen-es` | 2.10.2 | **Supported** | Used |
| `@connectrpc/protoc-gen-connect-es` | 1.7.0 | **NOT Supported** (max Edition 2023) | **Not used** |

### Current Approach

Since `protoc-gen-connect-es` v1.7.0 does not support Edition 2024, we use a **manual implementation** approach:

1. **Generate protobuf types only** using `protoc-gen-es` v2.10.2
   - This generates message types and service descriptors in `src/gen/`
   - Files: `account_pb.ts`, `address_pb.ts`, `transaction_pb.ts`

2. **Manually implement ConnectRPC handlers** in `src/server.ts`
   - Uses generated message types and service descriptors
   - Maps between proto types and service layer types
   - Provides `createRouter()` and `createServerFetchHandler()` functions

### Code Generation Command

To generate the protobuf types (without Connect stubs):

```bash
# From workspace root
mkdir -p apps/xrpl-grpc-server/src/gen && protoc \
  --plugin=protoc-gen-es=apps/xrpl-grpc-server/node_modules/.bin/protoc-gen-es \
  --es_out=apps/xrpl-grpc-server/src/gen \
  --es_opt=target=ts \
  -I proto/rippleapi \
  proto/rippleapi/account.proto proto/rippleapi/address.proto proto/rippleapi/transaction.proto
```

Note: The standard `make proto-ts` command will fail because it also tries to run `protoc-gen-connect-es`.

### Error When Using protoc-gen-connect-es

If you try to generate Connect stubs with Edition 2024 proto files:

```
protoc-gen-connect-es: Edition EDITION_2024 is later than the maximum supported edition EDITION_2023
--connect-es_out: protoc-gen-connect-es: Plugin failed with status code 1
```

## Architecture

### Current (Manual Implementation)

```
proto/rippleapi/*.proto (Edition 2024)
        │
        ▼
  protoc-gen-es v2.10.2
        │
        ▼
src/gen/*_pb.ts (message types + service descriptors)
        │
        ▼
src/server.ts (manual ConnectRPC handlers)
        │
        ▼
src/index.ts (Bun.serve() entry point)
```

### Future (Generated Connect Stubs)

```
proto/rippleapi/*.proto (Edition 2024)
        │
        ├──────────────────────────┐
        ▼                          ▼
  protoc-gen-es v2.x      protoc-gen-connect-es v2.x (when released)
        │                          │
        ▼                          ▼
src/gen/*_pb.ts           src/gen/*_connect.ts (generated stubs)
        │                          │
        └──────────────────────────┤
                                   ▼
                       src/server.ts (uses generated stubs)
```

## Future Migration Plan

### When to Migrate

Migrate when `@connectrpc/protoc-gen-connect-es` releases a version that supports Edition 2024 (likely v2.x).

Check for updates:
- [ConnectRPC GitHub](https://github.com/connectrpc/connect-es)
- [npm package](https://www.npmjs.com/package/@connectrpc/protoc-gen-connect-es)

### Migration Steps

1. **Update package.json**

   ```json
   {
     "devDependencies": {
       "@connectrpc/protoc-gen-connect-es": "^2.x.x"
     }
   }
   ```

2. **Run full proto generation**

   ```bash
   make proto-ts
   ```

   This should now succeed and generate both:
   - `src/gen/*_pb.ts` - Message types
   - `src/gen/*_connect.ts` - Connect service stubs

3. **Simplify server.ts**

   Replace manual handler implementations with generated stubs:

   ```typescript
   // Before (manual)
   import { RippleAccountAPI } from './gen/account_pb';
   
   export function createRouter(): ConnectRouter {
     return createConnectRouter()
       .service(RippleAccountAPI, {
         getAccountInfo: async (request) => {
           // Manual mapping code...
         },
       });
   }
   
   // After (generated stubs)
   import { RippleAccountAPI } from './gen/account_connect';
   
   export function createRouter(): ConnectRouter {
     return createConnectRouter()
       .service(RippleAccountAPI, accountServiceImpl);
   }
   ```

4. **Update imports**

   Service descriptors will come from `*_connect.ts` files instead of `*_pb.ts` files.

5. **Remove manual type mappings**

   The generated stubs will handle type conversions automatically.

6. **Update documentation**

   - Update this document to reflect the migration
   - Update MIGRATION-GUIDE.md if needed

### Verification After Migration

```bash
# Regenerate proto files
make proto-ts

# Verify no TypeScript errors
bun run typecheck

# Verify no lint errors
bun run lint

# Test the server
bun run dev
```

## Related Files

| File | Description |
|------|-------------|
| `proto/rippleapi/*.proto` | Edition 2024 proto definitions |
| `src/gen/*_pb.ts` | Generated message types (protoc-gen-es) |
| `src/server.ts` | Manual ConnectRPC handler implementation |
| `src/services/*.ts` | Service layer implementations |
| `make/codegen_proto.mk` | Proto generation makefile |
| `package.json` | Dependencies including proto tools |

## References

- [Protobuf Editions](https://protobuf.dev/editions/)
- [ConnectRPC Documentation](https://connectrpc.com/docs/)
- [@bufbuild/protoc-gen-es](https://github.com/bufbuild/protobuf-es)
- [@connectrpc/protoc-gen-connect-es](https://github.com/connectrpc/connect-es)
