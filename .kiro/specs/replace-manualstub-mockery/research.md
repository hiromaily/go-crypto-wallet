# Research & Design Decisions

---

## Summary

- **Feature**: `replace-manualstub-mockery`
- **Discovery Scope**: Extension (existing mockery infrastructure, adding missing entries)
- **Key Findings**:
  - mockery v3.6.1 is already installed; `.mockery.yaml` and `make mockery` are operational.
  - All target interfaces (XRP sub-interfaces, ETH segregated interfaces, `DescriptorFileWriter`) are correctly defined in `application/ports/` — only `.mockery.yaml` entries and test file updates are required.
  - `stubDescriptorGenerator` implements a use case interface (`keygenusecase.GenerateDescriptorUseCase`), not a ports interface; it is explicitly out of scope for mockery generation per the placement convention.

---

## Research Log

### Topic: Existing Mockery Configuration

- **Context**: Determine what is already configured vs. what is missing.
- **Sources Consulted**: `.mockery.yaml`, `go.mod`, existing `mocks/` directories.
- **Findings**:
  - 5 packages configured; 17 mock files generated.
  - XRP API (`Rippler`, `RippleAPIer`, `RipplePublicer`, `RippleAdminer`) configured but NOT generated on disk.
  - ETH package (`ports/api/eth`) absent entirely from `.mockery.yaml`.
  - `DescriptorFileWriter` (`ports/file`) absent from `.mockery.yaml` despite existing entry for the same package.
  - XRP sub-interfaces (`TransactionSubmitter`, `AccountInfoProvider`, `TransactionPreparer`) absent from `.mockery.yaml`.
- **Implications**: Minimal config change — add 3 packages/entries to `.mockery.yaml`, then run `make mockery`.

### Topic: BTC Keygen Stub Analysis

- **Context**: Determine which BTC stubs can be replaced with already-generated mocks vs. which require new configuration.
- **Findings**:
  - `stubAccountRepo`, `stubAuthRepo`, `stubSeedRepo`, `stubAccountKeyRepo` → implement `BTCAccountKeyRepositorier`, `AuthFullPubkeyRepositorier`, `SeedRepositorier` (all already generated in `cold/mocks/`). Only test import updates needed.
  - `stubDescriptorFileWriter` → implements `file.DescriptorFileWriter` (defined at `ports/file/interface.go:58`). This interface is in the same package already listed in `.mockery.yaml` but under a different entry key. Action: add `DescriptorFileWriter` to the existing `ports/file` `.mockery.yaml` entry.
  - `stubDescriptorGenerator` → implements `keygenusecase.GenerateDescriptorUseCase` (defined in `application/usecase/keygen/interfaces.go`). This is a **use case interface**, not a ports interface. The project convention restricts mockery to ports interfaces only. Decision: keep as struct stub.
- **Implications**: Test file changes in `generate_descriptor_test.go` and `export_descriptor_test.go`; one new `.mockery.yaml` entry for `DescriptorFileWriter`.

### Topic: XRP Watch Manual Mock Analysis

- **Context**: Determine whether XRP manual mocks correspond to defined ports interfaces.
- **Findings**:
  - `TransactionSubmitter` → defined in `ports/api/xrp/transaction_submitter.go:33`, 3 methods.
  - `AccountInfoProvider` → defined in `ports/api/xrp/account_info.go:31`, 2 methods.
  - `TransactionPreparer` → defined in `ports/api/xrp/interface.go:262`, 1 method.
  - All three are legitimate ISP-compliant segregated interfaces, not duplicates of the monolithic `XRPer`.
  - Target mock directory: `internal/infrastructure/api/xrp/mocks/` (same as Rippler et al.).
- **Implications**: Add `TransactionSubmitter`, `AccountInfoProvider`, `TransactionPreparer` to the existing `ports/api/xrp` entry in `.mockery.yaml`.

### Topic: ETH Interface Inventory

- **Context**: Enumerate which ETH interfaces qualify for mockery generation.
- **Findings**:
  - Leaf interfaces (15 candidates): `ETHLifecycle`, `ETHKeyAccessor`, `ETHTransactionSigner`, `ETHRawKeyImporter`, `ETHNodeAPIClient`, `ERC20er`, `EtherTxMonitor`, `ChainConfigProvider`, `BalanceChecker`, `TxCreator`, `GasEstimator`, `TxSigner`, `TxSender`, `TxMonitor`, `AddressValidator`.
  - Composed interfaces (`ETHKeygenSignClient`, `ETHWatchClient`, `WatchTxCreationDeps`, `KeygenSignTxDeps`) embed leaf interfaces; no separate mocks needed.
  - `Ethereumer` — DI-only monolithic interface; must be excluded.
  - `ETHTransactionSender` — type alias for `TxSender` (`ETHTransactionSender = TxSender`); mockery cannot mock type aliases; mock `TxSender` instead.
  - Target mock directory: `internal/infrastructure/api/eth/mocks/` (does not exist yet; created by `make mockery`).
- **Implications**: New `.mockery.yaml` package entry listing 15 leaf interfaces.

