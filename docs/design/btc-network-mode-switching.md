# Design Document: Bitcoin Network Mode Switching (regtest ↔ signet)

## Overview

This document defines the design for enabling the Bitcoin E2E testing infrastructure to support both `regtest` and `signet` network modes, with easy switching between them via flags, environment variables, or configuration files.

## Background

### Current State

- E2E tests currently run exclusively in `regtest` mode
- Configuration is hardcoded for regtest in multiple places:
  - `compose.btc.yaml` - Docker compose configuration
  - `docker/nodes/btc/*/bitcoin.conf` - Node configuration files
  - `config/wallet/btc/*.yaml` - Wallet configuration files
  - `scripts/operation/btc/btc_common.sh` - E2E test scripts

### Goals

1. Support both `regtest` and `signet` modes for E2E testing
2. Enable easy mode switching via environment variables or flags
3. Maintain backward compatibility (default to `regtest`)
4. Minimize code duplication between modes
5. Support future network modes (testnet3, mainnet) with minimal changes

### Benefits of Signet Support

| Aspect | Regtest | Signet |
|--------|---------|--------|
| Block Generation | Local (instant) | Network (controlled, ~10 min) |
| Faucet Required | No (can mine) | Yes |
| Network Conditions | Artificial | Near-production |
| External Dependencies | None | Signet network |
| Testing Speed | Fast | Slower |
| Realism | Low | High |

## Technical Specifications

### Network Configuration Differences

| Parameter | Regtest | Signet |
|-----------|---------|--------|
| Default RPC Port | 18443 | 38332 |
| bitcoin.conf flag | `regtest=1` | `signet=1` |
| Config section | `[regtest]` | `[signet]` |
| Chain magic bytes | 0xFABFB5DA | 0x0A03CF40 |
| Address prefix (bech32) | `bcrt1` | `tb1` |
| Block generation | Local mining | Network sync |
| Wallet data directory | `regtest/` | `signet/` |

## Design

### Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                     Mode Selection Layer                         │
│  (Environment Variable / CLI Flag / Config File)                 │
│                                                                   │
│  BTC_NETWORK_MODE=regtest|signet                                 │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                  Configuration Provider                          │
│                                                                   │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐      │
│  │ Docker       │    │ bitcoin.conf │    │ wallet.yaml  │      │
│  │ Compose      │    │ (per node)   │    │ (per wallet) │      │
│  └──────────────┘    └──────────────┘    └──────────────┘      │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Infrastructure Layer                           │
│                                                                   │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐      │
│  │ btc-watch    │    │ btc-keygen   │    │ btc-sign*    │      │
│  │ container    │    │ container    │    │ container    │      │
│  └──────────────┘    └──────────────┘    └──────────────┘      │
└─────────────────────────────────────────────────────────────────┘
```

### Mode Selection Priority

Mode is determined in the following priority order:

1. **CLI Flag** (`--network=signet`) - Highest priority
2. **Environment Variable** (`BTC_NETWORK_MODE=signet`)
3. **Default Value** (`regtest`) - Lowest priority

### File Structure Changes

```
docker/
└── nodes/
    └── btc/
        ├── watch/
        │   ├── bitcoin.regtest.conf     # NEW: regtest-specific config
        │   ├── bitcoin.signet.conf      # NEW: signet-specific config
        │   └── bitcoin.conf             # Symlink or copied from mode-specific
        ├── keygen/
        │   ├── bitcoin.regtest.conf
        │   ├── bitcoin.signet.conf
        │   └── bitcoin.conf
        ├── sign1/
        │   ├── bitcoin.regtest.conf
        │   ├── bitcoin.signet.conf
        │   └── bitcoin.conf
        └── sign2/
            ├── bitcoin.regtest.conf
            ├── bitcoin.signet.conf
            └── bitcoin.conf

config/
└── wallet/
    └── btc/
        ├── watch.yaml                   # Uses WALLET_BITCOIN_NETWORK_TYPE env override
        ├── keygen.yaml
        ├── sign1.yaml
        └── sign2.yaml

