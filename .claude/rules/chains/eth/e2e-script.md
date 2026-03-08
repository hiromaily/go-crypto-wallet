---
paths: ["scripts/operation/eth/e2e/**"]
---

# ETH E2E Script Development Rules

Rules applied when creating or modifying Ethereum E2E scripts.

## Required Documentation

Read the following documents before creating or modifying scripts:

| Document                                               | Contents                                             |
| ------------------------------------------------------ | ---------------------------------------------------- |
| `docs/chains/eth/transaction-patterns.md`              | Ethereum transaction types and application patterns  |
| `docs/chains/eth/README.md`                            | ETH chain overview and architecture                  |
| `scripts/operation/eth/eth_common.sh`                  | ETH-specific utility functions                       |
| `scripts/operation/common.sh`                          | Common utility functions (sourced by eth_common.sh)  |
| `pkg/config/README.md`                                 | Configuration override via environment variables     |

## Script Structure Conventions

### Header Comments

Each script must include header comments in the following format:

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

### Script Initialization

Every script MUST:
1. `source "${SCRIPT_DIR}/../eth_common.sh"` (which also sources `common.sh`)
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

ETH E2E scripts support two node backends via the `NODE_TYPE` environment variable:

| NODE_TYPE           | Description                         | Docker Profile | Port |
| ------------------- | ----------------------------------- | -------------- | ---- |
| `anvil` (**default**) | Foundry Anvil local devnet         | `anvil`        | 8546 |
| `geth`              | Geth `--dev` PoA chain (chain 1337) | `geth-dev`     | 8545 |

**Port Note**: Anvil defaults to port **8546** to avoid conflict with Geth on 8545.

## Database Configuration

E2E scripts support two database backends via the `DB_TYPE` environment variable:

| DB_TYPE                | Description        | Docker MySQL | Use Case              |
| ---------------------- | ------------------ | ------------ | --------------------- |
| `sqlite` (**default**) | Local SQLite files | Not required | Fast testing, CI/CD   |
| `mysql`                | Docker MySQL       | Required     | Full integration test |

### SQLite File Paths

Pattern-isolated SQLite paths (set automatically by `eth_init_database`):

| Wallet  | Path                                              |
| ------- | ------------------------------------------------- |
| Watch   | `./data/sqlite/eth/watch-e2e-{pattern}.db`        |
| Keygen  | `./data/sqlite/eth/keygen-e2e-{pattern}.db`       |

### Usage

```bash
# Anvil + SQLite (default — fastest, no Docker MySQL needed)
make eth-e2e-pN-reset

# Geth + MySQL
make eth-e2e-pN-reset NODE_TYPE=geth DB=mysql
```

### Database Query Functions

Use the helper functions from `eth_common.sh` for DB access:

```bash
# Get an unallocated payment address from the watch DB
addr=$(eth_get_payment_address "payment")

# Direct SQLite query (when helper functions are not available)
sqlite3 "${SQLITE_WATCH_DB_PATH}" "SELECT wallet_address FROM address WHERE coin='eth' LIMIT 5"
sqlite3 "${SQLITE_KEYGEN_DB_PATH}" "SELECT address FROM eth_account_key LIMIT 5"
```

| DB_TYPE  | Command                                                                               |
| -------- | ------------------------------------------------------------------------------------- |
| `sqlite` | `sqlite3 ./data/sqlite/eth/{wallet}-e2e-{pattern}.db "SELECT ..."`                   |
| `mysql`  | `docker compose exec -T wallet-mysql mysql -u root -proot watch -e "SELECT ..."`     |

## Wallet Command Wrappers

Always use the wrappers from `eth_common.sh` instead of calling binaries directly.
The wrappers inject the wallet-specific SQLite path automatically.

```bash
# ✅ CORRECT: use wrappers
eth_watch_cmd  --coin eth create tx ...
eth_keygen_cmd --coin eth key ...

# ❌ WRONG: calling binary directly (skips DB path injection)
"${GOPATH}/bin/watch" --coin eth create tx ...
```

## Configuration File Policy

