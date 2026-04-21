# README Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the 400-line README with a lean landing page (4 sections) that leads with multi-chain breadth, shows E2E quick-start commands, and links to the docs site for all detail.

**Architecture:** Rewrite `template/sections/project/overview.md`, create three new section files (`chain-coverage.md`, `quick-start.md`, `documentation-links.md`), update `template/pages/README.tpl.md` to include only those four sections, then rebuild with `docs-ssot build`.

**Tech Stack:** docs-ssot (Markdown template engine), GitHub Flavored Markdown

---

## File Map

| Action | File |
|--------|------|
| Modify | `template/sections/project/overview.md` |
| Create | `template/sections/product/chain-coverage.md` |
| Create | `template/sections/reference/quick-start.md` |
| Create | `template/sections/reference/documentation-links.md` |
| Modify | `template/pages/README.tpl.md` |
| Generated (do not edit) | `README.md` |

---

## Task 1: Rewrite `template/sections/project/overview.md`

**Files:**
- Modify: `template/sections/project/overview.md`

- [ ] **Step 1: Write the new overview section**

Replace the entire file with:

```markdown
<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/README.tpl.md · Run `make docs` to regenerate.
-->

<p align="center">
  <img src="https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/bitcoin-img.svg?sanitize=true" alt="Bitcoin" width="140px">
  <img src="https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/ethereum-img.png?raw=true" alt="Ethereum" width="140px">
  <img src="https://raw.githubusercontent.com/hiromaily/go-crypto-wallet/main/images/xrp-img.jpg?raw=true" alt="XRP" width="140px">
</p>

# go-crypto-wallet

**A production-grade multi-chain cold wallet system built in Go.**

[![Go Report Card](https://goreportcard.com/badge/github.com/hiromaily/go-crypto-wallet)](https://goreportcard.com/report/github.com/hiromaily/go-crypto-wallet)
[![Test](https://github.com/hiromaily/go-crypto-wallet/actions/workflows/lint-test.yml/badge.svg)](https://github.com/hiromaily/go-crypto-wallet/actions/workflows/lint-test.yml)
[![GitHub release](https://img.shields.io/badge/release-v6.2.0-blue.svg)](https://github.com/hiromaily/go-crypto-wallet/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Documentation](https://img.shields.io/badge/docs-GitHub%20Pages-blue)](https://hiromaily.github.io/go-crypto-wallet/)

**[📖 Documentation Site](https://hiromaily.github.io/go-crypto-wallet/)**

---

Managing cryptocurrency at scale means two things in conflict: keys must stay offline to be secure, and wallets must stay online to be useful. Most implementations compromise one for the other.

go-crypto-wallet resolves this with a **three-wallet architecture**: an offline **Keygen** wallet that generates and holds private keys, an offline **Sign** wallet held by independent authorizers for multisig, and an online **Watch** wallet that creates and broadcasts transactions without ever touching a private key.

```text
┌──────────────────┐      ┌──────────────────┐
│  Keygen Wallet   │      │   Sign Wallet    │
│   (OFFLINE)      │      │   (OFFLINE)      │
│                  │      │                  │
│  HD key gen      │      │  Auth signing    │
│  Multisig addrs  │      │  2nd+ signature  │
└────────┬─────────┘      └────────┬─────────┘
         │  export pubkeys / sign  │
         └────────────┬────────────┘
                      │
             ┌────────▼─────────┐
             │   Watch Wallet   │
             │    (ONLINE)      │
             │                  │
             │  Create tx       │
             │  Broadcast tx    │
             │  Monitor         │
             └──────────────────┘
```

Every transaction pattern — from legacy P2PKH to Taproot MuSig2 and MPC-TSS — is implemented in Go and verified end-to-end against real local nodes.

---
```

- [ ] **Step 2: Validate includes**

```bash
cd /Users/hiroki.yasui/go/src/github.com/hiromaily/go-crypto-wallet
docs-ssot validate
```

Expected: `OK`

- [ ] **Step 3: Commit**

```bash
git add template/sections/project/overview.md
git commit -m "docs: rewrite overview section for 2026 README redesign"
```

---

## Task 2: Create `template/sections/product/chain-coverage.md`

**Files:**
- Create: `template/sections/product/chain-coverage.md`

- [ ] **Step 1: Create the chain coverage section**

