# Descriptor Compatibility with Bitcoin Core

This guide describes how to validate descriptor compatibility between this wallet and Bitcoin Core.

## Version Requirements
- Minimum: Bitcoin Core v0.17.0 (descriptor support)
- Taproot: v22.0+ required
- MuSig2: v25.0+ (limited/experimental support)

## Environment Variables
- `BITCOIN_CLI` (optional): path to `bitcoin-cli` (default: `bitcoin-cli`)
- `BITCOIN_CLI_ARGS` (optional): space-separated arguments for `bitcoin-cli` (e.g., `-regtest -rpcuser=test -rpcpassword=test`)
- `BTC_CORE_COMPAT_DESCRIPTOR_FILE` (optional): descriptor file (Bitcoin Core JSON format or plain text) used by integration tests and scripts

## Integration Tests
Integration tests are behind the `integration` build tag and require a running Bitcoin Core node.

```bash
BTC_CORE_COMPAT_DESCRIPTOR_FILE=/path/to/descriptors.json \
BITCOIN_CLI_ARGS="-regtest -rpcuser=test -rpcpassword=test" \
go test -tags=integration ./test/integration -run TestDescriptorCompatibility_ImportAndDerive -v
```

Tests will:
1) Import descriptors from the supplied file using `bitcoin-cli importdescriptors`
2) Derive a small address range with `bitcoin-cli deriveaddresses`

Tests skip if `BTC_CORE_COMPAT_DESCRIPTOR_FILE` is unset.

## Manual Smoke Test Script
A lightweight script is provided for quick checks:

```
scripts/test/test_descriptor_compatibility.sh /path/to/descriptors.json
```

The script:
1) Shows `bitcoin-cli -version`
2) Runs `importdescriptors` with the given file
3) Extracts the first descriptor and runs `deriveaddresses` for two addresses

Configure node access via `BITCOIN_CLI_ARGS` if needed.

## Known Limitations
- Multisig descriptors and MuSig2 may require additional setup in Bitcoin Core; integration tests currently target single-key descriptors and small ranges.
- Descriptor content must match the network of the target Bitcoin Core node (mainnet/testnet/regtest).
- The watch wallet’s descriptor import currently treats complex multisig as unsupported; document and test accordingly before production use.