### ❌ Do NOT Edit Config Files Directly

Do **not** edit `eth/watch.yaml`, `eth/keygen.yaml`, etc.
Use **environment variables** to override settings.

### ✅ Override via Environment Variables

ETH config keys use the `WALLET_` prefix:

```bash
export WALLET_ETHEREUM_PORT="8546"       # Override RPC port
export WALLET_ETHEREUM_HOST="127.0.0.1"  # Override RPC host
```

Priority: Environment Variables > Config File > Default Values

## ETH Has No Address Type Configuration

**Unlike BTC, ETH does NOT have `address_type` or `key_type` configuration.**
All Ethereum addresses are standard secp256k1 EOA addresses (`0x...` checksummed).
There is no `WALLET_ADDRESS_TYPE` or `WALLET_KEY_TYPE` for ETH scripts.

## Pattern-Specific Environment Variables

Pattern-specific variables beyond the standard `NODE_TYPE` / `DB_TYPE`:

| Pattern | Extra Variables                                   | Description                                 |
| ------- | ------------------------------------------------- | ------------------------------------------- |
| P2      | `HYT_CONTRACT_ADDRESS`, `HYT_DEPLOYER_KEY`        | ERC-20 HYT token contract and deployer key  |
| P3      | `DEPLOYER_KEY`                                    | Anvil deployer key for Safe deployment      |
| P4      | `MPC_THRESHOLD`, `MPC_PARTY_IDS`, `MPC_PEER_ADDRS` | MPC-TSS ceremony parameters               |

## Transaction Patterns

| Pattern | Description                    | Signing Method     | Notes                                 |
| ------- | ------------------------------ | ------------------ | ------------------------------------- |
| P1      | Single-sig EIP-1559            | Keygen HD wallet   | Standard EOA transfer                 |
| P2      | ERC-20 HYT token transfer      | Keygen HD wallet   | Requires HYT contract deployment      |
| P3      | Safe 2-of-2 multisig payment   | EIP-712 (offline)  | Requires Safe contract deployment     |
| P4      | MPC-TSS 2-of-3 threshold sig   | Distributed MPC    | `keygen dkg`, `keygen serve mpc`, `watch send mpc` |

## ⚠️ MANDATORY: Always Use Makefile Targets

**AI Agents and developers MUST use Makefile targets to run E2E tests.**
Do NOT execute E2E scripts directly.

```bash
# ✅ CORRECT: Use Makefile target
make eth-e2e-p1-reset

# ❌ WRONG: Do not run scripts directly
./scripts/operation/eth/e2e/e2e-p1.sh --reset
```

### Why Makefile Targets?

1. **Automatic Build**: All targets depend on `build-all` (incremental, only rebuilds on Go source changes)
2. **Consistent Environment**: Properly passes `NODE_TYPE` and `DB_TYPE` variables
3. **Validated Targets**: Targets are defined per-pattern in `make/wallet/eth_e2e.mk`

### Available Makefile Targets

| Target                      | Description                          |
| --------------------------- | ------------------------------------ |
| `make eth-e2e-pN-reset`     | Fresh start with reset (recommended) |
| `make eth-e2e-pN`           | Run without reset                    |
| `make eth-e2e-pN-verbose`   | Run with verbose output              |
| `make eth-e2e-pN-ci`        | Run in non-interactive mode          |
| `make eth-e2e-pN-cleanup`   | Cleanup only                         |
| `make eth-e2e-parallel`     | Run all patterns concurrently        |
| `make eth-e2e-ci-all`       | Run all patterns in CI mode          |

### Parallel Testing

```bash
# Run all patterns in parallel (default: P1-P3)
make eth-e2e-parallel

# Run specific patterns
make eth-e2e-parallel PATTERNS=1,2

# Limit concurrency
make eth-e2e-parallel MAX_PARALLEL=2

# CI mode (all patterns, non-interactive)
make eth-e2e-ci-all
```

## Makefile Targets for New Scripts

Add targets to `make/wallet/eth_e2e.mk` when creating new scripts.