```markdown
## Supported Chains

| Chain | Address Types | Highlights | E2E Patterns |
|-------|--------------|------------|--------------|
| **[BTC](./docs/chains/btc/README.md)** | P2PKH → P2TR | Taproot (BIP341), MuSig2 (BIP327), Descriptor Wallets (BIP380), PSBT | [11 patterns](./docs/chains/btc/operations/e2e-transaction-patterns.md) |
| **[BCH](./docs/chains/bch/README.md)** | P2PKH, P2SH | CashAddr format, SIGHASH\_FORKID replay protection | 3 patterns |
| **[ETH](./docs/chains/eth/README.md)** | EOA `0x...` | ERC-20 tokens, Safe multisig (v1.4.1), MPC-TSS threshold signing | 4 patterns |
| **[XRP](./docs/chains/xrp/README.md)** | Classic `r...` | Ed25519 / secp256k1, multisig, offline keygen | 2 patterns |

Each pattern is implemented in Go and verified through real transactions on regtest or a local node — not mocks.
```

- [ ] **Step 2: Validate**

```bash
docs-ssot validate
```

Expected: `OK`

- [ ] **Step 3: Commit**

```bash
git add template/sections/product/chain-coverage.md
git commit -m "docs: add chain-coverage section for README redesign"
```

---

## Task 3: Create `template/sections/reference/quick-start.md`

**Files:**
- Create: `template/sections/reference/quick-start.md`

- [ ] **Step 1: Create the quick-start section**

```markdown
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
```

- [ ] **Step 2: Validate**

```bash
docs-ssot validate
```

Expected: `OK`

- [ ] **Step 3: Commit**

```bash
git add template/sections/reference/quick-start.md
git commit -m "docs: add quick-start section for README redesign"
```

---

## Task 4: Create `template/sections/reference/documentation-links.md`

**Files:**
- Create: `template/sections/reference/documentation-links.md`

- [ ] **Step 1: Create the documentation links section**

```markdown
## Documentation

| Topic | Link |
|-------|------|
| Getting Started | [Installation & Setup](./docs/getting-started/installation.md) |
| Architecture | [System Design](https://hiromaily.github.io/go-crypto-wallet/architecture/) |
| CLI Reference | [Command × Chain × UseCase matrix](./internal/interface-adapters/cli/README.md) |
| BTC Guide | [Bitcoin](./docs/chains/btc/README.md) |
| BCH Guide | [Bitcoin Cash](./docs/chains/bch/README.md) |
| ETH Guide | [Ethereum](./docs/chains/eth/README.md) |
| XRP Guide | [XRP Ledger](./docs/chains/xrp/README.md) |
| E2E Patterns | [Transaction Pattern Guide](./docs/chains/btc/operations/e2e-transaction-patterns.md) |
| Full Docs | [Documentation Site](https://hiromaily.github.io/go-crypto-wallet/) |
```

- [ ] **Step 2: Validate**

```bash
docs-ssot validate
```

Expected: `OK`

- [ ] **Step 3: Commit**

```bash
git add template/sections/reference/documentation-links.md
git commit -m "docs: add documentation-links section for README redesign"
```

---

## Task 5: Update `template/pages/README.tpl.md`

**Files:**
- Modify: `template/pages/README.tpl.md`

- [ ] **Step 1: Replace the template with the 4-section structure**

Replace the entire file with:

```markdown
<!-- @include: ../sections/project/overview.md -->

<!-- @include: ../sections/product/chain-coverage.md -->

<!-- @include: ../sections/reference/quick-start.md -->

<!-- @include: ../sections/reference/documentation-links.md -->
```

- [ ] **Step 2: Validate all includes resolve**

```bash
docs-ssot validate
```

Expected: `OK`

- [ ] **Step 3: Build and inspect output**

```bash
docs-ssot build
git diff README.md
```

Expected: README.md is regenerated. Visually verify:
- Auto-generated comment at top
- Center-aligned images
- Bold tagline
- ASCII diagram
- Narrative prose with `---` separators
- Chain coverage table
- Quick start code block
- Documentation link table

- [ ] **Step 4: Commit**

```bash
git add template/pages/README.tpl.md README.md template/INDEX.md
git commit -m "docs: update README template and regenerate lean landing page"
```

---

## Self-Review

**Spec coverage check:**

| Spec requirement | Task |
|-----------------|------|
| Auto-generated comment at top | Task 1 |
| Center-aligned hero images | Task 1 |
| Bold one-liner tagline | Task 1 |
| ASCII security model diagram | Task 1 |
| Narrative prose with `---` separators | Task 1 |
| Chain coverage table (BTC/BCH/ETH/XRP) | Task 2 |
| E2E quick-start commands by chain | Task 3 |
| Documentation link grid | Task 4 |
| README.tpl.md reduced to 4 includes | Task 5 |
| `docs-ssot build` produces clean output | Task 5 |

**Placeholder scan:** None — all section file content is written in full in each task.

**No broken references:** Removed sections (requirements, wallet-types, directory-structure, etc.) still exist as section files and continue to be used by other templates. Only `README.tpl.md` no longer includes them.
