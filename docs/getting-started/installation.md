# Installation

This guide covers setting up the development environment on macOS.

## Requirements

| Tool | Version | Purpose |
|------|---------|---------|
| [Go](https://go.dev/dl/) | 1.26.2+ | Build the wallet binaries |
| [Docker](https://www.docker.com/get-started) | latest | Blockchain nodes and databases |
| Docker Compose | latest | Container orchestration |
| [Foundry](https://getfoundry.sh/) | latest | ETH E2E only: deploy ERC-20 and Safe contracts (P2, P3), cast for MPC-TSS (P4) |

> **For E2E tests (recommended entry point):** Go, Docker, and Docker Compose are sufficient for BTC, BCH, and XRP. ETH E2E patterns P2, P3, and P4 additionally require Foundry.

Install Foundry (macOS/Linux):

```bash
curl -L https://foundry.paradigm.xyz | bash
foundryup
```
## Build the Wallets

Build all three wallet binaries:

```bash
make build
```

This produces `watch`, `keygen`, and `sign` binaries. The sign binary embeds the authorizer name at build time:

```bash
# Manual build (equivalent to make build)
go build -v -o ${GOPATH}/bin/watch  ./cmd/watch/
go build -v -o ${GOPATH}/bin/keygen ./cmd/keygen/
go build -ldflags "-X main.authName=auth1" -v -o ${GOPATH}/bin/sign1 ./cmd/sign/
go build -ldflags "-X main.authName=auth2" -v -o ${GOPATH}/bin/sign2 ./cmd/sign/
```

Configuration files are in [`./config/wallet/*.toml`](https://github.com/hiromaily/go-crypto-wallet/tree/main/config/wallet).
## E2E Test Setup

E2E tests are **self-contained**: each script starts and stops the required blockchain node containers automatically. You do not need to manually start nodes or configure databases before running an E2E test.

**Default database:** SQLite — no Docker database container required.

```bash
# BTC: build wallets, start regtest node, run full flow, stop node
make btc-e2e P=1

# ETH (Anvil): build wallets, start Anvil container, run full flow, stop container
make eth-e2e-p1

# XRP: build wallets, start rippled container, run full flow, stop container
make xrp-e2e-p1
```

### ETH-specific: Foundry required for P2, P3, P4

| Pattern | Reason |
|---------|--------|
| P2 (ERC-20 HYT) | `forge` deploys the HYT ERC-20 contract onto Anvil |
| P3 (Safe multisig) | `forge` deploys the Safe v1.4.1 proxy contract onto Anvil |
| P4 (MPC-TSS) | `cast` is used to verify on-chain balances |

Install Foundry before running ETH P2, P3, or P4:

```bash
curl -L https://foundry.paradigm.xyz | bash && foundryup
```

### Database options

All E2E tests default to SQLite. To use PostgreSQL or MySQL, pass `DB=`:

```bash
make btc-e2e P=1 DB=postgres   # requires Docker PostgreSQL container
make btc-e2e P=1 DB=mysql      # requires Docker MySQL container
```
## Bitcoin / Bitcoin Cash Node Setup (Manual Operation)

> **For E2E tests:** Node startup is handled automatically by the E2E scripts. This section is only needed for manual operation workflows.

Run BTC node containers via Docker Compose:

```bash
docker compose -f compose.btc.yaml up btc-watch btc-keygen btc-sign
```

Set up `bitcoin-cli` aliases:

```bash
alias bitcoin-cli-watch='docker exec -it btc-watch bitcoin-cli'
alias bitcoin-cli-keygen='docker exec -it btc-keygen bitcoin-cli'
alias bitcoin-cli-sign='docker exec -it btc-sign bitcoin-cli'
```

Create and load wallets:

```bash
# Create wallets (first time only)
./scripts/operation/create-bitcoind-wallet.sh

# Load wallets (required after container restart)
./scripts/operation/load-bitcoind-wallet.sh
```
## Ethereum Node Setup (Manual Operation)

> **For E2E tests:** Node startup is handled automatically by the E2E scripts. This section is only needed for manual operation workflows.

Two node options are supported:

### A. Anvil (default for E2E and local development)

[Anvil](https://getfoundry.sh/anvil/overview/) is part of Foundry and is the default node for all ETH E2E patterns.

```bash
docker compose -f compose.eth.yaml up anvil
```

### B. go-ethereum (Geth)

```bash
docker compose -f compose.eth.yaml up geth
# or
make up-docker-geth
```

Pass `NODE_TYPE=geth` to E2E scripts to use Geth instead of Anvil:

```bash
make eth-e2e-p1 NODE_TYPE=geth
```
## ERC-20 Contract Deployment

The HYT ERC-20 contract is deployed automatically by `make eth-e2e-p2` using Foundry (`forge`). No manual deployment step is required for E2E testing.

For manual deployment (advanced use):

```bash
cd ./apps/eth-contracts
forge build
forge script script/DeployHYT.s.sol --broadcast --rpc-url http://localhost:8546
```

Requires Foundry to be installed (`curl -L https://foundry.paradigm.xyz | bash && foundryup`).
## XRP Ledger Node Setup (Manual Operation)

> **For E2E tests:** Node startup is handled automatically by the E2E scripts. This section is only needed for manual operation workflows.

Run a local `rippled` node in standalone mode (equivalent to regtest):

```bash
docker compose -f compose.xrp.yaml up rippled
```

The standalone mode allows manual ledger advancement:

```bash
# Advance the ledger (equivalent to mining a block)
./scripts/operation/xrp/ledger-accept.sh
```