scripts/
└── operation/
    └── btc/
        ├── btc_common.sh               # Mode-aware utility functions
        ├── setup-network-mode.sh       # NEW: Network mode setup script
        └── e2e/
            └── *.sh                    # E2E scripts with mode support

compose.btc.yaml                        # Mode-aware compose file (uses env vars)
compose.btc.regtest.yaml               # NEW: regtest-specific overrides (optional)
compose.btc.signet.yaml                # NEW: signet-specific overrides (optional)
```

### Environment Variables

| Variable | Values | Default | Description |
|----------|--------|---------|-------------|
| `BTC_NETWORK_MODE` | `regtest`, `signet` | `regtest` | Primary mode selector |
| `BTC_RPC_PORT` | Number | Mode-dependent | RPC port override |
| `WALLET_BITCOIN_NETWORK_TYPE` | `regtest`, `signet` | From BTC_NETWORK_MODE | Wallet config override |

### Docker Compose Changes

```yaml
# compose.btc.yaml (updated)
services:
  btc-watch:
    image: bitcoin/bitcoin:29.2
    container_name: btc-watch
    volumes:
      - ./docker/nodes/btc/watch:/home/bitcoin/.bitcoin
    ports:
      # Dynamic port mapping based on mode
      - "${BTC_WATCH_RPC_PORT:-18332}:${BTC_INTERNAL_RPC_PORT:-18443}"
    environment:
      - BTC_NETWORK_MODE=${BTC_NETWORK_MODE:-regtest}
    healthcheck:
      test:
        [
          "CMD",
          "bitcoin-cli",
          "-${BTC_NETWORK_MODE:-regtest}",
          "-rpcuser=xyz",
          "-rpcpassword=xyz",
          "getblockchaininfo",
        ]
    command: -printtoconsole -deprecatedrpc=create_bdb
```

### Shell Script Changes (btc_common.sh)

```bash
# Mode configuration
BTC_NETWORK_MODE="${BTC_NETWORK_MODE:-regtest}"

# Mode-specific defaults
case "${BTC_NETWORK_MODE}" in
  regtest)
    BTC_INTERNAL_RPC_PORT="${BTC_INTERNAL_RPC_PORT:-18443}"
    BTC_CLI_FLAG="-regtest"
    BTC_WALLET_DATA_DIR="regtest"
    BTC_CAN_MINE=true
    ;;
  signet)
    BTC_INTERNAL_RPC_PORT="${BTC_INTERNAL_RPC_PORT:-38332}"
    BTC_CLI_FLAG="-signet"
    BTC_WALLET_DATA_DIR="signet"
    BTC_CAN_MINE=false
    ;;
  *)
    log_error "Invalid BTC_NETWORK_MODE: ${BTC_NETWORK_MODE}"
    exit 1
    ;;
esac
```

### E2E Script Changes

#### Mode-Specific Block Generation

```bash
# In E2E scripts
generate_test_utxos() {
  local address="$1"

  if [ "${BTC_CAN_MINE}" = "true" ]; then
    # Regtest: Generate blocks locally
    btc_cli "btc-watch" generatetoaddress 101 "$address"
  else
    # Signet: Request from faucet and wait for confirmation
    log_info "Signet mode: Please fund address via faucet: $address"
    log_info "Faucet URL: https://signetfaucet.com/"
    btc_wait_for_utxo "$address" 600  # Wait up to 10 minutes
  fi
}
```

#### Mode Setup Function

```bash
# New function in btc_common.sh
btc_setup_network_mode() {
  local mode="${BTC_NETWORK_MODE:-regtest}"
  log_step "Setting up Bitcoin network mode: ${mode}"

  # Copy mode-specific bitcoin.conf to active config
  for node in watch keygen sign1 sign2; do
    local conf_dir="docker/nodes/btc/${node}"
    local mode_conf="${conf_dir}/bitcoin.${mode}.conf"
    local active_conf="${conf_dir}/bitcoin.conf"

    if [ -f "$mode_conf" ]; then
      cp "$mode_conf" "$active_conf"
      log_info "Configured ${node} for ${mode} mode"
    fi
  done

  # Export mode-specific environment variables
  export WALLET_BITCOIN_NETWORK_TYPE="${mode}"
  export BTC_CLI_FLAG="-${mode}"
}
```

### Wallet Configuration Updates

The wallet configuration files will use environment variable overrides:

```yaml
# config/wallet/btc/watch.yaml
bitcoin:
  host: "127.0.0.1:18332/wallet/watch"  # Default port
  network_type: "regtest"  # Can be overridden by WALLET_BITCOIN_NETWORK_TYPE
