### XRP Ledger Node Setup (Manual Operation)

> **For E2E tests:** Node startup is handled automatically by the E2E scripts. This section is only needed for manual operation workflows.

Run a local `rippled` node in standalone mode (equivalent to regtest):

```bash
docker compose -f compose.xrp.yaml up rippled
```

The standalone mode allows manual ledger advancement:

```bash
# Advance the ledger (equivalent to mining a block)
./scripts/operation/xrp/ledger-accept.sh
```
