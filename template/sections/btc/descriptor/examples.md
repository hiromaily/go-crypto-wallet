# Descriptor Examples

This page collects runnable examples for descriptor workflows.

## Export and Import a Taproot Account

1. **Export from keygen wallet (receive + change)**

   ```bash
   keygen --coin btc descriptor export \
     --account deposit \
     --output /tmp/deposit_tr.json \
     --format bitcoin-core \
     --include-change
   ```

2. **Inspect the Bitcoin Core JSON**

   ```json
   [
     {
       "desc": "tr([a1b2c3d4/86'/0'/0']xpub.../0/*)",
       "timestamp": "now",
       "range": [0, 1000],
       "watchonly": true
     },
     {
       "desc": "tr([a1b2c3d4/86'/0'/0']xpub.../1/*)",
       "timestamp": "now",
       "range": [0, 1000],
       "watchonly": true
     }
   ]
   ```

3. **Import into the watch wallet**

   ```bash
   watch --coin btc descriptor import \
     --file /tmp/deposit_tr.json \
     --account deposit \
     --start 0 \
     --count 20
   ```

4. **Optional: Import into Bitcoin Core**

   ```bash
   bitcoin-cli importdescriptors "$(cat /tmp/deposit_tr.json)"
   bitcoin-cli deriveaddresses "tr([a1b2c3d4/86'/0'/0']xpub.../0/*)" "[0,2]"
   ```

## Validate a Descriptor File

```bash
watch --coin btc descriptor validate \
  --file ./deposit_all_descriptors.json
```

- Skips address derivation and only validates syntax/structure.
- Useful before sharing files externally.

## Export All Address Types

```bash
keygen --coin btc descriptor export-all \
  --account payment \
  --output ./payment_all.json
```

Outputs descriptors for Taproot, Bech32, P2SH-SegWit, and legacy (receive + change when available).

## Quick Compatibility Smoke Test (Bitcoin Core)

Use the helper script added for compatibility checks:

```bash
scripts/test/test_descriptor_compatibility.sh ./payment_all.json
```

Environment variables:

- `BITCOIN_CLI` (default: `bitcoin-cli`)
- `BITCOIN_CLI_ARGS` (e.g., `-regtest -rpcuser=test -rpcpassword=test`)

## Backup Strategy Example

```bash
# Export to a secure location
keygen --coin btc descriptor export \
  --account deposit \
  --output /secure/backup/deposit_descriptors.txt

# Encrypt for offsite storage
gpg --encrypt --recipient ops@example.com /secure/backup/deposit_descriptors.txt

# Verify checksum after transfer
sha256sum /secure/backup/deposit_descriptors.txt
```

## Handling Multisig Descriptors

- Generation is available for WSH and Taproot script-path descriptors.
- The watch wallet currently rejects multisig imports (`multisig descriptors are not supported yet`).
- Use Export + manual Bitcoin Core import for monitoring until watch wallet multisig support is enabled.
