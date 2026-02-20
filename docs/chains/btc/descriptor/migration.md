# Migrating to Descriptor Wallets

This guide outlines how to move from legacy address management to descriptor-based wallets without disrupting existing operations.

## Should You Migrate?

Migrate if you need Bitcoin Core compatibility, explicit address derivation, or simpler backups. Stay on legacy flows if the current setup is stable and you do not need descriptor portability yet.

## Phased Migration Plan

### Phase 1: Parallel Operation

- Generate descriptors alongside existing legacy addresses.
- Import descriptors into the watch wallet with a limited range to validate balances.
- Keep using legacy addresses for production transactions while monitoring descriptor-derived addresses.

### Phase 2: Gradual Cutover

- Start issuing new receive/change addresses from descriptor exports.
- Keep legacy addresses active for existing UTXOs.
- Export descriptor files regularly and back them up offline.

### Phase 3: Full Migration

- Move remaining funds to descriptor-derived addresses.
- Update operational runbooks to reference descriptor exports and validation.
- Decommission legacy address generation paths only after reconciliation.

## Backward Compatibility

- Legacy wallet flows continue to work; descriptors are opt-in.
- Descriptor imports currently target single-key descriptors; multisig imports are rejected until support is added.
- Bitcoin Core imports can coexist with the watch wallet imports for monitoring.

## Rollback Plan

1. Stop using descriptor-derived addresses.
2. Continue operating with legacy address generation.
3. Investigate and resolve descriptor issues (syntax, network mismatch, derivation path).
4. Re-run validation and import after fixes.

## Operational Checklist

- Export both receive and change branches (`--include-change`) for parity with legacy behavior.
- Encrypt descriptor backups and store them offline.
- Validate descriptor files before import (`watch --coin btc descriptor validate`).
- Align derivation paths with the intended network (mainnet/testnet/regtest) when exporting/importing.
- Document the Bitcoin Core import command and RPC credentials if operators rely on Core for monitoring.
