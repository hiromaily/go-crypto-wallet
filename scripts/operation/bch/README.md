# Bitcoin Cash Operation Scripts

This directory contains Bitcoin Cash operation and E2E workflow scripts.

## Directory Structure

```
bch/
├── bch_common.sh         # BCH-specific common utilities
├── README.md             # This file
└── e2e/                  # E2E test scripts
    ├── e2e-p1-p2pkh-singlesig.sh  # Pattern 1: P2PKH Single-sig
    ├── e2e-p2-p2sh-2of3.sh        # Pattern 2: P2SH 2-of-3 Multisig
    ├── e2e-p3-p2sh-3of3.sh        # Pattern 3: P2SH 3-of-3 Multisig
    └── e2e-parallel-runner.sh     # Parallel test runner for CI/CD
```

## BCH Protocol Limitations

Bitcoin Cash does **NOT** support the following Bitcoin features:

| Feature | BCH Support | Notes |
|---------|-------------|-------|
| SegWit (P2WPKH, P2WSH) | No | P2PKH and P2SH only |
| Taproot (P2TR) | No | |
| Descriptor wallets | No | Cannot use descriptor APIs |
| PSBT format | No | Raw transaction hex only |
| Schnorr / MuSig2 | No | ECDSA only |
| BIP49/84/86 derivation | No | BIP44 only (coin type 145) |

BCH uses **CashAddr** format (`bitcoincash:q...`) for mainnet addresses.

## Common Utilities

### bch_common.sh

BCH-specific common functions for E2E scripts. This file automatically sources `../common.sh`, so you don't need to source both.

**Usage in E2E scripts:**

```bash
# Source BCH common utilities
source "${SCRIPT_DIR}/../bch_common.sh"

# Initialize config paths
bch_get_config_paths

# Initialize RPC hosts (required for parallel execution)
bch_init_rpc_hosts
```

**Available Functions:**

| Function | Description |
|----------|-------------|
| `bch_get_config_paths` | Set standard config file paths |
| `bch_init_rpc_hosts` | Set wallet-specific RPC host addresses |
| `bch_check_prerequisites` | Check Docker and CLI commands |
| `bch_setup_infrastructure` | Start database and BCH nodes |
| `bch_setup_wallets` | Create wallets in BCH nodes |
| `bch_cleanup` | Stop containers |
| `bch_full_reset` | Full reset with volume deletion |
| `bch_watch_cmd` | Wrapper for watch commands |
| `bch_keygen_cmd` | Wrapper for keygen commands |
| `bch_sign1_cmd` | Wrapper for sign1 commands |
| `bch_sign2_cmd` | Wrapper for sign2 commands |
| `bch_extract_file_path` | Extract file path from command output |
| `sqlite_init_db` | Initialize SQLite database with schema |
| `sqlite_clean_db` | Remove SQLite database files |
| `sqlite_query` | Execute SQLite query directly |
| `mysql_query` | Execute MySQL query via Docker |
| `db_query` | Execute database query (abstraction) |
| `db_execute` | Execute database command (abstraction) |
| `db_is_sqlite` | Check if using SQLite database |
| `db_is_mysql` | Check if using MySQL database |
| `bch_get_sender_address` | Get sender address from watch DB |
| `bch_generate_receiver_addresses` | Generate receiver addresses from BCH node |
| `bch_insert_payment_requests` | Insert payment requests directly into DB |

### Database Configuration

The scripts support two database backends:

#### SQLite (Default)

Uses local SQLite files. No Docker database container required.

```bash
# Run with SQLite (default)
make bch-e2e-reset P=1
# or explicitly:
DB_TYPE=sqlite ./scripts/operation/bch/e2e/e2e-p1-p2pkh-singlesig.sh --reset
```

Benefits:
- Faster test startup (no Docker database container)
- Parallel test execution (each pattern uses isolated timestamped DB files)
- Lighter CI/CD environments

#### MySQL

Uses Docker MySQL container. Set `DB_TYPE=mysql`.

```bash
# Run with MySQL
DB_TYPE=mysql make bch-e2e-reset P=1
```

Start the container with: `docker compose --profile mysql up`

