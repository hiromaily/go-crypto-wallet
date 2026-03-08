---
paths: ["scripts/operation/xrp/e2e/**"]
---

# XRP E2E Script Development Rules

Rules applied when creating or modifying XRP Ledger E2E scripts.
**Also load**: `.claude/rules/chains/e2e-script.md` for universal E2E rules (DB config, Makefile policy, security, etc.).

## Required Documentation

| Document                                                     | Contents                                            |
| ------------------------------------------------------------ | --------------------------------------------------- |
| `docs/chains/xrp/README.md`                                  | XRP chain overview, architecture, wallet design     |
| `docs/chains/xrp/transaction-flow.md`                        | XRP transaction lifecycle and signing flow          |
| `docs/chains/xrp/setup-docker-compose-standalone-xrpl.md`   | rippled standalone mode setup                       |
| `scripts/operation/xrp/xrp_common.sh`                       | XRP-specific utility functions                      |
| `scripts/operation/common.sh`                                | Common utility functions (sourced by xrp_common.sh) |
| `pkg/config/README.md`                                       | Configuration override via environment variables    |

## Script Header Template

```bash
#!/usr/bin/env bash

# XRP E2E Workflow Script - Pattern N: [Pattern Name]
# This script automates the complete XRP workflow for [description]
# Usage: ./scripts/operation/xrp/e2e/e2e-pN.sh [OPTIONS]
# Options:
#   --cleanup          Stop containers and cleanup state
#   --reset            Full reset and run from scratch
#   --verbose          Enable verbose output (set -x)
#   --non-interactive  Run without prompts (for CI/CD)
#   -h, --help         Display help message
#
# Transaction Pattern:
#   Pattern N: XRP [Pattern Name]
#   - Address Type: XRP classic address (ed25519 or secp256k1)
#   - Address Format: r... (base58)
#   - Signing: [Keygen wallet / Multisig with signer list]
#
# Environment Variables:
#   DB_TYPE      Database type: sqlite (default) or mysql
#   XRP_WS_HOST  rippled WebSocket host (default: 127.0.0.1)
#   XRP_WS_PORT  rippled WebSocket port (default: 6006)
#
# Infrastructure:
#   rippled v3.1.0 standalone mode (compose.xrp.yaml)
#   WebSocket admin port 6006 (ws://) — direct connection, no xrpl-grpc-server
#   Ledger must be manually advanced via 'rippled ledger_accept' after each tx
```

## Script Initialization

Every script MUST:
1. Source `xrp_common.sh` (which also sources `common.sh`)
2. Set `export E2E_PATTERN="pN"` immediately after sourcing
3. Call `xrp_init_database` to configure DB paths
4. Call `xrp_get_config_paths` to set config file variables

```bash
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=../xrp_common.sh
source "${SCRIPT_DIR}/../xrp_common.sh"

export E2E_PATTERN="pN"
xrp_init_database
xrp_get_config_paths
```

## Infrastructure: rippled Standalone Mode

XRP uses a single `rippled` node in **standalone mode** — no peer discovery, no real network.

| Property          | Value                                  |
| ----------------- | -------------------------------------- |
| Docker compose    | `compose.xrp.yaml`                     |
| WebSocket port    | 6006 (`ws://`) — admin + public        |
| HTTP admin port   | 5005 (localhost inside container only) |
| Connection style  | Direct WebSocket, **no xrpl-grpc-server** |
| Ledger advance    | Manual — `xrp_ledger_accept` after each tx |

**No `NODE_TYPE` variable** — there is only one XRP node backend (rippled standalone).

### Ledger Acceptance

In standalone mode the ledger does NOT close automatically. After every transaction submission call:

```bash
xrp_ledger_accept   # advances ledger; defined in xrp_common.sh
```

Forgetting `ledger_accept` is the most common cause of "transaction not found" errors.

## Environment Variables

