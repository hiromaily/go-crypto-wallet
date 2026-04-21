# XRP Ledger (XRP) Technical Reference

This document provides a technical reference for XRP Ledger implementation in the go-crypto-wallet system.

## Table of Contents

1. [Overview](#overview)
2. [Core Specifications](#core-specifications)
3. [Address Types & Key Derivation](#address-types--key-derivation)
4. [Transaction Architecture](#transaction-architecture)
5. [Signing Mechanism](#signing-mechanism)
6. [Network Configuration](#network-configuration)
7. [Wallet Implementation](#wallet-implementation)
8. [Testing Resources](#testing-resources)
9. [Project Documentation](#project-documentation)

---

## Overview

### What is XRP Ledger?

XRP Ledger (XRPL) is a decentralized public blockchain built for payments. Unlike Bitcoin's PoW or Ethereum's PoS, XRPL uses the **XRP Ledger Consensus Protocol** (a form of federated Byzantine agreement) to achieve finality in 3–5 seconds without mining.

### Key Characteristics (2026)

| Property | Value |
|----------|-------|
| **Launch Date** | 2012 |
| **Ledger Close Time** | 3–5 seconds |
| **Consensus** | XRP Ledger Consensus Protocol (federated BFT) |
| **Native Currency** | XRP |
| **Smallest Unit** | Drop (1 XRP = 1,000,000 drops) |
| **Cryptographic Curves** | ed25519 (recommended), secp256k1 |
| **Address Format** | Classic address: `r...` (base58check) |
| **Multisig** | Native multi-signing via SignerList |

---

## Core Specifications

### Cryptographic Algorithms

| Algorithm | Description |
|-----------|-------------|
| **ed25519** | Recommended — shorter signatures, deterministic, faster verification |
| **secp256k1** | Legacy — same curve as Bitcoin |

This project defaults to **ed25519** for new key generation (`key_algorithm: "ed25519"` in config).

### Address Derivation

Classic address (`r...`) is derived from the public key via:

```
1. Generate key pair (ed25519 or secp256k1)
2. Compute SHA-512Half of the public key
3. Base58check-encode with XRPL alphabet
```

---

## Address Types & Key Derivation

### Address Format: Classic vs X-address

XRPL has two address encodings:

| Format | Prefix | Level | Status |
|--------|--------|-------|--------|
| **Classic address** | `r...` (base58check) | Protocol-native | Standard — used on-ledger |
| **X-address** | `X...` mainnet / `T...` testnet | Presentation only | Proposed (XLS-5d), not widely adopted |

**Classic addresses are the standard as of 2026.** X-addresses (proposed in [XLS-5d](https://github.com/XRPLF/XRPL-Standards/discussions/18)) are a convenience encoding that wraps a classic address together with an optional destination tag into a single string. They operate entirely at the application/presentation layer — the XRPL protocol and `xrpl-go` library use classic addresses internally.

This project uses classic addresses (`r...`). This is correct and not legacy.

#### Destination Tag (known limitation)

The problem X-addresses tried to solve is **destination tag routing**: many exchanges share a single deposit address and use a destination tag (a 32-bit integer) to identify individual users. X-addresses encode address + tag into one string to prevent users from forgetting the tag.

This project does not currently handle destination tags. This is a known feature gap for exchange-style deposit flows, but does not affect direct wallet-to-wallet transfers (P2P).

---

### HD Wallet Derivation Path

**Standard:** BIP44 with XRP coin type

```
m / 44' / 144' / account' / 0 / address_index
```

| Component | Value | Notes |
|-----------|-------|-------|
| **Purpose** | 44' | BIP44 |
| **Coin Type** | 144' | SLIP-0044 XRP |
| **Account** | See table below | |
| **Change** | 0 | External chain only |
| **Index** | 0..N | Sequential |

### Account Types

| Account | Index | Path | Use |
|---------|-------|------|-----|
| **deposit** | 0 | `m/44'/144'/0'/0/i` | Aggregate client funds |
| **payment** | 1 | `m/44'/144'/1'/0/i` | Outgoing payments |
| **stored** | 2 | `m/44'/144'/2'/0/i` | Cold storage |

---

## Transaction Architecture

### Account Model

XRP Ledger uses an **account-based model** (similar to Ethereum, unlike Bitcoin's UTXO). Each account has:

- A balance in drops
- A sequence number (incremented per transaction)
- An optional destination tag for payment routing

### Sequence Number

Every transaction requires the sender's current sequence number. Unlike Ethereum nonce, XRPL sequence must be fetched from the validated ledger state.

### Reserve Requirements

Accounts must maintain a **base reserve** (currently 10 XRP) to remain active. Sending the full balance would deactivate the account.

---

## Signing Mechanism

### Offline Signing (This Project)

XRP uses **offline signing via the Keygen wallet**. The Sign wallet is not required for XRP.

```
Watch Wallet  →  create unsigned tx  →  file
Keygen Wallet →  sign offline        →  signed tx file
Watch Wallet  →  submit signed tx    →  XRPL node
```

The signing library used is `github.com/xrpscan/xrpl-go`.

### Multi-Signature

XRPL supports native multi-signing via `SignerList`. This project currently implements single-sig only (Pattern 1).

---

## Network Configuration

### Supported Network Types

| `network_type` | Description | WebSocket URL |
|----------------|-------------|---------------|
| `mainnet` | Production network | `wss://s1.ripple.com:51234` (auto) |
| `testnet` | Public testnet (Altnet) | `wss://s.altnet.rippletest.net:51233` (auto) |
| `devnet` | Developer network (experimental amendments) | `wss://s.devnet.rippletest.net:51233` (auto) |
| `standalone` | Local isolated rippled node | Must set `websocket_public_url` in config |

For `mainnet`, `testnet`, and `devnet` the public WebSocket URL is resolved automatically from the network type.
For `standalone`, `websocket_public_url` must be set explicitly (e.g. `ws://127.0.0.1:6006`).

### Config Reference (`config/wallet/xrp/watch.yaml`)

```yaml
ripple:
  websocket_public_url: "ws://127.0.0.1:6006"  # required for standalone
  websocket_admin_url:  "ws://127.0.0.1:6006"  # required for ledger_accept in standalone
  network_type: "standalone"  # mainnet, testnet, devnet, standalone
```

---

## Wallet Implementation

### Wallet Roles

| Wallet | Role | Network |
|--------|------|---------|
| **Watch** | Create transactions, submit, monitor | Online |
| **Keygen** | Generate keys, sign transactions (offline) | Offline (air-gapped) |

> The **Sign wallet is not used for XRP**. The Keygen wallet handles offline signing directly.

### Keygen Wallet Operations

1. **`create seed`** — Generate BIP39 mnemonic and store encrypted seed
2. **`create hdkey`** — Derive BIP44 HD keys for specified account
3. **`export address`** — Export addresses to file for Watch Wallet

### Watch Wallet Operations

1. **`import address`** — Import addresses from Keygen Wallet file
2. **`create transfer`** — Build unsigned payment transaction
3. **`send`** — Submit signed transaction to XRPL node
4. **`monitor`** — Track confirmation status

---

## Testing Resources

### E2E Test Environment: rippled Standalone Mode

All E2E tests run against a **local `rippled` node in standalone mode** (Docker container).
Standalone mode is XRPL's equivalent of Bitcoin regtest:

- Fresh genesis ledger on every startup
- Ledgers do **not** close automatically — `ledger_accept` must be called after each transaction
- No external network dependency — fully deterministic
- Uses the genesis account (`rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh`) to fund test accounts

```bash
# Start rippled in standalone mode
docker compose up -d rippled

# Advance ledger (required after each transaction)
docker compose exec -T rippled rippled --silent ledger_accept
```

### Running E2E Tests

```bash
make xrp-e2e-p1           # Pattern 1: Single-sig Payment Transfer
make xrp-e2e-p1-reset     # Full reset + run
make xrp-e2e-p1-ci        # CI mode (non-interactive)
```

### E2E Transaction Patterns

| Pattern | Type | Network | Status |
|---------|------|---------|--------|
| **P1** | Single-sig Payment Transfer | rippled standalone | Verified |

### Supported Network Types by Environment

| Environment | `network_type` | Node |
|-------------|----------------|------|
| E2E / Local Dev | `standalone` | Docker rippled (local) |
| Integration testing | `testnet` | Public Altnet or self-hosted |
| Experimental features | `devnet` | XRPL Devnet |
| Production | `mainnet` | Public mainnet servers |

---

## Project Documentation

| Document | Description |
|----------|-------------|
| [testing-strategy.md](../../../../docs/chains/xrp/testing-strategy.md) | CI/E2E node strategy, standalone mode rationale |
| [setup-docker-compose-standalone-xrpl.md](../../../../docs/chains/xrp/setup-docker-compose-standalone-xrpl.md) | Docker Compose setup for local rippled |
| [transaction-flow.md](../../../../docs/chains/xrp/transaction-flow.md) | XRP transaction flow details |
| [architecture-2026.md](../../../../docs/chains/xrp/architecture-2026.md) | Current architecture overview |
| [xrpl-go.md](../../../../docs/chains/xrp/xrpl-go.md) | xrpl-go library usage |
| [network-devnet.md](../../../../docs/chains/xrp/network-devnet.md) | Devnet configuration |

---

**Document Version:** 1.0
**Last Updated:** 2026-03-06
**Maintainer:** go-crypto-wallet team