### Topic: DescriptorFileWriter Placement

- **Context**: `.mockery.yaml` has only one allowed entry per Go package import path.
- **Findings**:
  - `ports/file` already has an entry (`dir: infrastructure/storage/file/transaction/mocks/`, interfaces: `TransactionFileRepositorier`, `AddressFileRepositorier`).
  - `DescriptorFileWriter` is in the same package; must be added to the same entry (not a new entry).
  - The `transaction/mocks/` directory name is slightly imprecise for a descriptor-related mock, but adding a new entry is not possible in mockery v3 YAML (one entry per package).
  - Alternative directory rename (`file/mocks/`) would break existing imports in 2 test files; not worth the disruption.
- **Implications**: Append `DescriptorFileWriter:` to the existing `ports/file` interfaces list. Accept minor naming imprecision.

---

## Architecture Pattern Evaluation

| Option | Description | Strengths | Risks | Notes |
|--------|-------------|-----------|-------|-------|
| Extend existing .mockery.yaml | Add missing entries to existing config | Minimal change, consistent tooling | Directory naming imprecision for DescriptorFileWriter | Selected |
| Manual mocks for missing interfaces | Keep manual mocks, do not extend mockery | No config change | Defeats the purpose; maintenance burden grows | Rejected |
| Rename mock directories | Unify all file-storage mocks under `file/mocks/` | Cleaner naming | Breaks existing test imports | Deferred to future cleanup |

---

## Design Decisions

### Decision: Keep `stubDescriptorGenerator` as struct stub

- **Context**: `export_descriptor_test.go` uses `stubDescriptorGenerator` implementing `keygenusecase.GenerateDescriptorUseCase`, a use case interface (not a ports interface).
- **Alternatives Considered**:
  1. Add `application/usecase/keygen` package to `.mockery.yaml` → breaks placement convention; puts generated code in application layer.
  2. Keep struct stub → consistent with convention; lightweight.
- **Selected Approach**: Keep `stubDescriptorGenerator` as a simple struct stub.
- **Rationale**: The placement convention restricts mockery to `application/ports/` interfaces only. Use case interfaces injected into other use cases are a testing seam that struct stubs handle adequately.
- **Trade-offs**: Minor inconsistency within `export_descriptor_test.go` (one stub, one generated mock). Documented as exception in Claude rule and skill.
- **Follow-up**: If `export_descriptor_test.go` grows complex, consider promoting `GenerateDescriptorUseCase` to a ports interface.

### Decision: Exclude composed ETH interfaces from mockery

- **Context**: `WatchTxCreationDeps`, `KeygenSignTxDeps`, `ETHKeygenSignClient`, `ETHWatchClient` compose leaf interfaces. Generating mocks for them would duplicate method sets already available through leaf mocks.
- **Alternatives Considered**:
  1. Include composed interfaces → more mocks, but tests can use a single mock that satisfies all requirements.
  2. Exclude composed interfaces → tests compose leaf mocks; minimal generated code.
- **Selected Approach**: Exclude composed interfaces.
- **Rationale**: Go's implicit interface satisfaction means any struct implementing all leaf methods satisfies any composition. Generating composed mocks adds redundancy.
- **Trade-offs**: Tests must inject multiple mock arguments instead of one combined mock. Acceptable given the ISP design intent of the composed interfaces.

### Decision: Append `DescriptorFileWriter` to existing `ports/file` .mockery.yaml entry

- **Context**: `.mockery.yaml` permits exactly one entry per package import path. `DescriptorFileWriter` lives in `ports/file`, same as `TransactionFileRepositorier`.
- **Alternatives Considered**:
  1. Append to existing entry (same `dir: transaction/mocks/`) → minor naming imprecision.
  2. Rename directory to `file/mocks/` → breaks imports in existing test files.
- **Selected Approach**: Append to existing entry.
- **Rationale**: Least disruption; naming imprecision is acceptable since the directory already contains non-transaction mocks (`AddressFileRepositorier`).
- **Trade-offs**: `transaction/mocks/` hosts a descriptor mock — imprecise but harmless.

---

## Risks & Mitigations

- `make mockery` regenerates all configured mocks simultaneously — if an interface changes, all mocks are regenerated correctly. Risk: low.
- ETH `infrastructure/api/eth/mocks/` directory does not exist yet — mockery creates it automatically. Risk: none.
- BTC test refactoring may introduce `EXPECT()` call mismatches — mitigated by running `make go-test` after each test file update.
- `stubDescriptorGenerator` left as struct stub — documents an exception. Risk of future contributors adding more use-case stubs. Mitigated by Claude rule documenting the exception explicitly.

---

## References

- [mockery v3 configuration docs](https://vektra.github.io/mockery/latest/configuration/) — YAML schema and package entry format
- [testify/mock EXPECT() pattern](https://pkg.go.dev/github.com/stretchr/testify/mock) — builder pattern for mock expectations
- `.kiro/specs/replace-manualstub-mockery/gap-analysis.md` — full gap analysis with option evaluations
- `.claude/rules/internal/mockery.md` — project rule encoding the placement convention
