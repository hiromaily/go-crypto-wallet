### Ethereum Node Setup (Manual Operation)

> **For E2E tests:** Node startup is handled automatically by the E2E scripts. This section is only needed for manual operation workflows.

Two node options are supported:

#### A. Anvil (default for E2E and local development)

[Anvil](https://getfoundry.sh/anvil/overview/) is part of Foundry and is the default node for all ETH E2E patterns.

```bash
docker compose -f compose.eth.yaml up anvil
```

#### B. go-ethereum (Geth)

```bash
docker compose -f compose.eth.yaml up geth
# or
make up-docker-geth
```

Pass `NODE_TYPE=geth` to E2E scripts to use Geth instead of Anvil:

```bash
make eth-e2e-p1 NODE_TYPE=geth
```
