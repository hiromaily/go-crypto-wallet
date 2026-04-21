# Docker Compose setup to run a Standalone XRPL `rippled` node

## 1. `compose.yaml` (Standalone `rippled`)

```yaml
services:
  rippled:
    image: rippleci/rippled:3.1
    command: ["standalone"]

    # Expose WebSocket admin port (6006) to the host.
    # Port 6006 = ws:// WebSocket admin (used by the wallet CLI for direct connection)
    # Port 51233 = peer protocol (not needed for wallet CLI)
    ports:
      - "6006:6006"

    # Healthcheck ensures CI waits until rippled is ready.
    healthcheck:
      test: ["CMD", "rippled", "--silent", "server_info"]
      interval: 3s
      timeout: 2s
      retries: 30
```

Why this works:

* The image explicitly documents that passing `standalone` starts Standalone mode.
* In Standalone mode you typically need to manually advance the ledger (see `ledger_accept`) for tx to become validated.

---

## 2. Start / Stop commands (Compose only)

### Start

```bash
docker compose --profile mysql up -d
```

### Confirm it’s healthy

```bash
docker compose ps
```

### Stop (clean)

```bash
docker compose down -v
```

> For deterministic CI runs, `down -v` helps ensure a clean state each time.

---

## 3. The 3 commands you’ll use most in CI E2E

### A) `server_info` (readiness / debugging)

```bash
docker compose exec -T rippled rippled --silent server_info
```

### B) `ledger_accept` (advance ledger; required for validation in Standalone)

```bash
docker compose exec -T rippled rippled --silent ledger_accept
```

`ledger_accept` is the standard way to manually advance the ledger in Standalone testing workflows.

### C) `wallet_propose` (quick test wallet generation)

```bash
docker compose exec -T rippled rippled --silent wallet_propose
```

---

## 4. Connecting your Go app / CLI

Set the endpoint to the mapped port:

Set `WALLET_RIPPLE_WEBSOCKET_PUBLIC_URL` to the mapped port:

* If your Go CLI runs on the host:

  * `ws://localhost:6006`

* If your Go CLI runs as another Compose service:

  * `ws://rippled:6006` (service-to-service DNS)

---

## 5. Standalone-specific gotcha (must-do)

In Standalone:

* ledgers **don’t close automatically**
* “submit succeeded” does **not** mean “validated”

So your E2E flow should always include:

1. submit tx
2. `ledger_accept`
3. query tx by hash and confirm `validated == true` (and `tesSUCCESS`)

XRPL docs explicitly describe Standalone as a fresh genesis known state, and remind you that you must advance the ledger manually.
