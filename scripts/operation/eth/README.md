# Ethereum Operation Scripts

This directory contains Ethereum operation and E2E workflow scripts.

## Directory Structure

```
eth/
├── eth_common.sh         # ETH-specific common utilities
├── README.md             # This file
└── e2e/                  # E2E test scripts
    ├── e2e-p1.sh         # Pattern 1: Single-sig EIP-1559
    ├── e2e-p2.sh         # Pattern 2: ERC-20 HYT token transfer
    └── e2e-parallel-runner.sh  # Parallel test runner for CI/CD
```

## Common Utilities

### eth_common.sh

ETH-specific common functions for E2E scripts. This file automatically sources `../common.sh`, so you don't need to source both.

**Usage in E2E scripts:**

```bash
# Source ETH common utilities
source "${SCRIPT_DIR}/../eth_common.sh"

# Initialize config paths
eth_get_config_paths

# Use common functions
eth_check_prerequisites
eth_setup_infrastructure
```

**Available Functions:**

| Function | Description |
|----------|-------------|
| `eth_get_config_paths` | Set standard config file paths |
| `eth_init_database` | Initialize database configuration |
| `eth_init_sqlite_db` | Initialize SQLite databases with E2E schemas |
| `eth_check_prerequisites` | Check CLI commands are built |
| `eth_setup_infrastructure` | Start ETH node (Anvil or Geth) |
| `eth_cleanup` | Stop containers and remove SQLite databases |
| `eth_full_reset` | Full reset including keystore files |
| `eth_watch_cmd` | Wrapper for watch commands |
| `eth_keygen_cmd` | Wrapper for keygen commands |
| `eth_sign_cmd` | Wrapper for sign commands |
| `eth_fund_address` | Fund address with ETH (Anvil or Geth) |
| `eth_get_payment_address` | Get first unallocated address from watch DB |
| `eth_export_watch_address_csv` | Export keygen addresses as watch-compatible CSV |
| `eth_extract_file_path` | Extract file path from command output |
| `eth_extract_tx_id` | Extract transaction ID from send output |

### Database Configuration

The scripts support two database backends:

#### SQLite (Default)

Uses local SQLite files. No Docker database container required.

```bash
# Run with SQLite (default)
./scripts/operation/eth/e2e/e2e-p1.sh
# or explicitly:
DB_TYPE=sqlite ./scripts/operation/eth/e2e/e2e-p1.sh
```

Benefits:
- Faster test startup (no Docker database container)
- Parallel test execution (each pattern uses an isolated DB file)
- Lighter CI/CD environments

#### MySQL

Uses Docker MySQL container. Set `DB_TYPE=mysql`.

```bash
# Run with MySQL
DB_TYPE=mysql ./scripts/operation/eth/e2e/e2e-p1.sh
```

Start the container with: `docker compose --profile mysql up`

> Note: PostgreSQL support is not yet implemented for ETH scripts.

