## Supported Chains

| Chain | Address Types | Highlights | E2E Patterns |
|-------|--------------|------------|--------------|
| **[BTC](./docs/chains/btc/README.md)** | P2PKH → P2TR | Taproot (BIP341), MuSig2 (BIP327), Descriptor Wallets (BIP380), PSBT | [11 patterns](./docs/chains/btc/operations/e2e-transaction-patterns.md) |
| **[BCH](./docs/chains/bch/README.md)** | P2PKH, P2SH | CashAddr format, SIGHASH\_FORKID replay protection | 3 patterns |
| **[ETH](./docs/chains/eth/README.md)** | EOA `0x...` | ERC-20 tokens, Safe multisig (v1.4.1), MPC-TSS threshold signing | 4 patterns |
| **[XRP](./docs/chains/xrp/README.md)** | Classic `r...` | Ed25519 / secp256k1, multisig, offline keygen | 2 patterns |

Each pattern is implemented in Go and verified through real transactions on regtest or a local node — not mocks.
