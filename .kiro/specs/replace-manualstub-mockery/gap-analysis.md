# Gap Analysis: replace-manualstub-mockery

## 1. Current State Investigation

### Mock Infrastructure (Existing)

The project already has a working mockery v3 setup:

- **Tool**: `github.com/vektra/mockery/v3 v3.6.1` (in `go.mod`)
- **Config**: `.mockery.yaml` — global settings with per-package entries
- **Command**: `make mockery`
- **Global conventions**: `filename: "mock_{{.InterfaceName | snakecase}}.go"`, `structname: "Mock{{.InterfaceName}}"`, template: `testify`, formatter: `goimports`

**Currently configured packages and their mock directories:**

| Package (ports) | Mock directory | Generated? |
|---|---|---|
| `ports/api/btc` (`Bitcoiner`) | `infrastructure/api/btc/mocks/` | ✅ Yes (1 file) |
| `ports/api/xrp` (4 interfaces) | `infrastructure/api/xrp/mocks/` | ❌ Not on disk |
| `ports/file` (2 interfaces) | `infrastructure/storage/file/transaction/mocks/` | ✅ Yes (2 files) |
| `ports/repository/watch` (8 interfaces) | `infrastructure/repository/watch/mocks/` | ✅ Yes (8 files) |
| `ports/repository/cold` (6 interfaces) | `infrastructure/repository/cold/mocks/` | ✅ Yes (6 files) |

### Manual Test Doubles Inventory

#### A. BTC Keygen — Simple Struct Stubs (in `*_test.go` files)

| Stub type | File | Implements | .mockery.yaml? | Generated? |
|---|---|---|---|---|
| `stubAccountRepo` | `generate_descriptor_test.go` | `BTCAccountKeyRepositorier` | ✅ Yes | ✅ Yes (cold/mocks) |
| `stubAuthRepo` | `generate_descriptor_test.go` | `AuthFullPubkeyRepositorier` | ✅ Yes | ✅ Yes (cold/mocks) |
| `stubSeedRepo` | `generate_descriptor_test.go` | `SeedRepositorier` | ✅ Yes | ✅ Yes (cold/mocks) |
| `stubAccountKeyRepo` | `export_descriptor_test.go` | `BTCAccountKeyRepositorier` | ✅ Yes | ✅ Yes (cold/mocks) |
| `stubDescriptorGenerator` | `export_descriptor_test.go` | `keygenusecase.GenerateDescriptorUseCase` | ❌ No | ❌ No |
| `stubDescriptorFileWriter` | `export_descriptor_test.go` | `file.DescriptorFileWriter` | ❌ No | ❌ No |

#### B. XRP Watch — testify/mock Struct Mocks (in `*_test.go` files)

| Mock type | File | Interface | Package | .mockery.yaml? |
|---|---|---|---|---|
| `MockTransactionSubmitter` | `send_transaction_test.go` | `apixrp.TransactionSubmitter` | `ports/api/xrp/transaction_submitter.go` | ❌ No |
| `MockAccountInfoProvider` | `create_transaction_test.go` | `apixrp.AccountInfoProvider` | `ports/api/xrp/account_info.go` | ❌ No |
| `MockTransactionPreparer` | `create_transaction_test.go` | `apixrp.TransactionPreparer` | `ports/api/xrp/interface.go` | ❌ No |

#### C. Interface Compliance Stubs (NOT targets — keep as-is)

`segregated_interfaces_test.go` contains `mockAccountInfoProvider` and `mockTransactionSubmitter` as compile-time conformance checks only (no testify dependency). These are intentionally minimal and should NOT be replaced.

### ETH Interface Gap

`ports/api/eth/interface.go` defines 15+ interfaces but **zero** are in `.mockery.yaml`:

- Core segregated: `ETHLifecycle`, `ETHKeyAccessor`, `ETHTransactionSigner`, `ETHRawKeyImporter`, `ETHNodeAPIClient`
- EIP-1559 flow: `TxCreator`, `GasEstimator`, `TxSigner`, `TxSender`, `TxMonitor`, `AddressValidator`, `ChainConfigProvider`, `BalanceChecker`
- Multi-use: `ERC20er`, `EtherTxMonitor`
- Composed (DI-level): `ETHKeygenSignClient`, `ETHWatchClient`, `WatchTxCreationDeps`, `KeygenSignTxDeps`

### Placement Convention (Existing Pattern)

```
ports/api/{chain}/      → infrastructure/api/{chain}/mocks/
ports/file/             → infrastructure/storage/file/transaction/mocks/
ports/repository/cold/  → infrastructure/repository/cold/mocks/
ports/repository/watch/ → infrastructure/repository/watch/mocks/
```