```

The Go code already supports environment variable overrides with the `WALLET_` prefix.

### Makefile Integration

```makefile
# make/btc_e2e.mk additions

# Network mode (default: regtest)
BTC_NETWORK_MODE ?= regtest

# Export for child processes
export BTC_NETWORK_MODE

# Mode-specific targets
.PHONY: btc-e2e-regtest
btc-e2e-regtest:
	BTC_NETWORK_MODE=regtest $(E2E_SCRIPT_PATH) --reset

.PHONY: btc-e2e-signet
btc-e2e-signet:
	BTC_NETWORK_MODE=signet $(E2E_SCRIPT_PATH) --reset

# Run E2E in current mode
.PHONY: btc-e2e-mode
btc-e2e-mode: _btc-e2e-validate
	$(E2E_SCRIPT_PATH)
```

## Implementation Plan

### Phase 1: Infrastructure Preparation (3-4 days)

1. Create mode-specific `bitcoin.conf` files for each node
2. Update `compose.btc.yaml` with environment variable support
3. Add mode detection and setup functions to `btc_common.sh`

### Phase 2: E2E Script Refactoring (2-3 days)

1. Update E2E scripts to use mode-aware utility functions
2. Implement signet-specific UTXO generation (faucet integration)
3. Add mode-specific wait times and timeouts

### Phase 3: Wallet Configuration (1-2 days)

1. Verify environment variable overrides work for network_type
2. Update wallet config documentation
3. Add validation for mode consistency

### Phase 4: Makefile and CI Integration (2-3 days)

1. Add mode-specific Make targets
2. Update CI workflows for optional signet testing
3. Create mode setup/teardown scripts

### Phase 5: Documentation and Testing (2-3 days)

1. Update E2E README with mode switching instructions
2. Document signet faucet usage
3. Test all E2E patterns in both modes
4. Create troubleshooting guide for signet issues

## Considerations

### Signet-Specific Challenges

1. **Faucet Dependency**: Signet requires external faucet for test funds
2. **Wait Times**: Blocks take ~10 minutes, requiring longer timeouts
3. **Network Availability**: Tests depend on public signet availability
4. **Initial Sync**: First-time setup requires blockchain sync

### Recommended Solutions

| Challenge | Solution |
|-----------|----------|
| Faucet dependency | Cache testnet coins, auto-request with rate limiting |
| Long wait times | Increase timeouts, add progress indicators |
| Network issues | Graceful degradation, skip signet in CI by default |
| Initial sync | Document expected sync time, provide pre-synced volumes |

### Security Considerations

- Signet configurations should NOT be used with real funds
- Clearly label signet-mode outputs to prevent confusion
- Maintain separate data directories for each mode

## Backward Compatibility

- Default mode remains `regtest` for all existing scripts
- No changes required for existing E2E test invocations
- Environment variable override is optional

## Future Extensibility

This design supports adding additional network modes:

```bash
case "${BTC_NETWORK_MODE}" in
  regtest) ... ;;
  signet) ... ;;
  testnet3)  # Future support
    BTC_INTERNAL_RPC_PORT="${BTC_INTERNAL_RPC_PORT:-18332}"
    BTC_CLI_FLAG="-testnet"
    BTC_WALLET_DATA_DIR="testnet3"
    BTC_CAN_MINE=false
    ;;
esac
```

## References

- [Bitcoin Core signet documentation](https://en.bitcoin.it/wiki/Signet)
- [Signet faucet](https://signetfaucet.com/)
- [Bitcoin network ports](https://en.bitcoin.it/wiki/Running_Bitcoin)
- [Project E2E documentation](../crypto/btc/operations/e2e-transaction-patterns.md)
