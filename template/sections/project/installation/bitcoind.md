### Bitcoin / Bitcoin Cash Node Setup (Manual Operation)

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
