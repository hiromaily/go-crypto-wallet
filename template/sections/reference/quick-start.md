## Quick Start

Requires Docker. See the [Installation Guide](./docs/getting-started/installation.md) for full setup.

```bash
# Bitcoin (BTC) — 11 patterns from P2PKH to Taproot
make btc-e2e P=1    # P2PKH single-sig
make btc-e2e P=9    # P2TR Taproot (Schnorr)

# Bitcoin Cash (BCH)
make bch-e2e P=1    # P2PKH single-sig
make bch-e2e P=2    # P2SH 2-of-3 multisig

# Ethereum (ETH)
make eth-e2e-p1     # Single-sig EIP-1559
make eth-e2e-p3     # Safe 2-of-2 multisig
make eth-e2e-p4     # MPC-TSS 2-of-3 threshold signing

# XRP Ledger
make xrp-e2e-p1     # Single-sig payment
make xrp-e2e-p2     # 2-of-2 multisig payment
```

See the [E2E Transaction Patterns Guide](./docs/chains/btc/operations/e2e-transaction-patterns.md) for all patterns, CI mode (`make btc-e2e-ci P=1`), and reset options.
