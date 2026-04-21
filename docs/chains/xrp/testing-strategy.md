<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/docs/chains/xrp/testing-strategy.tpl.md · Run `make docs` to regenerate.
-->

# XRPL Development & Testing Strategy (2026 Edition)

## Context

* Language: **Golang**
* Library: `github.com/xrpscan/xrpl-go`
* Responsibilities:

  * Create transactions
  * Sign transactions
  * Submit transactions
  * Confirm validation
  * Provide CLI-based E2E testing in CI

This document defines the optimal node strategy and validation architecture for development and CI in 2026.

---

# Core Principle

> CI must be deterministic and free from external network dependencies.
> Validation must be explicitly confirmed, not assumed.

---

# Recommended Node Strategy by Environment

| Environment                         | Node Type                                  | Reason                                |
| ----------------------------------- | ------------------------------------------ | ------------------------------------- |
| Unit Tests                          | No node                                    | Pure logic testing                    |
| CI E2E                              | `rippled` Standalone                       | Deterministic, no external dependency |
| Local Dev                           | Standalone                                 | Fast iteration                        |
| Integration (real network behavior) | Self-hosted `rippled` connected to Testnet | Stability & realism                   |
| Experimental features               | Specific Devnet                            | Amendment testing                     |

---

# Why Standalone Mode Is Ideal for CI

Standalone mode provides:

* Fresh genesis ledger on every startup
* Fully deterministic state
* Manual ledger control
* No faucet dependency
* No rate limits
* No network instability
* No external failures

This makes it ideal for:

* CLI-based E2E tests
* Automated CI pipelines
* Reliable submission testing
* Balance and sequence verification

---

# Standalone Mode Behavior

## Start rippled in Standalone Mode

```bash
# Port 6006 is the WebSocket admin port used by the wallet CLI
docker run --rm -d --name rippled_ci -p 6006:6006 rippleci/rippled standalone
```

Or with Docker Compose (recommended):

```bash
docker compose -f compose.xrp.yaml up -d rippled
```

---

## Genesis Account (Pre-funded)

When a new genesis ledger is created, it contains a master account:

```
Address: rHb9CJAWyB4rj91VRWn96DkukG4bwdtyTh
Secret:  snoPBrXtMeMyMHUVTgbuqAfg1SUTb
```

This account holds all XRP.

It should be used ONLY in CI or standalone environments to fund test accounts.

Never use this in production.

---

## Ledger Advancement (Critical)

In Standalone mode:

* Ledgers do NOT close automatically.
* Transactions are not finalized until ledger closes.

You must explicitly advance the ledger:

```bash
rippled ledger_accept
```

This is required to:

* Move transactions into validated state
* Update balances
* Increment sequence numbers

---

# Required CI E2E Flow

The CI pipeline must follow this deterministic sequence:

```
1. Start rippled in standalone mode
2. Generate test wallet
3. Fund wallet from genesis account
4. ledger_accept
5. Execute Go CLI transaction
6. ledger_accept
7. Verify transaction result
```

---

# Reliable Submission Rule (Mandatory)

The AI Agent must follow this rule:

> A successful `submit` response does NOT mean the transaction succeeded.

Correct validation procedure:

1. Submit transaction
2. Advance ledger (if standalone)
3. Query transaction by hash
4. Confirm:

   * `validated == true`
   * `meta.TransactionResult == tesSUCCESS`

Only then is the transaction considered successful.

---

# Environment Configuration Requirement

The XRPL endpoint must be environment-configurable via:

```
WALLET_RIPPLE_WEBSOCKET_PUBLIC_URL
```

Examples:

```
ws://localhost:6006                          # Local standalone rippled (WebSocket admin port)
wss://s.altnet.rippletest.net:51233          # Testnet public WebSocket
```

> **Note**: The wallet CLI connects **directly** to rippled via WebSocket.
> No xrpl-grpc-server or intermediate server is required (that approach is deprecated).

This allows the same CLI binary to operate in:

* CI (Standalone)
* Local development
* Testnet integration
* Devnet testing

Hardcoding endpoints is prohibited.

---

# Network Layer Abstraction Requirement

The Go implementation must:

* Accept a WebSocket endpoint as configuration
* Gracefully handle connection failures
* Implement retry logic where appropriate
* Not assume persistent WebSocket availability

Because `xrpscan/xrpl-go` is a low-level WebSocket client, reliable submission must be implemented explicitly.

---

# CI Determinism Requirements

The CI environment must:

* Start with a fresh genesis ledger
* Not use persistent volumes
* Destroy containers after execution
* Avoid public XRPL servers
* Avoid Testnet faucets

Each run must be fully isolated.

---

# When to Use Testnet

Testnet should be used for:

* Observing real fee dynamics
* Validating network latency behavior
* Testing real ledger close timing
* Observing amendment activation states

Testnet should NOT be used for:

* CI
* Deterministic balance tests
* High-frequency automated testing

For stable integration testing:

> Run your own `rippled` connected to Testnet
> Avoid public shared nodes (rate limits, instability)

---

# When to Use Devnet

Devnet is required only when testing:

* Hooks
* Lending
* Batch transactions
* Experimental amendments

Normal payment and account operations do not require Devnet.

---

# Validator Mode Clarification

Validator mode is NOT required.

Application development requires only:

```
stock server mode
```

Validator mode is for:

* UNL participation
* Consensus voting
* Validation signing

It provides no advantage for transaction submission testing.

---

# AI Agent Operational Rules

The AI Agent must:

1. Never assume submit == success
2. Always confirm `validated == true`
3. Always confirm `tesSUCCESS`
4. Use `ledger_accept` in Standalone
5. Keep XRPL endpoint configurable via `WALLET_RIPPLE_WEBSOCKET_PUBLIC_URL`
6. Avoid public infrastructure in CI
7. Ensure every test starts from clean genesis
8. Separate logical tests from network tests

---

# Final Architecture Summary

## Development

```
Docker Standalone rippled
```

## CI

```
Docker Standalone rippled
Manual ledger control
CLI-driven E2E
```

## Integration Testing

```
Self-hosted rippled connected to Testnet
```

## Experimental Features

```
Specific Devnet
```

---

# Strategic Conclusion

As of 2026, the most robust XRPL development model is:

* Standalone for CI and deterministic testing
* Testnet for real network behavior validation
* Devnet only when necessary
* No validator required
* Explicit reliable submission verification

This architecture ensures:

* Stability
* Determinism
* CI reliability
* Scalable integration
* AI-agent-friendly clarity

---

If needed, the next step can be:

* A fully production-grade `e2e.sh` template
* A reliable submission implementation example in Go
* A CI workflow example (GitHub Actions, etc.)
* An AI Agent system prompt template for autonomous XRPL development