**Naming Convention**: `eth-e2e-pN` where N is the pattern number.

```makefile
###############################################################################
# E2E Tests
#
# Pattern N: [Description]
#   NODE_TYPE: anvil (default) or geth
#   DB:        sqlite (default) or mysql
#
# Usage:
#   make eth-e2e-pN-reset             # Run with Anvil + SQLite
#   make eth-e2e-pN-reset NODE_TYPE=geth DB=mysql
###############################################################################

.PHONY: eth-e2e-pN-reset
eth-e2e-pN-reset: build-all
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-pN.sh --reset

.PHONY: eth-e2e-pN
eth-e2e-pN: build-all
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-pN.sh

.PHONY: eth-e2e-pN-verbose
eth-e2e-pN-verbose: build-all
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-pN.sh --verbose

.PHONY: eth-e2e-pN-ci
eth-e2e-pN-ci: build-all
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-pN.sh --non-interactive

.PHONY: eth-e2e-pN-cleanup
eth-e2e-pN-cleanup:
	NODE_TYPE="$(NODE_TYPE)" DB_TYPE="$(DB)" ./scripts/operation/eth/e2e/e2e-pN.sh --cleanup
```

**Current Scripts**:

| Pattern | Script         | Make Target     | Status |
| ------- | -------------- | --------------- | ------ |
| 1       | `e2e-p1.sh`    | `eth-e2e-p1`    | ✅     |
| 2       | `e2e-p2.sh`    | `eth-e2e-p2`    | ✅     |
| 3       | `e2e-p3.sh`    | `eth-e2e-p3`    | ✅     |
| 4       | `e2e-p4.sh`    | `eth-e2e-p4`    | planned |

## Verification After Changes

### After Go Code Changes

```bash
make go-lint && make check-build && make go-test
make eth-e2e-pN-reset
```

### After Shell Script Changes

```bash
make shfmt
```

## Retry Limit

**If the fix-test cycle exceeds 5 iterations, organize progress and report.**

### Escalation Conditions

- Same error occurs repeatedly
- Deep understanding of Ethereum protocol required
- Large-scale Go code changes needed

## Common Errors

### Node Not Ready

```bash
# Verify node is running
curl -s -X POST http://127.0.0.1:8546 \
  -H "Content-Type: application/json" \
  --data '{"method":"eth_blockNumber","params":[],"id":1,"jsonrpc":"2.0"}'
```

### Wrong Port

Anvil uses **8546** by default (not 8545) to avoid conflict with Geth:
```bash
echo $ETH_RPC_PORT   # should be 8546 for anvil, 8545 for geth
```

### HYT Contract Not Deployed (Pattern 2)

```bash
make deploy-hyt
# Then re-run
make eth-e2e-p2-reset
```

### Safe Not Deployed (Pattern 3)

Ensure `apps/eth-contracts/node_modules` is installed:
```bash
cd apps/eth-contracts && bun install
```

### Database State from Previous Run

```bash
make eth-e2e-pN-reset   # use --reset to wipe and start fresh
```

### MPC Party Timeout (Pattern 4)

Ensure all party nodes are started before running `watch send mpc`.
All `keygen serve mpc` daemons must be listening before the coordinator initiates the signing session.

## Security Rules

- ❌ Do NOT log private keys or passphrases
- ❌ Do NOT use test credentials in production configs
- Reference: `docs/guidelines/security.md`

## Avoiding Impact on Other Patterns

### When Modifying eth_common.sh

- Always verify impact on all existing E2E patterns (P1–P3)
- Do not break existing pattern behavior

### Pattern-Specific Changes

- Set environment variables locally within each script (`export E2E_PATTERN=pN`)
- Each pattern uses an isolated SQLite DB via pattern suffix (`-e2e-pN`)

## Related Skills

| Skill             | Use Case                     |
| ----------------- | ---------------------------- |
| `shell-scripts`   | Script creation/modification |
| `go-development`  | Go code changes              |
| `makefile-update` | Makefile updates             |
| `git-workflow`    | Branch management            |
