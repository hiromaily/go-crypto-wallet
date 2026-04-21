### E2E Test Setup

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

#### ETH-specific: Foundry required for P2, P3, P4

| Pattern | Reason |
|---------|--------|
| P2 (ERC-20 HYT) | `forge` deploys the HYT ERC-20 contract onto Anvil |
| P3 (Safe multisig) | `forge` deploys the Safe v1.4.1 proxy contract onto Anvil |
| P4 (MPC-TSS) | `cast` is used to verify on-chain balances |

Install Foundry before running ETH P2, P3, or P4:

```bash
curl -L https://foundry.paradigm.xyz | bash && foundryup
```

#### Database options

All E2E tests default to SQLite. To use PostgreSQL or MySQL, pass `DB=`:

```bash
make btc-e2e P=1 DB=postgres   # requires Docker PostgreSQL container
make btc-e2e P=1 DB=mysql      # requires Docker MySQL container
```
