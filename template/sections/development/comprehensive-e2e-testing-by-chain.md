## ✨ Comprehensive E2E Testing by Chain

This project includes **fully automated E2E tests** covering transaction patterns from legacy to cutting-edge. Each pattern is implemented and verified through real transactions on regtest.

### Bitcoin (BTC) Transaction Patterns

BTC supports 11 transaction patterns from legacy P2PKH to cutting-edge Taproot with MuSig2.

| Pattern | Type | Address Format | Signature | Status |
|---------|------|----------------|-----------|--------|
| **P1** | P2PKH Single-sig | `m.../n...` | 1-of-1 | ✅ Verified |
| **P2** | P2PKH 2-of-3 Multisig | `2...` | 2-of-3 | ✅ Verified |
| **P3** | P2SH-P2WPKH Single-sig | `2...` | 1-of-1 | ✅ Verified |
| **P4** | P2SH-P2WSH 2-of-3 | `2...` | 2-of-3 | ✅ Verified |
| **P5** | P2WPKH Native SegWit | `bcrt1q...` | 1-of-1 | ✅ Verified |
| **P6** | P2WSH 2-of-3 | `bcrt1q...` | 2-of-3 | ✅ Verified |
| **P7** | P2WSH 3-of-3 | `bcrt1q...` | 3-of-3 | ✅ Verified |
| **P8** | P2SH-P2WSH 3-of-3 | `2...` | 3-of-3 | ✅ Verified |
| **P9** | P2TR Taproot Single-sig | `bcrt1p...` | Schnorr | ✅ Verified |
| **P10** | P2TR MuSig2 N-of-N | `bcrt1p...` | Aggregated | 🔧 Framework |
| **P11** | P2TR Tapscript M-of-N | `bcrt1p...` | Script Path | 🔧 Framework |

### Bitcoin Cash (BCH) Transaction Patterns

BCH supports P2PKH and P2SH patterns with CashAddr format. SegWit and Taproot are **not supported** on BCH.

| Pattern | Type | Address Format | Signature | Status |
|---------|------|----------------|-----------|--------|
| **P1** | P2PKH Single-sig | `bchreg:q...` | 1-of-1 | ✅ Verified |
| **P2** | P2SH 2-of-3 Multisig | `bchreg:p...` | 2-of-3 | ✅ Verified |
| **P3** | P2SH 3-of-3 Multisig | `bchreg:p...` | 3-of-3 | ✅ Verified |

> **Note**: BCH uses CashAddr format (`bitcoincash:q...` for P2PKH, `bitcoincash:p...` for P2SH on mainnet).
> BCH signing requires SIGHASH_FORKID (0x41) for replay protection.

### Ethereum (ETH) Transaction Patterns

ETH supports EIP-1559 (Type 2) transactions with EOA single-sig, ERC-20 token transfers, Safe multisig payments, and MPC-TSS threshold signing.

| Pattern | Type | Address Format | Transaction Type | Status |
|---------|------|----------------|------------------|--------|
| **P1** | Single-sig (EOA) | `0x...` | EIP-1559 (Type 2) | ✅ Verified |
| **P2** | ERC-20 HYT Token Transfer | `0x...` | EIP-1559 (Type 2) | ✅ Verified |
| **P3** | Safe 2-of-2 Multisig Payment | `0x...` | Safe `execTransaction` | ✅ Verified |
| **P4** | MPC-TSS 2-of-3 Threshold Signing | `0x...` | EIP-1559 (Type 2) | ✅ Verified |

> **Note**: Supports both Anvil and Geth nodes. Database backend is configurable (SQLite/MySQL/PostgreSQL).
> P2 deploys the HYT ERC-20 contract via Foundry (forge) onto the local Anvil node before running the transfer workflow.
> P3 deploys a [Safe v1.4.1](https://safe.global/) proxy contract via Foundry, then exercises the full 2-of-2 multisig
> workflow: propose → offline EIP-712 sign (signer 1) → offline EIP-712 sign (signer 2) → broadcast `execTransaction`.
> P4 implements threshold ECDSA signing via MPC-TSS (tss-lib): 3 nodes run distributed key generation (DKG) offline,
> then 2-of-3 nodes collaborate online to produce a single ECDSA signature without any party holding the full private key.

### XRP Ledger (XRP) Transaction Patterns

XRP supports single-sig payment transfer with classic addresses (ed25519 or secp256k1).

| Pattern | Type | Address Format | Signing | Status |
|---------|------|----------------|---------|--------|
| **P1** | Single-sig Payment Transfer | `r...` (base58) | Keygen wallet (offline) | ✅ Verified |
| **P2** | 2-of-2 Multisig Payment Transfer | `r...` (base58) | Keygen wallet (offline, sequential) | ✅ Verified |

> **Note**: XRP uses classic addresses (`r...` prefix). The Sign wallet is not required — the
> Keygen wallet handles offline signing directly. E2E runs against a local **rippled standalone mode**
> node (Docker), equivalent to Bitcoin regtest. Ledgers are advanced manually via `ledger_accept`.
> P2 exercises the full XRP multisig workflow: set-signer-list → create-multisig-tx →
> add-multisig-signature (×2) → submit-multisig-tx.

### Quick Start

```bash
# BTC: Run any pattern with a single command
make btc-e2e P=1    # P2PKH Single-sig
make btc-e2e P=9    # P2TR Taproot

# BCH: Run supported patterns
make bch-e2e P=1    # P2PKH Single-sig
make bch-e2e P=2    # P2SH 2-of-3 Multisig
make bch-e2e P=3    # P2SH 3-of-3 Multisig

# ETH: Run individual patterns
make eth-e2e-p1     # Single-sig EIP-1559
make eth-e2e-p2     # ERC-20 HYT Token Transfer
make eth-e2e-p3     # Safe 2-of-2 Multisig Payment
make eth-e2e-p4     # MPC-TSS 2-of-3 Threshold Signing

# XRP: Run patterns
make xrp-e2e-p1     # Single-sig Payment Transfer
make xrp-e2e-p2     # 2-of-2 Multisig Payment Transfer

# ETH: Run all patterns in parallel
make eth-e2e-parallel          # Run P1, P2, P3, and P4 in parallel
make eth-e2e-ci-all            # CI mode (non-interactive)

# XRP: Run all patterns in parallel
make xrp-e2e-parallel          # Run P1 and P2 in parallel
make xrp-e2e-ci-all            # CI mode (non-interactive)

# Fresh start with full reset
make btc-e2e-reset P=5
make bch-e2e-reset P=3
make eth-e2e-p1-reset
make eth-e2e-p2-reset
make eth-e2e-p3-reset
make eth-e2e-p4-reset
make xrp-e2e-p1-reset
make xrp-e2e-p2-reset

# CI/CD mode (non-interactive)
make btc-e2e-ci P=3
make bch-e2e-ci P=3
make eth-e2e-p1-ci
make eth-e2e-p2-ci
make eth-e2e-p3-ci
make eth-e2e-p4-ci
make xrp-e2e-p1-ci
make xrp-e2e-p2-ci
```

### Why This Matters

- 🔒 **Production-Ready**: Every transaction pattern is tested end-to-end
- 🔄 **Regression Testing**: Catch breaking changes before they reach production
- 📚 **Learning Resource**: Real working examples of all Bitcoin script types
- 🚀 **CI/CD Ready**: Automated testing for continuous integration
- 🧩 **Modular Design**: Shared utilities in `btc_common.sh` reduce code duplication by 80%

See [E2E Transaction Patterns Guide](../../../docs/chains/btc/operations/e2e-transaction-patterns.md) for detailed documentation.