**Environment Variables:**

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_TYPE` | `sqlite` | Database type: `sqlite` or `mysql` |
| `SQLITE_DB_DIR` | `./data/sqlite/bch` | SQLite database directory |
| `SQLITE_WATCH_DB_PATH` | auto-generated | Watch wallet database |
| `SQLITE_KEYGEN_DB_PATH` | auto-generated | Keygen wallet database |
| `SQLITE_SIGN_DB_PATH` | auto-generated | Sign1 wallet database |
| `SQLITE_SIGN2_DB_PATH` | auto-generated | Sign2 wallet database |

SQLite DB file naming: `{wallet}-e2e-{pattern}-{timestamp}.db`
(e.g., `watch-e2e-p1-20260101-120000.db`)

---

## E2E Patterns

BCH supports **3 patterns** only (no SegWit, Taproot, or descriptor patterns).

### Pattern 1: P2PKH Single-sig (`e2e-p1-p2pkh-singlesig.sh`)

**Transaction Pattern:**

- Address Type: P2PKH (BIP44)
- Address Format: `bitcoincash:q...` (mainnet) / `m.../n...` (regtest)
- Key Derivation: `m/44'/145'/account'/change/index`
- Signature: Single-sig (Keygen wallet only, no Sign wallets needed)

**Workflow Phases:**

```
Phase 1: Key Generation (Keygen wallet — offline)
  └─ create seed → create hdkey (client, deposit, payment, stored)

Phase 2: Address Export → Watch Import
  └─ keygen export address --account → CSV → watch import address

Phase 3: Test UTXO Generation (regtest only)
  └─ generatetoaddress 101 blocks → wait for balance

Phase 4: Payment Request Setup
  └─ bch_insert_payment_requests (direct DB insert via bch_common.sh helper)

Phase 5: Create Unsigned Transaction (Watch wallet — online)
  └─ watch create payment → unsigned hex file

Phase 6: Offline Signing (Keygen wallet — offline)
  └─ keygen sign signature --file → signed hex file

Phase 7: Broadcast (Watch wallet — online)
  └─ watch send tx --file → txID
```

### Pattern 2: P2SH 2-of-3 Multisig (`e2e-p2-p2sh-2of3.sh`)

**Transaction Pattern:**

- Address Type: P2SH (BIP44 + BIP11)
- Address Format: `bitcoincash:p...` (mainnet) / `2...` (regtest)
- Signature: 2-of-3 (any 2 of Keygen, Sign1, Sign2)

**Workflow Phases:**

Same as Pattern 1 plus multisig setup:

```
Phase 1: Key Generation (Keygen + Sign1 + Sign2 wallets)
  ├─ keygen: create seed → create hdkey → import privkey
  ├─ sign1:  create seed → create hdkey → import privkey
  └─ sign2:  create seed → create hdkey → import privkey

Phase 2: Multisig Setup
  ├─ sign1 export fullpubkey → fullpubkey_auth1.csv
  ├─ sign2 export fullpubkey → fullpubkey_auth2.csv
  ├─ keygen import fullpubkey (sign1, sign2)
  └─ keygen create multisig address → export to watch

Phase 3-7: Same as Pattern 1 (UTXO, create tx, sign, broadcast)
  └─ Signing: keygen sign → sign1 sign → watch send (2-of-3)
```

### Pattern 3: P2SH 3-of-3 Multisig (`e2e-p3-p2sh-3of3.sh`)

Same as Pattern 2 but requires **all 3** signatures (keygen + sign1 + sign2).

---

## Running E2E Tests

### Using Makefile (Recommended)

Always use Makefile targets. They automatically build binaries before running.

```bash
# Run with full reset (recommended for fresh start)
make bch-e2e-reset P=1    # Pattern 1
make bch-e2e-reset P=2    # Pattern 2
make bch-e2e-reset P=3    # Pattern 3

# Run without reset (continuing from previous state)
make bch-e2e P=1

# Run with verbose output
make bch-e2e-verbose P=1

# Run in non-interactive CI mode
make bch-e2e-ci P=1

# Cleanup only
make bch-e2e-cleanup P=1
```

### Parallel Test Runner

Runs all patterns in parallel using isolated SQLite databases.

```bash
# Run all patterns in parallel
make bch-e2e-parallel

# Run specific patterns
make bch-e2e-parallel PATTERNS=1,2

# Run all in CI mode
make bch-e2e-ci-all
```

Or directly:

```bash
./scripts/operation/bch/e2e/e2e-parallel-runner.sh --ci
./scripts/operation/bch/e2e/e2e-parallel-runner.sh --patterns 1,2 --verbose
```

**Parallel runner options:**

| Option | Description |
|--------|-------------|
| `--patterns <list>` | Patterns to run (e.g., `"1,2,3"` or `"1-3"`) |
| `--max-parallel <N>` | Limit concurrent processes (default: 3) |
| `--verbose` | Show real-time output |
| `--ci` | Non-interactive CI/CD mode |
| `-h, --help` | Display help |

