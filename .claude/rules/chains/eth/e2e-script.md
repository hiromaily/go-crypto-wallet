---
paths: ["scripts/operation/eth/e2e/**"]
---

# ETH E2E Script Development Rules

Rules applied when creating or modifying Ethereum E2E scripts.
**Also load**: `.claude/rules/chains/e2e-script.md` for universal E2E rules (DB config, Makefile policy, security, etc.).

## Required Documentation

| Document                                               | Contents                                             |
| ------------------------------------------------------ | ---------------------------------------------------- |
| `docs/chains/eth/transaction-patterns.md`              | Ethereum transaction types and application patterns  |
| `docs/chains/eth/README.md`                            | ETH chain overview and architecture                  |
| `scripts/operation/eth/eth_common.sh`                  | ETH-specific utility functions                       |
| `scripts/operation/common.sh`                          | Common utility functions (sourced by eth_common.sh)  |
| `pkg/config/README.md`                                 | Configuration override via environment variables     |

## Script Header Template

```bash
#!/usr/bin/env bash

# Ethereum E2E Workflow Script - Pattern N: [Pattern Name]
# This script automates the complete Ethereum workflow for [description]
# Usage: ./scripts/operation/eth/e2e/e2e-pN.sh [OPTIONS]
# Options:
#   --cleanup          Stop containers and cleanup state
#   --reset            Full reset and run from scratch
#   --verbose          Enable verbose output (set -x)
#   --non-interactive  Run without prompts (for CI/CD)
#   -h, --help         Display help message
#
# Reference Documentation:
#   docs/chains/eth/transaction-patterns.md
#
# Transaction Pattern:
#   Pattern N: ETH [Pattern Name]
#   - [Key characteristics]
#   - Transaction Type: EIP-1559 (Type 2) / [other type]
#   - Signing: [Keygen wallet / Safe multisig / MPC-TSS]
#
# Environment Variables:
#   NODE_TYPE    Node type: anvil (default) or geth
#   DB_TYPE      Database type: sqlite (default) or mysql
#   ETH_RPC_HOST Ethereum RPC host (default: 127.0.0.1)
#   ETH_RPC_PORT Ethereum RPC port (default: 8546 for anvil, 8545 for geth)
```

## Script Initialization

Every script MUST:
1. Source `eth_common.sh` (which also sources `common.sh`)
2. Set `export E2E_PATTERN="pN"` immediately after sourcing
3. Call `eth_init_database` to configure DB paths

```bash
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=../eth_common.sh
source "${SCRIPT_DIR}/../eth_common.sh"

export E2E_PATTERN="pN"
eth_init_database
```

## Node Types

| NODE_TYPE             | Description                         | Docker Profile | Port | Chain ID |
| --------------------- | ----------------------------------- | -------------- | ---- | -------- |
| `anvil` (**default**) | Foundry Anvil local devnet          | `anvil`        | 8546 | 31337    |
| `geth`                | Geth `--dev` PoA chain (geth --dev) | `geth-dev`     | 8545 | 1337     |

**Port Note**: Anvil defaults to port **8546** to avoid conflict with Geth on 8545.

## Database Configuration (ETH-Specific)

Pattern-isolated SQLite paths (set automatically by `eth_init_database`):

| Wallet | Path                                        |
| ------ | ------------------------------------------- |
| Watch  | `./data/sqlite/eth/watch-e2e-{pattern}.db`  |
| Keygen | `./data/sqlite/eth/keygen-e2e-{pattern}.db` |

### Database Query Functions

```bash
# Via eth_common.sh helpers
addr=$(eth_get_payment_address "payment")

# Direct SQLite queries
sqlite3 "${SQLITE_WATCH_DB_PATH}"  "SELECT wallet_address FROM address WHERE coin='eth' LIMIT 5"
sqlite3 "${SQLITE_KEYGEN_DB_PATH}" "SELECT address FROM eth_account_key LIMIT 5"
```

| DB_TYPE  | Manual query command                                                                 |
| -------- | ------------------------------------------------------------------------------------ |
| `sqlite` | `sqlite3 ./data/sqlite/eth/{wallet}-e2e-{pattern}.db "SELECT ..."`                  |
| `mysql`  | `docker compose exec -T wallet-mysql mysql -u root -proot watch -e "SELECT ..."`    |

## Wallet Command Wrappers

Always use the wrappers from `eth_common.sh` — they inject the wallet-specific SQLite path automatically.