**Gap**: No convention defined for:
- Use case interface mocks (`keygenusecase.GenerateDescriptorUseCase`)
- `ports/file/DescriptorFileWriter` (not yet in `.mockery.yaml`)

---

## 2. Requirements Feasibility Analysis

### Requirement 1: Replace BTC Keygen Manual Stubs

**Technical needs**:
- `stubAccountRepo`, `stubAuthRepo`, `stubSeedRepo`, `stubAccountKeyRepo` → already-generated mocks exist in `cold/mocks/`; only test code changes needed
- `stubDescriptorFileWriter` → `file.DescriptorFileWriter` exists in `ports/file/interface.go` (line 58) but is missing from `.mockery.yaml`; needs an entry under `infrastructure/storage/file/descriptor/mocks/` or merged into existing `infrastructure/storage/file/transaction/mocks/`
- `stubDescriptorGenerator` → `keygenusecase.GenerateDescriptorUseCase` is defined in `application/usecase/keygen/interfaces.go`. This is a **use case interface**, not a ports interface. Placement requires a decision (see Options below)

**Constraint**: Use case interface mocks have no established placement precedent in this codebase.

### Requirement 2: Replace XRP Watch Manual Mocks

**Technical needs**:
- All 3 interfaces (`TransactionSubmitter`, `AccountInfoProvider`, `TransactionPreparer`) are already **correctly defined** in `ports/api/xrp/`
- Only `.mockery.yaml` entries are missing; target `dir: internal/infrastructure/api/xrp/mocks/`
- After generation, test files need to import from `xrpmocks` and use `.EXPECT()` pattern

**No constraint issues** — straightforward extension of existing pattern.

### Requirement 3: Add ETH API Interface Mocks

**Technical needs**:
- Add `ports/api/eth` package to `.mockery.yaml` with `dir: internal/infrastructure/api/eth/mocks/`
- Decision needed: which interfaces to include (see Options)
  - Composed interfaces (`WatchTxCreationDeps`, `KeygenSignTxDeps`) embed multiple small ones; mocking them would produce large mocks
  - `Ethereumer` must be excluded (DI-layer only, per comment in interface.go)
- `ETHTransactionSender` is a type alias (`= TxSender`), not a named interface — mockery cannot generate a mock for it directly; only `TxSender` needs a mock

### Requirement 4: Generate Missing XRP API Mocks

**Technical needs**: Run `make mockery`. The 4 interfaces are already in `.mockery.yaml`; they just haven't been generated yet.

### Requirement 5: Enforce Placement Convention

**Technical needs**: Document the convention and handle the two gap cases:
- `DescriptorFileWriter`: `ports/file/` → should follow existing `ports/file/` pattern → `infrastructure/storage/file/transaction/mocks/` (merge with existing entry) OR a new `infrastructure/storage/file/descriptor/mocks/`
- Use case interface mocks: new convention needed

### Requirement 6: Create Claude Skill

**Technical needs**: New file in `.claude/skills/mockery/SKILL.md`. No code changes required.

---

## 3. Implementation Approach Options

### Issue A: Placement of `DescriptorFileWriter` mock

**Option A1 — Merge into existing `ports/file` entry**
Add `DescriptorFileWriter` to the existing `ports/file` section in `.mockery.yaml`, keeping `dir: internal/infrastructure/storage/file/transaction/mocks/`. The `transaction/mocks/` name becomes slightly misleading since it now also holds the descriptor writer mock.

Trade-offs: ✅ Minimal YAML change ❌ Directory name misleads

**Option A2 — New entry with dedicated descriptor mocks dir**
Add a new package entry for `ports/file` with `dir: internal/infrastructure/storage/file/descriptor/mocks/`.

Trade-offs: ✅ Semantically correct directory naming ❌ One more directory, `.mockery.yaml` has two entries for the same package (not allowed — each package can only appear once)

**Option A3 — Rename and consolidate file storage mocks dir**
Rename `infrastructure/storage/file/transaction/mocks/` to `infrastructure/storage/file/mocks/` (covering all file storage interfaces) and update the single `.mockery.yaml` entry.

Trade-offs: ✅ Single cohesive directory for all file-storage mocks ❌ Requires renaming existing directory (import path change in test files)

**Recommendation**: A1 (least disruption); or A3 if a cleaner long-term structure is preferred.

---

### Issue B: Placement of use case interface mocks (`GenerateDescriptorUseCase`)

**Option B1 — Co-locate with use case package**
Add `.mockery.yaml` entry with `dir: internal/application/usecase/keygen/mocks/`. Breaks the "ports → infrastructure" convention but is pragmatic.

Trade-offs: ✅ Mock is near its interface and consumer ❌ Puts generated code in application layer ❌ Inconsistent with existing convention