| Variable      | Default       | Description                              |
| ------------- | ------------- | ---------------------------------------- |
| `DB_TYPE`     | `sqlite`      | Database backend (`sqlite` or `mysql`)   |
| `XRP_WS_HOST` | `127.0.0.1`   | rippled WebSocket host                   |
| `XRP_WS_PORT` | `6006`        | rippled WebSocket port                   |

### Configuration Overrides

Config keys use the `WALLET_` prefix:

```bash
export WALLET_RIPPLE_WEBSOCKET_PUBLIC_URL="ws://${XRP_WS_HOST}:${XRP_WS_PORT}"
export WALLET_RIPPLE_WEBSOCKET_ADMIN_URL="ws://${XRP_WS_HOST}:${XRP_WS_PORT}"
```

The wrappers `xrp_watch_cmd` and `xrp_keygen_cmd` inject these automatically.

## Database Configuration (XRP-Specific)

Pattern-isolated SQLite paths (set automatically by `xrp_init_database`):

| Wallet | Path                                        |
| ------ | ------------------------------------------- |
| Watch  | `./data/sqlite/xrp/watch-e2e-{pattern}.db`  |
| Keygen | `./data/sqlite/xrp/keygen-e2e-{pattern}.db` |

### Database Query Functions

```bash
# Via xrp_common.sh helpers
addr=$(xrp_get_address "payment")

# Direct SQLite queries
sqlite3 "${SQLITE_WATCH_DB_PATH}"  "SELECT wallet_address FROM address WHERE coin='xrp' LIMIT 5"
sqlite3 "${SQLITE_KEYGEN_DB_PATH}" "SELECT address FROM xrp_account_key LIMIT 5"
```

| DB_TYPE  | Manual query command                                                                 |
| -------- | ------------------------------------------------------------------------------------ |
| `sqlite` | `sqlite3 ./data/sqlite/xrp/{wallet}-e2e-{pattern}.db "SELECT ..."`                  |
| `mysql`  | `docker compose exec -T wallet-mysql mysql -u root -proot watch -e "SELECT ..."`    |

## Wallet Command Wrappers

Always use the wrappers from `xrp_common.sh` — they inject both the SQLite path and the WebSocket URL automatically.

```bash
# ✅ CORRECT: use wrappers
xrp_watch_cmd  --coin xrp create tx ...
xrp_keygen_cmd --coin xrp key ...

# ❌ WRONG: calling binary directly (skips DB and WebSocket injection)
"${GOPATH}/bin/watch" --coin xrp create tx ...
```

## XRP Has No Address Type Configuration

XRP uses a single address type — classic address (`r...`, base58).
Do **not** set `WALLET_ADDRESS_TYPE` or `WALLET_KEY_TYPE` in XRP scripts.

Key type (ed25519 vs secp256k1) is configured in `keygen.yaml` but does not affect E2E script structure.

## XRP-Specific Concepts

### Account Reserve

XRP accounts do not exist until funded with at least **10 XRP** (base reserve).
Always fund new addresses before any transaction.

### Funding via Genesis Account

In standalone mode, use the well-known genesis account for funding:

```bash
# Via xrp_common.sh helper
xrp_fund_address "rXXX..." 100   # fund with 100 XRP

# Genesis credentials (standalone ONLY — never use in production)
# Address: rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh
# Secret:  snoPBrXtMeMyMHUVTgbuqAfg1SUTb
```

### XRP Amount Units

| Unit   | Value              | Usage                       |
| ------ | ------------------ | --------------------------- |
| XRP    | 1 XRP              | Human-readable amounts      |
| drops  | 1 XRP = 1,000,000  | Wire protocol, API calls    |

### Destination Tag

Payments to exchanges require a Destination Tag. Include in E2E scripts where relevant.

### SignerList for Multisig

Before sending a multisig transaction, the sender must have a SignerList on-ledger:
1. Create a `SignerListSet` transaction
2. Sign and submit it (requires `ledger_accept`)
3. Only then can multisig `Payment` transactions be submitted

## Transaction Patterns

