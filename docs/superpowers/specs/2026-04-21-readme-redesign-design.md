# README Redesign — Design Spec

**Date:** 2026-04-21
**Status:** Approved

## Goal

Refactor `README.md` from a dense 400-line reference document into a lean landing page that hooks readers in under 30 seconds and links out to the documentation site for all detail.

## Approach

Option 3 — Story Flow. Shortest structure, fastest to read, strongest "front door" feeling. Chain coverage is summarized (not exhaustive). Quick Start is immediately actionable via E2E commands. Everything else links to the docs site.

## Key Decisions

- **Primary value proposition:** Multi-chain breadth (BTC with Taproot/MuSig2, BCH, ETH with Safe/MPC-TSS, XRP)
- **Quick Start entry point:** E2E test commands (`make btc-e2e P=1`, etc.) — shows the wallet working end-to-end immediately
- **Docs site link:** `https://hiromaily.github.io/go-crypto-wallet/` — prominent in hero section and documentation links section

## Output Structure

The new `template/pages/README.tpl.md` will include exactly 4 sections:

```
<!-- @include: ../sections/project/overview.md -->
<!-- @include: ../sections/product/chain-coverage.md -->
<!-- @include: ../sections/reference/quick-start.md -->
<!-- @include: ../sections/reference/documentation-links.md -->
```

## Section Designs

### Section 1: `template/sections/project/overview.md` (rewrite)

Content:
- Title (`# go-crypto-wallet`)
- 3 coin images aligned right (BTC, ETH, XRP) — keep existing
- Badge row: Go Report Card, Test CI, Release, License, Documentation
- Docs site link: `**[📖 Documentation Site](https://hiromaily.github.io/go-crypto-wallet/)**`
- 2-sentence tagline:
  > A production-grade multi-chain cold wallet system built in Go.
  > Offline key generation and signing, online watch-only wallet, and E2E-verified transaction patterns across BTC, BCH, ETH, and XRP.

### Section 2: `template/sections/product/chain-coverage.md` (new file)

Content:
- `## Supported Chains` heading
- Single compact table — one row per chain:

| Chain | Address Types | Highlights | Verified Patterns |
|-------|--------------|------------|-------------------|
| **BTC** | P2PKH → P2TR | Taproot (BIP341), MuSig2 (BIP327), Descriptor Wallets (BIP380), PSBT | 11 |
| **BCH** | P2PKH, P2SH | CashAddr format, SIGHASH_FORKID replay protection | 3 |
| **ETH** | EOA `0x...` | ERC-20 tokens, Safe multisig, MPC-TSS threshold signing | 4 |
| **XRP** | Classic `r...` | Ed25519/secp256k1, multisig, offline keygen | 2 |

- Link to full E2E pattern docs

### Section 3: `template/sections/reference/quick-start.md` (new file)

Content:
- `## Quick Start` heading
- One-liner prerequisite note (Docker required, link to installation guide)
- E2E commands, one block per chain — core patterns only, no CI/reset variants:

```bash
# Bitcoin (BTC) — 11 patterns
make btc-e2e P=1    # P2PKH single-sig
make btc-e2e P=9    # P2TR Taproot

# Bitcoin Cash (BCH) — 3 patterns
make bch-e2e P=1    # P2PKH single-sig
make bch-e2e P=2    # P2SH 2-of-3 multisig

# Ethereum (ETH) — 4 patterns
make eth-e2e-p1     # Single-sig EIP-1559
make eth-e2e-p3     # Safe 2-of-2 multisig
make eth-e2e-p4     # MPC-TSS 2-of-3 threshold signing

# XRP Ledger — 2 patterns
make xrp-e2e-p1     # Single-sig payment
make xrp-e2e-p2     # 2-of-2 multisig payment
```

- Link to full E2E guide for all patterns and options

### Section 4: `template/sections/reference/documentation-links.md` (new file)

Content:
- `## Documentation` heading
- Clean link list organized by topic:

| Topic | Link |
|-------|------|
| Getting Started & Installation | docs site |
| Architecture | docs site |
| CLI Command Reference | `internal/interface-adapters/cli/README.md` |
| BTC Chain Guide | `docs/chains/btc/README.md` |
| BCH Chain Guide | `docs/chains/bch/README.md` |
| ETH Chain Guide | `docs/chains/eth/README.md` |
| XRP Chain Guide | `docs/chains/xrp/README.md` |
| E2E Transaction Patterns | `docs/chains/btc/operations/e2e-transaction-patterns.md` |
| Database Management | docs site |

## What Is Removed from README

All of the following sections are dropped from `README.tpl.md`. Their section files remain intact and continue to be used by other templates (docs site pages, ARCHITECTURE.md, etc.):

- `development/requirements.md` — requirements tables
- `product/current-features.md` — development status checklist
- `development/comprehensive-e2e-testing-by-chain.md` — full E2E pattern tables
- `product/use-cases.md` — deposit/payment/transfer descriptions
- `product/wallet-types.md` — watch/keygen/sign wallet detail
- `product/workflow-diagram.md` — key generation and transaction flow diagrams
- `architecture/directory-structure.md` — full directory tree
- `architecture/components.md` — xrpl-grpc-server, eth-contracts
- `development/environment.md` — DevContainer setup
- `development/installation.md` — installation steps
- `reference/operation-example.md` — BTC operation example links
- `reference/command-example.md` — Makefile and CLI module links
- `architecture/architecture.md` — full architecture deep-dive

## Testing / Validation

After implementation:
1. `docs-ssot validate` — all includes resolve
2. `docs-ssot build` — README.md regenerates cleanly
3. `git diff README.md` — visually verify output matches design
4. Check no broken links to removed sections remain in the README output
