# Protocol Buffer Configuration

This document describes the Protocol Buffer (protobuf) setup in go-crypto-wallet.

## Edition 2024

This project uses **Protobuf Edition 2024**, the latest edition of Protocol Buffers.

### What are Protobuf Editions?

Editions replace the old `syntax = "proto2"` / `syntax = "proto3"` declarations with a more flexible system:

- **Edition 2023**: First edition, unifies proto2 and proto3 behaviors
- **Edition 2024**: Latest edition with improved defaults and features

Reference: [Protobuf Editions Overview](https://protobuf.dev/editions/overview/)

### Why Edition 2024?

- Future-proof: Adopts the latest protobuf features
- Better defaults: Improved naming conventions and field presence semantics
- Compatibility: Uses feature flags for backward compatibility

## Proto Files

Location: `proto/rippleapi/`

| File | Description |
|------|-------------|
| `account.proto` | XRP account information API |
| `address.proto` | XRP address generation API |
| `transaction.proto` | XRP transaction API |

### Feature Flags

Each proto file uses the following feature flags for backward compatibility:

```protobuf
edition = "2024";

// Use proto3 semantics for field presence (implicit presence)
option features.field_presence = IMPLICIT;

// Allow existing camelCase field names without requiring snake_case
option features.enforce_naming_style = STYLE_LEGACY;
```

## Code Generation

### Recommended: protoc (for Edition 2024)

```bash
make proto
```

This uses `protoc` directly, which fully supports Edition 2024.

**Requirements:**

- `protoc` >= 33.0
- `protoc-gen-go` (Go code generator)
- `protoc-gen-go-grpc` (gRPC code generator)

### Alternative: buf (for future use)

```bash
make proto-buf
```

**Note:** As of January 2026, buf CLI v1.63.0 does not yet fully support Edition 2024.
The command will fail with:

```
edition "2024" not yet fully supported; latest supported edition "2023"
```

Once buf adds Edition 2024 support, `make proto-buf` can be used as the primary generation method.

### Proto Formatting (clang-format)

Since buf CLI v1.63.0 doesn't support edition 2024, use `clang-format` for proto formatting:

```bash
make proto-fmt         # Format proto files with clang-format (recommended)
make proto-fmt-check   # Check formatting (used in CI)
```

**Requirements:**
- `clang-format` (install via `apt install clang-format` or `brew install clang-format`)

### buf Utilities (Limited with Edition 2024)

**Note:** As of buf CLI v1.63.0, the following commands fail with edition 2024:

```bash
make proto-fmt-buf       # ❌ Fails - edition 2024 not supported
make proto-lint          # ❌ Fails - edition 2024 not supported
make breaking-proto      # ✅ Works with continue-on-error
```

buf lint is temporarily disabled in CI. Re-enable when buf >= 1.64.0 adds edition 2024 support.

## Generated Code

Output: `internal/infrastructure/api/xrp/protogen/`

### Opaque API Pattern

Edition 2024 generates Go code using the "opaque" API pattern:

**Before (Edition 2023 / proto3):**

```go
// Direct field access
msg := &protogen.Instructions{
    Fee:    "10",
    MaxFee: "100",
}
value := msg.Fee
```

**After (Edition 2024):**

```go
// Builder pattern for construction
msg := protogen.Instructions_builder{
    Fee:    "10",
    MaxFee: "100",
}.Build()

// Getter methods for access
value := msg.GetFee()

// Setter methods for mutation
msg.SetFee("20")
```

### Migration Notes

When upgrading from Edition 2023 to 2024:

1. Update struct literals to use `*_builder{}.Build()` pattern
2. Replace direct field access (`msg.Field`) with getters (`msg.GetField()`)
3. Use setters (`msg.SetField()`) for mutation

## Configuration Files

### buf.yaml

```yaml
version: v2
modules:
  - path: proto/rippleapi
    name: buf.build/hiromaily/go-crypto-wallet-rippleapi
lint:
  use:
    - BASIC
  except:
    - ENUM_VALUE_PREFIX
    - ENUM_ZERO_VALUE_SUFFIX
    - FIELD_LOWER_SNAKE_CASE
    - PACKAGE_DIRECTORY_MATCH
    - PACKAGE_SAME_DIRECTORY
    - DIRECTORY_SAME_PACKAGE
```

### buf.gen.yaml

```yaml
version: v2
managed:
  enabled: false
plugins:
  - remote: buf.build/protocolbuffers/go
    out: internal/infrastructure/api/xrp/protogen
    opt:
      - paths=source_relative
  - remote: buf.build/grpc/go
    out: internal/infrastructure/api/xrp/protogen
    opt:
      - paths=source_relative
inputs:
  - directory: proto/rippleapi
```

## Version Requirements

| Tool | Minimum Version | Edition 2024 Support |
|------|-----------------|---------------------|
| protoc | **>= 33.0** | Full support |
| buf CLI | >= 1.64.0 | Expected (v1.63.0 does not support) |
| protoc-gen-go | latest | Full support |
| protoc-gen-go-grpc | latest | Full support |

### Version Checks

The Makefile includes automatic version verification:

```bash
make proto-version-check  # Check protoc >= 33.0
make buf-version-check    # Check buf >= 1.64.0
```

These checks run automatically when using `make proto` or `make proto-buf`.

## References

- [Protobuf Editions Overview](https://protobuf.dev/editions/overview/)
- [Protobuf Edition 2024 Announcement](https://protobuf.dev/news/2025-06-27/)
- [buf Documentation](https://buf.build/docs/)
- [Go Protocol Buffers](https://protobuf.dev/reference/go/go-generated/)