**Environment Variables:**

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_TYPE` | `sqlite` | Database type: `sqlite` or `mysql` |
| `SQLITE_WATCH_DB_PATH` | `./data/sqlite/eth/watch-e2e-{pattern}.db` | Watch wallet database |
| `SQLITE_KEYGEN_DB_PATH` | `./data/sqlite/eth/keygen-e2e-{pattern}.db` | Keygen wallet database |

---

## ETH Node Configuration

The scripts support two Ethereum node types:

#### Anvil (Default)

Local development chain powered by Foundry Anvil. Fastest option for E2E testing.

```bash
# Run with Anvil (default)
./scripts/operation/eth/e2e/e2e-p1.sh
# or explicitly:
NODE_TYPE=anvil ./scripts/operation/eth/e2e/e2e-p1.sh
```

Start with: `docker compose -f compose.eth.yaml --profile anvil up -d anvil`

Default RPC port: `8546` (to avoid conflict with Geth)

#### Geth (Dev mode)

Local PoA development chain (`geth --dev`, chain ID 1337). Requires more startup time.

```bash
# Run with Geth dev node
NODE_TYPE=geth ./scripts/operation/eth/e2e/e2e-p1.sh
```

Start with: `docker compose -f compose.eth.yaml --profile geth-dev up -d`

Default RPC port: `8545`

**ETH Node Environment Variables:**

| Variable | Default | Description |
|----------|---------|-------------|
| `NODE_TYPE` | `anvil` | Node type: `anvil` or `geth` |
| `ETH_RPC_HOST` | `127.0.0.1` | Ethereum RPC host |
| `ETH_RPC_PORT` | `8546` (anvil) / `8545` (geth) | Ethereum RPC port |

---

## E2E Patterns

### Pattern 1: Single-sig EIP-1559 (`e2e-p1.sh`)

Automates the complete Ethereum single-sig transaction flow.

**Transaction Pattern:**
- Coin: `eth`
- Address Type: Standard secp256k1 EOA
- Transaction Type: EIP-1559 (Type 2)
- Signing: Keygen wallet (offline HD derivation, no Sign wallet needed)

**Workflow Phases:**

```
Phase 1: Key Generation (Keygen wallet — offline)
  └─ create seed → create hdkey (payment, deposit accounts)

Phase 2: Address Export → Watch Import
  └─ keygen DB → CSV file → watch import address

Phase 3: Fund Addresses (test environment only)
  └─ anvil_setBalance / eth_sendTransaction

Phase 4: Create Unsigned Transaction (Watch wallet — online)
  └─ watch create transfer → unsigned JSON file

Phase 5: Offline Signing (Keygen wallet — offline)
  └─ keygen sign signature --file → signed JSON file

Phase 6: Broadcast (Watch wallet — online)
  └─ watch send tx --file → txID

Phase 7: Confirmation Monitoring (Watch wallet — online)
  └─ watch monitor senttx
```

**Usage:**

```bash
# Run E2E workflow (default: Anvil + SQLite)
./scripts/operation/eth/e2e/e2e-p1.sh

# Fresh start with full reset
./scripts/operation/eth/e2e/e2e-p1.sh --reset

# Run with verbose output
./scripts/operation/eth/e2e/e2e-p1.sh --verbose

# Stop containers and cleanup
./scripts/operation/eth/e2e/e2e-p1.sh --cleanup
```

---

### Pattern 2: ERC-20 HYT Token Transfer (`e2e-p2.sh`)

Automates the complete HYT ERC-20 token transfer flow. Requires Foundry toolchain (`forge`, `cast`).

**Transaction Pattern:**
- Coin: `hyt` (ERC-20 on Ethereum)
- Address Type: Standard secp256k1 EOA (same HD derivation as Pattern 1)
- Transaction Type: EIP-1559 (Type 2) with ABI-encoded `transfer(address,uint256)` data field
- Signing: Keygen wallet (offline HD derivation, no Sign wallet needed)

**Workflow Phases:**

```
Phase 0: Deploy HYT Contract
  └─ forge script DeployHYT.s.sol → HYT_CONTRACT_ADDRESS

Phase 1: Key Generation (Keygen wallet — offline)
  └─ create seed → create hdkey (payment, deposit accounts)

Phase 2: Address Export → Watch Import
  └─ keygen DB → CSV file → watch import address --coin hyt

Phase 3: Fund ETH (gas) + HYT Tokens
  └─ anvil_setBalance (ETH for gas) + cast send transfer() (HYT tokens)

Phase 4: Create Unsigned Transaction (Watch wallet — online)
  └─ watch create transfer --coin hyt → unsigned JSON file

Phase 5: Offline Signing (Keygen wallet — offline)
  └─ keygen sign signature --coin hyt --file → signed JSON file

Phase 6: Broadcast (Watch wallet — online)
  └─ watch send tx --coin hyt --file → txID

Phase 7: HYT Balance Verification
  └─ cast call balanceOf() → verify recipient received tokens

