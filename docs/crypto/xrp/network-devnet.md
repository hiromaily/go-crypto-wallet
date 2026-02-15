Below is a **Devnet-focused design document** written for an AI Agent.

# XRPL Devnet Usage Guide (2026 Edition)

## Purpose

This document defines how an AI Agent should understand and use **XRPL Devnet** within a Golang-based transaction system using `xrpscan/xrpl-go`.

Devnet is **not** a replacement for Standalone or Testnet.
It serves a specific purpose: **testing features not yet active on Mainnet/Testnet.**

---

## 1. What Is XRPL Devnet?

Devnet is a public XRPL network that:

* Runs newer or experimental amendments
* Enables features not yet active on Mainnet
* Is unstable by design
* May reset periodically
* May change behavior without notice

It is intended for:

* Protocol feature experimentation
* Amendment validation testing
* Pre-mainnet feature integration

It is NOT intended for:

* Deterministic CI
* Stable integration testing
* Production-like uptime expectations

---

## 2. When the AI Agent Should Use Devnet

The AI Agent should use Devnet only if:

* The feature being tested requires a specific amendment
* The amendment is not active on Testnet or Mainnet
* The system depends on experimental XRPL functionality

Examples of Devnet-only features (depending on current network state):

* Hooks
* Lending protocol features
* Batch transactions
* Newly proposed amendments
* Advanced AMM functionality (if not activated elsewhere)

If no amendment dependency exists, Devnet should NOT be used.

---

## 3. Devnet vs Other Environments

| Environment | Deterministic | External Dependency | Amendment Coverage | CI Suitable |
| ----------- | ------------- | ------------------- | ------------------ | ----------- |
| Standalone  | Yes           | No                  | Manual control     | Yes         |
| Testnet     | No            | Yes                 | Stable amendments  | Limited     |
| Devnet      | No            | Yes                 | Experimental       | No          |
| Mainnet     | No            | Yes                 | Fully activated    | No          |

Key Rule:

> Devnet is for feature capability validation, not behavioral reliability validation.

---

## 4. Devnet Operational Characteristics

The AI Agent must assume:

* Ledgers close automatically
* Fees may fluctuate
* Network may restart
* Faucet may be unreliable
* Historical data may be wiped
* Amendments may toggle between resets

Devnet must be treated as **ephemeral infrastructure**.

---

## 5. Devnet Connection Strategy

Devnet must be accessed via WebSocket endpoint.

The endpoint must be configurable:

```
XRPL_ENDPOINT
```

Example (placeholder):

```
wss://s.devnet.rippletest.net:51233
```

Hardcoding endpoints is prohibited.

The AI Agent must:

* Retry connection attempts
* Handle network interruptions
* Not assume ledger timing consistency

---

## 6. Devnet Testing Model

Devnet testing must follow this structure:

```
1. Connect to Devnet
2. Create or fund account via faucet
3. Submit transaction
4. Poll transaction by hash
5. Wait for validated == true
6. Confirm meta.TransactionResult == tesSUCCESS
```

Unlike Standalone:

* Do NOT use ledger_accept
* Do NOT assume ledger timing
* Implement timeout logic

---

## 7. Reliable Submission on Devnet

The AI Agent must implement:

### Submission Flow

1. Submit transaction
2. Record transaction hash
3. Poll `tx` endpoint
4. Wait for:

   * validated == true
   * meta.TransactionResult == tesSUCCESS

### Required Safeguards

* Timeout mechanism
* Retry mechanism
* Failure classification:

  * `tef*` (final failure)
  * `tel*` (local failure)
  * `tem*` (malformed)
  * `ter*` (retryable)

Submit success alone is insufficient.

---

## 8. Funding Model

Devnet requires a faucet for new accounts.

The AI Agent must assume:

* Faucet may rate-limit
* Faucet may fail
* Funding may take time to reflect

Best practice:

* Cache funded test wallets for reuse
* Avoid creating excessive new accounts

Devnet funding should NOT be part of deterministic CI.

---

## 9. Amendment Awareness Requirement

Devnet exists primarily for amendment testing.

The AI Agent must:

* Query `server_info`
* Inspect `amendment_blocked`
* Confirm required amendment is enabled

If an amendment is missing:

* Abort feature test
* Report unsupported environment

---

## 10. Devnet Is Not CI-Compatible

Devnet must NOT be used for:

* Automated deterministic CI
* Balance-sensitive assertions
* Performance baselines
* High-frequency test runs

Reasons:

* Unpredictable resets
* Fee changes
* Ledger timing variation
* External dependency risk

---

## 11. Recommended Usage Pattern

### Development Workflow

```
Standalone → Deterministic logic testing
Testnet → Real network behavior validation
Devnet → Experimental amendment validation
```

Devnet is the final validation step before Mainnet integration if the feature depends on new protocol capabilities.

---

## 12. AI Agent Operational Rules for Devnet

The AI Agent must:

1. Never assume stability
2. Never use Devnet in CI
3. Always confirm validated == true
4. Always confirm tesSUCCESS
5. Always handle network failure
6. Implement timeout and retry
7. Confirm amendment activation before feature usage
8. Avoid excessive faucet requests

---

## 13. Risk Model

Devnet risks include:

* Ledger resets
* Amendment toggling
* Faucet downtime
* Node instability
* Transaction expiration

The AI Agent must classify Devnet failures as:

```
Environment instability
OR
Feature incompatibility
```

Not automatically as application bugs.

---

## 14. Final Devnet Strategy Summary

Devnet exists to answer one question:

> "Does this feature work under the new amendment?"

It does NOT answer:

* Is my submission logic deterministic?
* Is my balance logic correct?
* Is my retry logic stable?
* Is my application production-ready?

Those questions are answered by:

* Standalone (determinism)
* Testnet (network realism)

---

## Final Architecture Including Devnet

### CI

```
Standalone only
```

### Local Development

```
Standalone (primary)
```

### Integration

```
Self-hosted Testnet-connected rippled
```

### Experimental Feature Validation

```
Devnet
```

---

## Strategic Conclusion (2026)

Devnet is:

* Feature-driven
* Amendment-driven
* Unstable by design
* Necessary for forward compatibility

It must be treated as:

> A feature capability sandbox — not a stability benchmark.

---

If desired, the next document can define:

* A feature-detection layer design in Go
* An amendment-aware transaction builder strategy
* A multi-network abstraction model for XRPL
* A unified AI Agent operational specification across Standalone/Testnet/Devnet