---

## Configuration Files

| File | Purpose |
|------|---------|
| `config/wallet/bch/watch.yaml` | Watch wallet configuration |
| `config/wallet/bch/keygen.yaml` | Keygen wallet configuration |
| `config/wallet/bch/sign1.yaml` | Sign1 wallet configuration |
| `config/wallet/bch/sign2.yaml` | Sign2 wallet configuration |
| `config/wallet/account/account.yaml` | Account config (single-sig, Pattern 1) |
| `config/wallet/account/account_2of3.yaml` | Account config (2-of-3, Pattern 2) |
| `config/wallet/account/account_3of3.yaml` | Account config (3-of-3, Pattern 3) |

## Generated Files

### Address Files (`data/address/bch/`)

- `client_*.csv` - Client addresses (non-multisig)
- `deposit_*.csv` - Deposit addresses
- `payment_*.csv` - Payment addresses
- `stored_*.csv` - Stored addresses

### Public Key Files (`data/fullpubkey/bch/`)

- `auth1_*.csv` - Full public keys from sign1 (multisig patterns)
- `auth2_*.csv` - Full public keys from sign2 (multisig patterns)

### Transaction Files (`data/tx/bch/`)

- `*.hex` - Unsigned and signed raw transaction hex files

### SQLite Database Files (`data/sqlite/bch/`)

- `watch-e2e-p1-{timestamp}.db` - Watch wallet for Pattern 1
- `keygen-e2e-p1-{timestamp}.db` - Keygen wallet for Pattern 1
- (similar files per pattern)

## Account Types

| Account | Purpose | Multisig |
|---------|---------|----------|
| `client` | Client-facing addresses for receiving funds | No |
| `deposit` | Deposit addresses for initial fund receipt | Yes (P2/P3) |
| `payment` | Payment addresses for outgoing transactions | Yes (P2/P3) |
| `stored` | Cold storage addresses for long-term holding | Yes (P2/P3) |

## Docker Container Names

| Container | Port | Purpose |
|-----------|------|---------|
| `bch-watch` | 28332 | Watch wallet node |
| `bch-keygen` | 29332 | Keygen wallet node |
| `bch-sign1` | 30332 | Sign1 wallet node |
| `bch-sign2` | 31332 | Sign2 wallet node |

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_TYPE` | `sqlite` | Database type: `sqlite` or `mysql` |
| `BCH_RPC_USER` | `xyz` | BCH node RPC username (regtest) |
| `BCH_RPC_PASSWORD` | `xyz` | BCH node RPC password (regtest) |
| `BCH_WALLET_PASSPHRASE` | `test` | Wallet passphrase (if encrypted) |
| `BCH_ENCRYPTED` | `false` | Enable wallet encryption |

## E2E Script Verification Status

### Pattern 1: P2PKH Single-sig (`e2e-p1-p2pkh-singlesig.sh`)

**Status**: Under investigation — UTXO query issue

**Known Issue**: Transaction creation fails with "No utxo" error despite UTXOs existing in the watch wallet. Requires investigation of watch wallet UTXO query logic for BCH watch-only addresses.

### Pattern 2: P2SH 2-of-3 (`e2e-p2-p2sh-2of3.sh`)

**Status**: Not yet verified

### Pattern 3: P2SH 3-of-3 (`e2e-p3-p2sh-3of3.sh`)

**Status**: Not yet verified

## Troubleshooting

### "No UTXOs available"

```bash
# Check balance in watch wallet
docker exec bch-watch bitcoin-cli -regtest -rpcwallet=watch getbalance "*" 1 true

# List UTXOs
docker exec bch-watch bitcoin-cli -regtest -rpcwallet=watch listunspent

# Generate 101 blocks manually
docker exec bch-watch bitcoin-cli -regtest generatetoaddress 101 <address>
```

### Address Format Issues

BCH uses different formats per network:

| Network | P2PKH Format | P2SH Format |
|---------|-------------|-------------|
| Mainnet | `bitcoincash:q...` | `bitcoincash:p...` |
| Testnet | `bchtest:q...` | `bchtest:p...` |
| Regtest | `m.../n...` | `2...` |

Ensure `address_type = "bch-cashaddr"` is set in config files.

### Container Health Issues

```bash
docker logs bch-watch
make bch-e2e-cleanup P=1
make bch-e2e-reset P=1
```

### Build CLI Binaries

```bash
make build-all
```