| Pattern | Description               | Signing                           | Status                     |
| ------- | ------------------------- | --------------------------------- | -------------------------- |
| P1      | Single-sig payment        | Keygen HD wallet (offline)        | ✅ Full E2E                |
| P2      | 2-of-2 multisig payment   | Keygen as both signers (workaround) | ⚠️ Partial (see note below) |

### P2 Known Limitation

The sign wallet binary does not yet support XRP coin (`NewSigner` returns "not implemented").
As a workaround, P2 uses the keygen wallet to simulate both signers.
Full E2E requires:
1. Enable XRP in the sign wallet (`NewSigner` in `container.go`)
2. Update the keygen sign use case to detect and handle JSON-format multisig files

## Configuration Files

| File                            | Purpose                     |
| ------------------------------- | --------------------------- |
| `config/wallet/xrp/watch.yaml`  | Watch wallet configuration  |
| `config/wallet/xrp/keygen.yaml` | Keygen wallet configuration |

## Makefile Targets

Add targets to `make/wallet/xrp_e2e.mk`. Naming convention: `xrp-e2e-pN` (no `P=` parameter — pattern is in the target name).

Targets pass `DB_TYPE`, `XRP_WS_HOST`, and `XRP_WS_PORT`:

```makefile
.PHONY: xrp-e2e-pN-reset
xrp-e2e-pN-reset: build-all
	DB_TYPE="$(DB)" XRP_WS_HOST="$(XRP_WS_HOST)" XRP_WS_PORT="$(XRP_WS_PORT)" ./scripts/operation/xrp/e2e/e2e-pN.sh --reset
```

| Target                    | Description                          |
| ------------------------- | ------------------------------------ |
| `make xrp-e2e-pN-reset`   | Fresh start with reset (recommended) |
| `make xrp-e2e-pN`         | Run without reset                    |
| `make xrp-e2e-pN-verbose` | Run with verbose output              |
| `make xrp-e2e-pN-ci`      | Run in non-interactive mode          |
| `make xrp-e2e-pN-cleanup` | Cleanup only                         |
| `make xrp-e2e-parallel`   | Run all patterns concurrently        |
| `make xrp-e2e-ci-all`     | Run all patterns in CI mode          |

### Parallel Testing

XRP uses distinct parameter names to avoid Makefile conflicts with ETH parallel targets:

```bash
make xrp-e2e-parallel                             # all patterns (default P1)
make xrp-e2e-parallel XRP_PATTERNS=1,2            # specific patterns
make xrp-e2e-parallel XRP_MAX_PARALLEL=2          # limit concurrency
make xrp-e2e-ci-all                               # CI mode, all patterns
```

**Note**: A single shared rippled node is used across all patterns (unlike ETH which isolates nodes per pattern).

**Current Scripts**:

| Pattern | Script      | Make Target   | Status                     |
| ------- | ----------- | ------------- | -------------------------- |
| 1       | `e2e-p1.sh` | `xrp-e2e-p1`  | ✅ Full E2E                |
| 2       | `e2e-p2.sh` | `xrp-e2e-p2`  | ⚠️ Partial (signing workaround) |

## Common Errors

### Transaction Not Confirmed

Most likely cause: forgot to call `ledger_accept` after submitting the transaction.

```bash
# Manually advance ledger
docker compose -f compose.xrp.yaml exec -T rippled rippled --silent ledger_accept
```

### Account Not Found / "actNotFound"

The address has not been funded yet. All XRP accounts need minimum 10 XRP before use.

```bash
xrp_fund_address "rXXX..." 100
```

### WebSocket Connection Refused

```bash
# Verify rippled is running
docker compose -f compose.xrp.yaml ps
# Check WebSocket port
echo $XRP_WS_PORT   # should be 6006
```

### SignerList Not Set (Pattern 2)

The multisig sender must have a SignerList on-ledger before sending a multisig Payment.
Submit and confirm the `SignerListSet` transaction first.

### Database State from Previous Run

```bash
make xrp-e2e-pN-reset   # --reset wipes and restarts
```
