## E2E Transaction Patterns

This section describes the E2E (End-to-End) transaction patterns available for Bitcoin Cash in the go-crypto-wallet system.

> **Common flow reference**: The 3-wallet setup, signing, and monitoring flows are defined in
> [docs/transaction-flow.md](../../transaction-flow.md). The patterns below describe
> BCH-specific address formats, signing algorithms, and protocol constraints.

### BCH Pattern Limitations vs BTC

Due to protocol differences, BCH supports **fewer patterns** than BTC:

| Feature | BTC Support | BCH Support | Impact |
|---------|-------------|-------------|--------|
| **SegWit** | ✅ Activated 2017 | ❌ Not implemented | No P2WPKH, P2WSH patterns |
| **Taproot** | ✅ Activated 2021 | ❌ Not implemented | No P2TR patterns |
| **Schnorr Signatures** | ✅ BIP340 | ❌ Not available | ECDSA only |
| **MuSig2** | ✅ BIP327 | ❌ Not available | No signature aggregation |
| **PSBT** | ✅ BIP174 | ⚠️ Limited | Raw transaction format |

### BCH E2E Pattern Matrix

| Pattern | Key Type | Signature Pattern | Address Format | E2E Script Status |
|---------|----------|-------------------|----------------|-------------------|
| **1** | **CashAddr P2PKH** | **Single-sig** | **`bitcoincash:q...`** | **🔶 Manual testing** |
| **2** | **CashAddr P2SH** | **2-of-3 Multisig** | **`bitcoincash:p...`** | **❌ Not implemented** |
| **3** | **CashAddr P2SH** | **3-of-3 Multisig** | **`bitcoincash:p...`** | **✅ e2e-workflow.sh** |

### Pattern 1: BCH CashAddr P2PKH Single-sig

```
Address Type: CashAddr P2PKH (BIP44)
Signing Requirements: Single-sig (Keygen only)
Address Format: bitcoincash:q... (P2PKH), bchtest:q... (Testnet)
```

**Workflow:**

