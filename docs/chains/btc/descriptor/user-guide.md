<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/docs/chains/btc/descriptor/user-guide.tpl.md · Run `make docs` to regenerate.
-->

# Descriptor Wallet User Guide

This guide explains how to generate, export, and import Bitcoin descriptor wallets using the CLI tools in this repository. It targets operators of the keygen (offline) and watch (online) wallets.

## Quick Start

1. **Generate a descriptor (keygen wallet)**

   ```bash
   keygen --coin btc descriptor generate \
     --account deposit \
     --address-type taproot
   ```

   Example output: `tr([a1b2c3d4/86'/0'/0']xpub.../0/*)`

2. **Export descriptors for Bitcoin Core**

   ```bash
   keygen --coin btc descriptor export \
     --account deposit \
     --output deposit_descriptors.json \
     --format bitcoin-core \
     --include-change
   ```

3. **Import descriptors into the watch wallet**

   ```bash
   watch --coin btc descriptor import \
     --file deposit_descriptors.json \
     --account deposit \
     --start 0 \
     --count 100
   ```

4. **Optional: Import into Bitcoin Core for monitoring**

   ```bash
   bitcoin-cli importdescriptors "$(cat deposit_descriptors.json)"
   ```

## Descriptor Types

| Address Type | Descriptor | Use Case |
|--------------|------------|----------|
| Legacy       | `pkh([fingerprint/44'/0'/0']xpub.../0/*)` | Backward compatibility |
| P2SH-SegWit  | `sh(wpkh([fingerprint/49'/0'/0']xpub.../0/*))` | Wrapped SegWit |
| Bech32       | `wpkh([fingerprint/84'/0'/0']xpub.../0/*)` | Native SegWit |
| Taproot      | `tr([fingerprint/86'/0'/0']xpub.../0/*)` | BIP86 single key |
| Multisig (WSH) | `wsh(sortedmulti(...))` | Supported for generation; watch wallet import currently rejects multisig descriptors |

## Common Operations

- **Generate all descriptors for an account**

  ```bash
  keygen --coin btc descriptor export-all \
    --account deposit \
    --output deposit_all_descriptors.json
  ```

- **Validate before importing**

  ```bash
  watch --coin btc descriptor validate \
    --file deposit_all_descriptors.json
  ```

- **Backup descriptors**

  ```bash
  keygen --coin btc descriptor export \
    --account deposit \
    --output /secure/backup/deposit_descriptors.txt

  # Optional: encrypt the backup
  gpg --encrypt --recipient you@example.com /secure/backup/deposit_descriptors.txt
  ```

- **Rotate/change branch**: rerun export with `--include-change` to keep receive and change branches in sync.

## Troubleshooting

- **"Invalid descriptor"**: Confirm BIP380 syntax, correct derivation path, and a valid extended public key. Remove any stray whitespace or comments in text files.
- **"Addresses don't match"**: Ensure the same descriptor file is used for both watch wallet and Bitcoin Core. Verify the network (mainnet/testnet/regtest) and derivation indices.
- **"multisig descriptors are not supported yet"**: Watch wallet import currently rejects multisig WSH descriptors. Track multisig support before production use.

## Security Notes

- Never export or log extended private keys (xprv/tprv). Descriptor flows use extended public keys.
- Store descriptor backups offline and encrypt them at rest.
- Descriptors reveal address structure and history; treat files as sensitive.
- Validate descriptor files from external sources before import.

## Related Material

- `docs/descriptor_examples.md` for runnable samples
- `docs/descriptor_migration.md` for moving from legacy wallets
- `docs/api/descriptor_api.md` for CLI/API reference
