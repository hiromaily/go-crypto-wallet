# Descriptor Architecture

Descriptor functionality follows the project's Clean Architecture pattern.

```
CLI (interface-adapters)
    ↓
Application use cases (keygen/watch)
    ↓
Domain (descriptor types, validation)
    ← Infrastructure (descriptor parser/service, file IO)
```

## Components

- **Domain**: `internal/domain/wallet.Descriptor`, validation helpers, descriptor types (PKH, SH(WPKH), WPKH, TR, WSH). No infrastructure dependencies.
- **Application (keygen)**: `GenerateDescriptorUseCase` builds descriptors from account/auth repositories; `ExportDescriptorUseCase` formats and writes descriptor sets.
- **Application (watch)**: `ImportDescriptorUseCase` parses/validates descriptors, derives addresses, and stores them in the watch wallet repository.
- **Infrastructure**: `DescriptorService` (generator), `DescriptorParser` (parse + extract keys), file writer implementations, Bitcoin Core compatibility helpers.
- **Interface Adapters**: CLI commands under `internal/interface-adapters/cli/keygen/btc` and `internal/interface-adapters/cli/watch/btc`.

## Data Flows

### Export Flow

```
CLI keygen descriptor export
    → ExportDescriptorUseCase
      → GenerateDescriptorUseCase
        → DescriptorService (format derivation paths)
      → FileWriter (text/JSON/Bitcoin Core)
```

### Import Flow

```
Descriptor file (text or Bitcoin Core JSON)
    → ImportDescriptorUseCase
      → DescriptorParser (syntax + key extraction)
      → Address derivation (HD key + descriptor type)
      → Address repository (watch wallet)
```

## Key Decisions

- **Separation of concerns**: Domain types are free of infrastructure; parsing/generation lives in infrastructure; orchestration lives in use cases.
- **Format flexibility**: Export supports text, JSON array, and Bitcoin Core JSON; import accepts Bitcoin Core JSON or line-delimited text.
- **Multisig policy**: Generation supports multisig descriptors; watch import currently rejects WSH multisig for safety.
- **Change branch**: Export can include change descriptors to keep parity with legacy flows.

## Testing Strategy

- **Unit**: Descriptor parsing/formatting, domain validation, CLI flag parsing, use case constructors.
- **Integration**: `test/integration/descriptor_workflow_test.go` exercises export → import → storage with in-memory repositories; `test/integration/descriptor_compatibility_test.go` verifies Bitcoin Core compatibility (requires external node, `integration` tag).
- **Benchmark**: `internal/infrastructure/api/btc/btc/descriptor_bench_test.go` measures generation and parsing throughput.

## Security Considerations

- Only extended public keys are used for descriptors; private keys are never logged or written.
- Validate checksum and derivation paths before import; reject malformed descriptors.
- Treat descriptor files as sensitive metadata (privacy/leakage risk).
- Multisig imports remain disabled in the watch wallet until the validation and storage path are hardened.