```bash
# ✅ CORRECT: use wrappers
eth_watch_cmd  --coin eth create tx ...
eth_keygen_cmd --coin eth key ...

# ❌ WRONG: calling binary directly (skips DB path injection)
"${GOPATH}/bin/watch" --coin eth create tx ...
```

## ETH Has No Address Type Configuration

**Unlike BTC, ETH does NOT have `address_type` or `key_type` configuration.**
All Ethereum addresses are standard secp256k1 EOA addresses (`0x...` checksummed).
Do not set `WALLET_ADDRESS_TYPE` or `WALLET_KEY_TYPE` in ETH scripts.

## Pattern-Specific Environment Variables

| Pattern | Extra Variables                                     | Description                                |
| ------- | --------------------------------------------------- | ------------------------------------------ |
| P2      | `HYT_CONTRACT_ADDRESS`, `HYT_DEPLOYER_KEY`          | ERC-20 HYT token contract and deployer key |
| P3      | `DEPLOYER_KEY`                                      | Anvil deployer key for Safe deployment     |
| P4      | `MPC_THRESHOLD`, `MPC_PARTY_IDS`, `MPC_PEER_ADDRS` | MPC-TSS ceremony parameters                |

## Transaction Patterns

| Pattern | Description                  | Signing Method   | Notes                                              |
| ------- | ---------------------------- | ---------------- | -------------------------------------------------- |
| P1      | Single-sig EIP-1559          | Keygen HD wallet | Standard EOA transfer                              |
| P2      | ERC-20 HYT token transfer    | Keygen HD wallet | Requires HYT contract deployment                   |
| P3      | Safe 2-of-2 multisig payment | EIP-712 offline  | Requires Safe contract deployment                  |
| P4      | MPC-TSS 2-of-3 threshold sig | Distributed MPC  | `keygen dkg`, `keygen serve mpc`, `watch send mpc` |

## Makefile Targets

Add targets to `make/wallet/eth_e2e.mk`. Naming convention: `eth-e2e-pN` (no `P=` parameter — pattern is in the target name).

| Target                    | Description                          |
| ------------------------- | ------------------------------------ |
| `make eth-e2e-pN-reset`   | Fresh start with reset (recommended) |
| `make eth-e2e-pN`         | Run without reset                    |
| `make eth-e2e-pN-verbose` | Run with verbose output              |
| `make eth-e2e-pN-ci`      | Run in non-interactive mode          |
| `make eth-e2e-pN-cleanup` | Cleanup only                         |
| `make eth-e2e-parallel`   | Run all patterns concurrently        |
| `make eth-e2e-ci-all`     | Run all patterns in CI mode          |

### Parallel Testing

```bash
make eth-e2e-parallel                     # all patterns (default P1-P3)
make eth-e2e-parallel PATTERNS=1,2        # specific patterns
make eth-e2e-parallel MAX_PARALLEL=2      # limit concurrency
make eth-e2e-ci-all                       # CI mode, all patterns
```

**Current Scripts**:

| Pattern | Script      | Make Target   | Status  |
| ------- | ----------- | ------------- | ------- |
| 1       | `e2e-p1.sh` | `eth-e2e-p1`  | ✅      |
| 2       | `e2e-p2.sh` | `eth-e2e-p2`  | ✅      |
| 3       | `e2e-p3.sh` | `eth-e2e-p3`  | ✅      |
| 4       | `e2e-p4.sh` | `eth-e2e-p4`  | planned |

## Common Errors

### Node Not Ready

```bash
curl -s -X POST http://127.0.0.1:8546 \
  -H "Content-Type: application/json" \
  --data '{"method":"eth_blockNumber","params":[],"id":1,"jsonrpc":"2.0"}'
```

### Wrong Port

Anvil defaults to **8546** (not 8545) to avoid conflict with Geth:
```bash
echo $ETH_RPC_PORT   # 8546 for anvil, 8545 for geth
```

### HYT Contract Not Deployed (Pattern 2)

```bash
make deploy-hyt
make eth-e2e-p2-reset
```

### Safe Not Deployed (Pattern 3)

```bash
cd apps/eth-contracts && bun install
```

### Database State from Previous Run

```bash
make eth-e2e-pN-reset   # --reset wipes and restarts
```

### MPC Party Timeout (Pattern 4)

All `keygen serve mpc` daemons must be listening **before** `watch send mpc` is called.
