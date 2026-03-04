# Verify E2E flow $ARGUMENTS

Verify that the E2E scripts for the specified chain follow the documented transaction flow.

## Usage

```
/e2e:verify-e2e-flow btc
/e2e:verify-e2e-flow eth
/e2e:verify-e2e-flow xrp
/e2e:verify-e2e-flow bch
/e2e:verify-e2e-flow all
```

## Reference Documents

**MUST read first:**

- `docs/transaction-flow.md` — canonical 3-wallet flow (Watch / Keygen / Sign) and transaction lifecycle

**Chain-specific (load based on argument):**

| Chain | Scripts Directory | Chain Doc |
|-------|-------------------|-----------|
| `btc` | `scripts/operation/btc/e2e/` | `docs/chains/btc/README.md` |
| `eth` | `scripts/operation/eth/e2e/` | `docs/chains/eth/README.md` |
| `xrp` | `scripts/operation/xrp/e2e/` | `docs/chains/xrp/README.md` |
| `bch` | `scripts/operation/bch/e2e/` | `docs/chains/bch/README.md` |
| `all` | All of the above | All of the above |

## Verification Steps

### Step 1: List all e2e scripts

Use Glob to list all scripts in `scripts/operation/{chain}/e2e/*.sh` (excluding parallel runner).

### Step 2: Read each script

Read all pattern scripts (not `e2e-parallel-runner.sh`).

### Step 3: Check each script against the documented flow

For each script, verify the following phases are present and correct:

#### Setup Flow

**Single-sig patterns:**

- [ ] Keygen: `create seed`
- [ ] Keygen: `create hdkey` (for all required accounts)
- [ ] Keygen: export address / descriptor
- [ ] Watch: import address / descriptor

**Multisig patterns (additional steps):**

- [ ] Sign wallet(s): `create seed`
- [ ] Sign wallet(s): `create hdkey`
- [ ] Sign wallet(s): `export fullpubkey`
- [ ] Keygen: `import fullpubkey` (for each sign wallet)

#### Transaction Flow

- [ ] Watch: `create [payment|transfer|deposit]` → produces unsigned tx file
- [ ] Keygen: `sign signature --file <unsigned_tx>` → produces signed/partial tx file
- [ ] Sign wallet(s): `sign signature --file <partial_tx>` (multisig only, repeated until threshold met)
- [ ] Watch: `send tx --file <signed_tx>` → broadcasts and returns txID

#### Monitoring Flow

- [ ] Watch: `monitor senttx` (optional but preferred)

### Step 4: Identify deviations and special cases

Flag the following as findings (not necessarily errors):

| Category | Description |
|----------|-------------|
| **Test-only phases** | UTXO generation, payment request creation, contract deployment, balance verification — not part of production flow |
| **Placeholder** | Signing steps marked as `# placeholder` or `TODO` (e.g., MuSig2, Tapscript) |
| **Missing monitoring** | `monitor senttx` not called |
| **Address export mechanism** | ETH reads keygen DB directly instead of using CLI export command |
| **Transaction type coverage** | Only `payment` or `transfer` tested (deposit/transfer may not be exercised) |

## Output Format

Report results using this structure:

### Compliance Matrix

| Pattern | Type | Wallets | Setup | Create TX | Sign (Keygen) | Sign (Sign wallets) | Send | Monitor | Status |
|---------|------|---------|-------|-----------|----------------|---------------------|------|---------|--------|
| p1 ... | single/multi | watch+keygen[+sign1+sign2] | ✅/❌ | ✅/❌ | ✅/❌ | ✅/⚠️/— | ✅/❌ | ✅/❌ | ✅/⚠️/❌ |

**Status key**: ✅ fully compliant · ⚠️ partial/placeholder · ❌ missing or incorrect

### Findings

List each deviation with:

- **Pattern**: which script(s) are affected
- **Category**: from the deviation table above
- **Detail**: specific description

### Summary

One paragraph conclusion: overall compliance level, any critical gaps, and patterns that are incomplete vs. production-ready.
