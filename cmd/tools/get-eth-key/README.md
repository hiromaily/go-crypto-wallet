# get-eth-key

A command-line tool to extract Ethereum private keys from keystore files or generate them from mnemonic phrases using BIP39/BIP44 standards.

## Overview

`get-eth-key` provides two modes of operation:

1. **Extract private key from keystore**: Retrieve a private key from an existing keystore file using an Ethereum address and keystore directory.
2. **Generate from mnemonic**: Derive a private key and address from a BIP39 mnemonic phrase and BIP44 HD wallet path.

## Installation

### Building

```bash
cd cmd/tools/get-eth-key
go build -o get-eth-key main.go wallet.go
```

Or using Makefile:

```bash
make build
```

## Usage

### Command-line Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `-addr` | Mode 1 | Ethereum address (hex format with 0x prefix) |
| `-dir` | Mode 1 | Keystore directory path |
| `-mnemonic` | Mode 2 | BIP39 mnemonic phrase (12 or 24 words) |
| `-hdpath` | Mode 2 | BIP44 HD wallet derivation path (e.g., `m/44'/60'/0'/0/0`) |

### Mode 1: Extract Private Key from Keystore

This mode extracts a private key from an existing keystore file.

**Requirements:**

- The keystore file must exist at `{dir}/{address}` (where `address` is the Ethereum address without the `0x` prefix)
- The keystore file must be in standard Ethereum keystore format (JSON)

**Example:**

```bash
./get-eth-key --addr 0x71678cd07cfac46c2dc427f999abf46aae115925 --dir ./keystore
```

**Output:**

```
addr: 0x71678cd07cfac46c2dc427f999abf46aae115925
keyDir: ./keystore
mnemonic: 
hdPath: 
[Mode] From address and keystore directory
&{...private key object...}
```

### Mode 2: Generate from Mnemonic and HD Path

This mode generates a private key and address from a BIP39 mnemonic phrase using BIP44 derivation.

**Requirements:**

- Valid BIP39 mnemonic phrase (12 or 24 words)
- Valid BIP44 HD wallet path format: `m/44'/{coin_type}'/{account}'/{change}/{index}`
  - Purpose must be `44` (BIP44)
  - Coin type for Ethereum is `60`
  - Change must be `0` (external) or `1` (internal)
  - Example: `m/44'/60'/0'/0/0` (Ethereum, account 0, external, index 0)

**Example:**

```bash
./get-eth-key --mnemonic "math razor capable expose worth grape metal sunset metal sudden usage scheme" --hdpath "m/44'/60'/0'/0/0"
```

**Output:**

```
addr: 
keyDir: 
mnemonic: math razor capable expose worth grape metal sunset metal sudden usage scheme
hdPath: m/44'/60'/0'/0/0
From mnemonic and hd path
privateKey: 0x..., address: 0x...
```

### Using Makefile

```bash
# Mode 1: Extract from keystore
make run

# Mode 2: Generate from mnemonic
make run2
```

## HD Wallet Path Format

The HD wallet path follows BIP44 standard:

```
m / purpose' / coin_type' / account' / change / index
```

- `m`: Master key
- `purpose'`: Hardened derivation, must be `44` for BIP44
- `coin_type'`: Hardened derivation, `60` for Ethereum
- `account'`: Hardened derivation, account index (typically `0`)
- `change`: `0` for external addresses, `1` for internal addresses
- `index`: Address index (typically `0` for the first address)

**Common Ethereum paths:**

- `m/44'/60'/0'/0/0`: First external address of account 0
- `m/44'/60'/0'/0/1`: Second external address of account 0
- `m/44'/60'/0'/1/0`: First internal address of account 0

## Notes

- **Mode 1**: The keystore file must be located at `{dir}/{address}` where `address` is the Ethereum address without the `0x` prefix. The password is currently set to an empty string (no password required).
- **Mode 2**: The mnemonic phrase must be valid BIP39 format. The HD path must follow BIP44 standard with purpose `44`.
- The tool outputs debug information showing all provided arguments before processing.
- If neither mode's required arguments are provided, the tool will output "[Mode] nothing to run".

## Troubleshooting

### Error: "[Mode] nothing to run"

Either provide both `-addr` and `-dir` for Mode 1, or both `-mnemonic` and `-hdpath` for Mode 2.

### Error: "private key file is not found" (Mode 1)

- Ensure the keystore file exists at `{dir}/{address}`
- Verify the address format (should match the filename without `0x` prefix)
- Check that the directory path is correct

### Error: "invalid mnemonic" (Mode 2)

- Ensure the mnemonic phrase is valid BIP39 format (12 or 24 words)
- Check for typos or missing words
- Verify the mnemonic phrase is complete

### Error: "invalid path level" or "prefix should be 'm'" (Mode 2)

- Ensure the HD path starts with `m/`
- Verify the path format: `m/44'/{coin_type}'/{account}'/{change}/{index}`
- Check that hardened derivation levels (purpose, coin_type, account) have apostrophes (`'`)

### Error: "purpose should be 44" (Mode 2)

The HD path must use BIP44 standard with purpose `44`. Other purposes are not supported.

### Error: "change should be 0 or 1" (Mode 2)

The change value must be either `0` (external addresses) or `1` (internal addresses).

## Related Files

- `main.go`: Main implementation file with command-line interface
- `wallet.go`: BIP39/BIP44 HD wallet implementation
- `Makefile`: Makefile for building and running
