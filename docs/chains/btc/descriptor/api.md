<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/docs/chains/btc/descriptor/api.tpl.md · Run `make docs` to regenerate.
-->

# Descriptor API (CLI Surface)

Descriptor functionality is exposed through CLI commands for keygen (offline) and watch (online) wallets.

## Keygen Wallet Commands

### Generate a Descriptor

```
keygen --coin btc descriptor generate \
  --account <account> \
  --address-type <taproot|bech32|p2sh-segwit|legacy> \
  [--change] \
  [--required-sigs <n>]  # multisig only
```

- Outputs a single descriptor string to stdout.
- `--change` switches to the change branch (`/1/*`); default is receive (`/0/*`).

### Export Descriptors

```
keygen --coin btc descriptor export \
  --account <account> \
  --output <path> \
  [--format <text|json|bitcoin-core>] \
  [--include-change]
```

- Writes Taproot, Bech32, P2SH-SegWit, and Legacy descriptors to `<path>`.
- `--include-change` writes both receive and change descriptors.
- `--format bitcoin-core` emits JSON suitable for `bitcoin-cli importdescriptors`.

### Export All (Receive + Change)

```
keygen --coin btc descriptor export-all \
  --account <account> \
  --output <path> \
  [--format <text|json|bitcoin-core>]
```

- Convenience wrapper that always includes change descriptors.

## Watch Wallet Commands

### Import Descriptors

```
watch --coin btc descriptor import \
  --account <account> \
  --file <path>|"" \
  [--start <n>] \
  [--count <n>]
```

- Accepts Bitcoin Core JSON (`[{ "desc": "...", ...}]`) or line-delimited descriptors.
- Reads from stdin when `--file` is omitted.
- Derives `count` addresses starting at `start` for each descriptor.
- Current behavior rejects multisig WSH descriptors with an explicit error.

### Validate Descriptors

```
watch --coin btc descriptor validate \
  --file <path>|""
```

- Parses and validates descriptors without writing addresses.
- Useful for CI or pre-flight checks.

## File Formats

- **Text**: One descriptor per line; `#` starts a comment.
- **JSON**: `["desc1", "desc2"]`
- **Bitcoin Core JSON** (default for export): `[{ "desc": "...", "timestamp": "now", "range": [0,1000], "watchonly": true }]`

## Environment Variables

For Bitcoin Core compatibility tests and scripts:

- `BTC_CORE_COMPAT_DESCRIPTOR_FILE`: Descriptor file used by integration tests and smoke script.
- `BITCOIN_CLI`: Path to `bitcoin-cli` (default: `bitcoin-cli`).
- `BITCOIN_CLI_ARGS`: Space-separated args (e.g., `-regtest -rpcuser=test -rpcpassword=test`).

## Error Handling

- Errors are wrapped with context using `%w` for traceability.
- Import collects per-descriptor errors and reports them after processing.
- Validation fails fast on malformed descriptors; multisig imports report explicit unsupported errors.

## References

- CLI definitions: `internal/interface-adapters/cli/keygen/btc/descriptor.go`, `internal/interface-adapters/cli/watch/btc/descriptor.go`
- Use cases: `internal/application/usecase/keygen/btc`, `internal/application/usecase/watch/btc`
- Integration tests: `test/integration/descriptor_workflow_test.go`, `test/integration/descriptor_compatibility_test.go`