**Option B2 — Keep stub as-is (no mockery for use case interfaces)**
Treat `stubDescriptorGenerator` as legitimate; it tests how `ExportDescriptorUseCase` orchestrates `GenerateDescriptorUseCase`. Since both are use cases (not ports), the interface stub is lightweight and appropriate.

Trade-offs: ✅ No convention violation ✅ Less configuration ❌ Inconsistency (one use case interface has a stub, others use mocks)

**Option B3 — Refactor to eliminate the dependency**
`ExportDescriptorUseCase` depends on `GenerateDescriptorUseCase` as a use case (unusual pattern). Consider refactoring so it depends on a ports-level interface instead, then mock that.

Trade-offs: ✅ Cleanest architecture ❌ Requires use case refactoring beyond the scope of this spec

**Recommendation**: B2 for now (keep stub for use case interface); document as exception in the skill. B3 is a separate refactoring spec.

---

### Issue C: Which ETH interfaces to mock

**Option C1 — Mock only leaf (non-composed) interfaces**
Include: `ETHLifecycle`, `ETHKeyAccessor`, `ETHTransactionSigner`, `ETHRawKeyImporter`, `ETHNodeAPIClient`, `TxCreator`, `GasEstimator`, `TxSigner`, `TxSender`, `TxMonitor`, `AddressValidator`, `ChainConfigProvider`, `BalanceChecker`, `ERC20er`, `EtherTxMonitor`.
Exclude: `Ethereumer` (DI-only), `ETHTransactionSender` (type alias), composed interfaces (`ETHKeygenSignClient`, `ETHWatchClient`, `WatchTxCreationDeps`, `KeygenSignTxDeps`).

Trade-offs: ✅ Minimal, ISP-aligned mocks ✅ Consistent with small-interface principle ❌ Tests using composed interfaces must compose mock calls manually

**Option C2 — Mock leaf + composed interfaces**
Include all from C1 plus composed interfaces.

Trade-offs: ✅ Tests depending on composed interfaces get ready-made mocks ❌ More generated code ❌ Composed mocks duplicate leaf method sets

**Recommendation**: C1. Composed interfaces are thin wrappers; Go's implicit satisfaction means a struct implementing all leaf interfaces automatically satisfies any composition.

---

## 4. Requirement-to-Asset Map

| Requirement | Existing Assets | Gaps | Complexity |
|---|---|---|---|
| Req 1: BTC stubs → mocks | 4 of 6 interfaces already have generated mocks | `DescriptorFileWriter` not in .mockery.yaml; `GenerateDescriptorUseCase` needs policy decision | S |
| Req 2: XRP watch mocks | All 3 interfaces defined in ports | 3 interfaces missing from .mockery.yaml | S |
| Req 3: ETH interface mocks | All ETH interfaces defined in ports | Entire ETH package missing from .mockery.yaml; `infrastructure/api/eth/mocks/` dir does not exist | S |
| Req 4: Generate XRP API mocks | 4 interfaces already in .mockery.yaml | Just not generated; run `make mockery` | XS |
| Req 5: Placement convention | Convention implicit in existing structure | Not documented; 2 gap cases need policy | S |
| Req 6: Claude mockery skill | `go-development` skill exists as model | New file only | S |

**Overall effort**: M (3–7 days)
**Risk**: Low — extends well-established mockery patterns; no architectural changes

---

## 5. Implementation Complexity & Risk

**Effort: M**
- Multiple files to update (`.mockery.yaml`, test files, new skill file)
- Regenerating all mocks with `make mockery` is fast but requires verification
- Test file changes require careful `.EXPECT()` call construction per interface

**Risk: Low**
- mockery v3 is already installed and operational
- All target interfaces are already cleanly defined in the correct layer
- No domain or infrastructure behavior changes required
- Build/lint verification is automated

---

## 6. Recommendations for Design Phase

### Preferred approach: Option C (Hybrid)

1. **Extend `.mockery.yaml`**: Add XRP segregated interfaces, `DescriptorFileWriter`, and ETH interfaces (leaf only) using A1 for file storage and C1 for ETH.
2. **Keep `stubDescriptorGenerator` as-is** (Option B2): Use case interface stubs are a legitimate pattern; document the exception in the skill.
3. **New `infrastructure/api/eth/mocks/` directory**: Created implicitly by `make mockery` once .mockery.yaml is updated.
4. **Test file updates**: Replace 4 BTC stubs and 3 XRP manual mocks with generated mock constructors and `.EXPECT()` patterns.
5. **Run `make mockery`** once for all changes.
6. **New skill file**: `.claude/skills/mockery/SKILL.md`.

### Research items to carry forward

- **None critical.** All interfaces exist; mockery configuration is well-understood.
- Minor: Confirm whether `.mockery.yaml` permits multiple package entries for the same Go import path (for the DescriptorFileWriter issue). If not, Option A1 (add to existing entry) is the only valid path.