Phase 8: Confirmation Monitoring (Watch wallet — online)
  └─ watch monitor senttx --coin hyt
```

**HYT-specific Environment Variables:**

| Variable | Default | Description |
|----------|---------|-------------|
| `HYT_CONTRACT_ADDRESS` | auto-detected from broadcast/ | HYT contract address |
| `HYT_DEPLOYER_KEY` | Anvil account 0 | Private key of HYT deployer/master account |
| `HYT_MASTER_ADDRESS` | Anvil account 0 address | Address of the HYT deployer |

**Usage:**

```bash
# Run E2E workflow (requires Foundry: forge, cast)
./scripts/operation/eth/e2e/e2e-p2.sh

# With specific HYT contract address
HYT_CONTRACT_ADDRESS=0x... ./scripts/operation/eth/e2e/e2e-p2.sh --reset
```

---

## Parallel Test Runner (`e2e-parallel-runner.sh`)

Runs all ETH E2E patterns in parallel using a shared Anvil node. Each pattern uses an isolated SQLite database.

```bash
# Run all patterns in parallel (CI mode)
./scripts/operation/eth/e2e/e2e-parallel-runner.sh --ci

# Run specific patterns
./scripts/operation/eth/e2e/e2e-parallel-runner.sh --patterns 1,2

# Run with verbose output
./scripts/operation/eth/e2e/e2e-parallel-runner.sh --verbose
```

**Options:**

| Option | Description |
|--------|-------------|
| `--patterns <list>` | Run specific patterns (e.g., `"1,2"` or `"1-2"`) |
| `--max-parallel <N>` | Limit concurrent processes (default: 2) |
| `--verbose` | Show real-time output from all processes |
| `--ci` | Non-interactive mode for CI/CD |
| `-h, --help` | Display help message |

---

## Configuration Files

| File | Purpose |
|------|---------|
| `config/wallet/eth/watch.yaml` | Watch wallet configuration |
| `config/wallet/eth/keygen.yaml` | Keygen wallet configuration |
| `config/wallet/eth/sign.yaml` | Sign wallet configuration |

## Generated Files

### Address Files (`data/address/eth/`)

- `payment_*.csv` - Payment account addresses (ETH or HYT)
- `deposit_*.csv` - Deposit account addresses
- `*_hyt_*.csv` - HYT-specific address CSV files

### Transaction Files (`data/tx/eth/`)

- `*.json` - Unsigned and signed transaction files (ETH JSON format)

### SQLite Database Files (`data/sqlite/eth/`)

- `watch-e2e-p1.db` - Watch wallet database for Pattern 1
- `keygen-e2e-p1.db` - Keygen wallet database for Pattern 1
- `watch-e2e-p2.db` - Watch wallet database for Pattern 2
- `keygen-e2e-p2.db` - Keygen wallet database for Pattern 2

## Account Types

| Account | Purpose |
|---------|---------|
| `payment` | Sender addresses for outgoing transactions |
| `deposit` | Receiver addresses for incoming funds |

## Wallet Architecture

ETH uses only two wallets (no separate sign wallets):

| Wallet | Role | Environment |
|--------|------|-------------|
| Keygen | Key generation + offline signing | Offline |
| Watch | Transaction creation + broadcast + monitoring | Online |

The keygen wallet performs offline HD key derivation to sign transactions — no `sign1`/`sign2` wallets are needed for ETH single-sig.

## Troubleshooting

### Anvil not starting

```bash
docker compose -f compose.eth.yaml --profile anvil up -d anvil
docker logs anvil
```

### HYT contract not found (Pattern 2)

Ensure Foundry is installed and HYT has been deployed:

```bash
make deploy-hyt
# or manually:
cd apps/eth-contracts && forge script script/DeployHYT.s.sol --rpc-url http://127.0.0.1:8546 --broadcast
```

### Build CLI binaries

```bash
make build-all
```