Follows the [common single-sig flow](../../transaction-flow.md#single-sig-flow).
BCH-specific signing: **ECDSA** (no Schnorr).

**Characteristics:**

- Simple and fast (completed with single signature)
- Sign1/Sign2 wallets not required
- Uses BIP44 key derivation path (`m/44'/145'/account'/change/index` for mainnet)
- CashAddr P2PKH format (`bitcoincash:q...` prefix)
- Suitable for `client` account type

**Key Derivation:**

```
Mainnet: m/44'/145'/account'/change/index
Testnet/Regtest: m/44'/1'/account'/change/index
```

**Implementation Status:** 🔶 Manual testing (no dedicated E2E script)

### Pattern 2: BCH CashAddr P2SH 2-of-3 Multisig

```
Address Type: CashAddr P2SH (BIP44 + BIP11)
Signing Requirements: 2-of-3 Multisig (any 2 of Keygen, Sign1, Sign2)
Address Format: bitcoincash:p... (P2SH), bchtest:p... (Testnet)
```

**Workflow:**

Follows the [common multisig flow](../../transaction-flow.md#multisig-flow-m-of-n) with M=2, N=3.
Signing stops after 2 signatures (`isCompleted: true`); Sign2 is not required.
BCH-specific signing: **ECDSA** (no Schnorr).

**Characteristics:**

- Requires 2 of 3 possible signatures
- Provides redundancy - losing one key doesn't prevent access
- Uses BIP44 key derivation with BIP11 multisig
- CashAddr P2SH format (`bitcoincash:p...` prefix)
- Suitable for `deposit`, `payment` accounts

**Redeem Script:**

```
2 <PubKey1> <PubKey2> <PubKey3> 3 OP_CHECKMULTISIG
```

**Implementation Status:** ❌ Not implemented (planned)

### Pattern 3: BCH CashAddr P2SH 3-of-3 Multisig (Current E2E)

```
Address Type: CashAddr P2SH (BIP44 + BIP11)
Signing Requirements: 3-of-3 Multisig (Keygen + Sign1 + Sign2)
Address Format: bitcoincash:p... (P2SH), bchtest:p... (Testnet)
E2E Script: scripts/operation/bch/e2e-workflow.sh
```

**Workflow:**

Follows the [common multisig flow](../../transaction-flow.md#multisig-flow-m-of-n) with M=3, N=3.
All 3 signatures required (Keygen + Sign1 + Sign2).
BCH-specific signing: **ECDSA** (no Schnorr).

**Characteristics:**

- All 3 signatures required (maximum security)
- Uses BIP44 key derivation with BIP11 multisig
- CashAddr P2SH format (`bitcoincash:p...` prefix)
- Suitable for `stored` account type (cold storage)

**Redeem Script:**

```
3 <PubKey1> <PubKey2> <PubKey3> 3 OP_CHECKMULTISIG
```

**Implementation Status:** ✅ Implemented in `scripts/operation/bch/e2e-workflow.sh`

### BCH vs BTC Pattern Comparison

| BTC Pattern | BCH Equivalent | Notes |
|-------------|----------------|-------|
| Pattern 1 (P2PKH Single-sig) | **BCH Pattern 1** | Same concept, CashAddr format |
| Pattern 2 (P2PKH 2-of-3) | **BCH Pattern 2** | Same concept, CashAddr format |
| Pattern 3 (P2SH-P2WPKH) | ❌ N/A | BCH has no SegWit |
| Pattern 4 (P2SH-P2WSH 2-of-3) | ❌ N/A | BCH has no SegWit |
| Pattern 5 (P2WPKH) | ❌ N/A | BCH has no SegWit |
| Pattern 6 (P2WSH 2-of-3) | ❌ N/A | BCH has no SegWit |
| Pattern 7 (P2WSH 3-of-3) | **BCH Pattern 3** | Similar but no SegWit benefits |
| Pattern 8 (P2SH-P2WSH 3-of-3) | **BCH Pattern 3** | Similar but no SegWit benefits |
| Pattern 9 (P2TR Single-sig) | ❌ N/A | BCH has no Taproot |
| Pattern 10 (P2TR MuSig2) | ❌ N/A | BCH has no Taproot/Schnorr |
| Pattern 11 (P2TR Tapscript) | ❌ N/A | BCH has no Taproot |

### Account Types and Recommended Patterns

| Account | Purpose | Recommended BCH Pattern | Reason |
|---------|---------|-------------------------|--------|
| **client** | Customer deposit addresses | Pattern 1 (Single-sig) | Simple, customer-side operations |
| **deposit** | Deposit aggregation | Pattern 2 or 3 (Multisig) | Enhanced security |
| **payment** | Payments | Pattern 2 or 3 (Multisig) | Approval workflow |
| **stored** | Long-term storage | Pattern 3 (3-of-3) | Highest security |

### E2E Script Reference

| Script | Pattern | Description |
|--------|---------|-------------|
| `scripts/operation/bch/e2e-workflow.sh` | Pattern 3 | 3-of-3 Multisig workflow |

**Usage:**

```bash
# Run complete E2E workflow
./scripts/operation/bch/e2e-workflow.sh

# Run with reset (fresh state)
./scripts/operation/bch/e2e-workflow.sh --reset

# Run with verbose output
./scripts/operation/bch/e2e-workflow.sh --verbose

# Cleanup only
./scripts/operation/bch/e2e-workflow.sh --cleanup
```

### Future BCH Upgrades (2025-2026)

While not currently implemented in go-crypto-wallet, BCH has planned upgrades:

| Upgrade | CHIP | Description | go-crypto-wallet Impact |
|---------|------|-------------|-------------------------|
| **CashTokens** | CHIP-2022-02 | Native fungible/NFT tokens | Not in scope |
| **VM Limits** | CHIP-2021-05 | Extended script capabilities | Not in scope |
| **BigInt** | CHIP-2024-07 | High-precision arithmetic | Not in scope |
| **Loops** | CHIP-2021-05 | Script loops (OP_BEGIN/OP_UNTIL) | Not in scope |
| **Functions** | CHIP-2025-05 | Reusable script functions | Not in scope |
| **Pay-2-Script** | CHIP-2024-12 | Direct P2S outputs | Not in scope |

**Note:** go-crypto-wallet focuses on standard transaction patterns. Advanced contract features are out of scope.

### Implementation Priority

| Priority | Pattern | Reason |
|----------|---------|--------|
| ✅ Completed | Pattern 3 (3-of-3 Multisig) | Security-first approach |
| 🔶 Manual | Pattern 1 (Single-sig) | Basic functionality available |
| 🔜 Planned | Pattern 2 (2-of-3 Multisig) | Operational flexibility |

---
