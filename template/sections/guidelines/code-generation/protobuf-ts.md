### Protocol Buffer Code (TypeScript) [DEPRECATED]

> **⚠️ DEPRECATED**: XRP gRPC server (`apps/xrpl-grpc-server/`) is no longer used.
>
> **Status**: All XRP functionality migrated to native Go implementation. TypeScript protobuf code generation is no longer needed.

**Tool**: protoc with protoc-gen-es and protoc-gen-connect-es [NO LONGER USED]
**Source**: `proto/rippleapi/*.proto` [DEPRECATED]
**Command**: `make proto-ts` [NO LONGER NEEDED]

**Generated Files** [DEPRECATED]:

- `apps/xrpl-grpc-server/src/protogen/*.ts` - TypeScript protocol buffer code [No longer generated]
  - `account_pb.ts` - Account message types [Deprecated]
  - `account_connect.ts` - Account service client [Deprecated]
  - `address_pb.ts` - Address message types [Deprecated]
  - `address_connect.ts` - Address service client [Deprecated]
  - `transaction_pb.ts` - Transaction message types [Deprecated]
  - `transaction_connect.ts` - Transaction service client [Deprecated]
