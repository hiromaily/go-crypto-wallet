### Protocol Buffer Code (Go) [DEPRECATED for XRP]

> **⚠️ DEPRECATED**: XRP protocol buffers (`proto/rippleapi/`) are no longer used.
>
> **Status**: `apps/xrpl-grpc-server/` and `proto/xrpapi/` are deprecated. XRP functionality now uses native Go implementation with xrpl-go libraries.
>
> **Migration**: See `.kiro/specs/xrp-transaction-flow-alignment/` for current architecture.

**Tool**: protoc (or buf when Edition 2024 is supported)
**Source**: `proto/rippleapi/*.proto` [DEPRECATED]
**Command**: `make proto` [NO LONGER NEEDED FOR XRP]

**Generated Files** [DEPRECATED]:

- `internal/infrastructure/api/xrp/xrp/*.pb.go` (6 files) [No longer generated]
  - `account.pb.go` - Account message types [Deprecated]
  - `account_grpc.pb.go` - Account gRPC service code [Deprecated]
  - `address.pb.go` - Address message types [Deprecated]
  - `address_grpc.pb.go` - Address gRPC service code [Deprecated]
  - `transaction.pb.go` - Transaction message types [Deprecated]
  - `transaction_grpc.pb.go` - Transaction gRPC service code [Deprecated]

**Note**: Protocol buffers were previously used for XRP (Ripple) gRPC communication. This is no longer the case.
