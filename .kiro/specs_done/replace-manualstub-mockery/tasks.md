# Implementation Plan

- [x] 1. Extend mockery configuration and regenerate all mocks

- [x] 1.1 Add missing XRP sub-interfaces and file-storage interface to the mockery configuration
  - Append `TransactionSubmitter`, `AccountInfoProvider`, and `TransactionPreparer` to the existing XRP API package entry in `.mockery.yaml`
  - Append `DescriptorFileWriter` to the existing file-storage package entry in `.mockery.yaml`
  - Confirm each interface name exactly matches the exported identifier in its source package before saving
  - _Requirements: 1.2, 2.1, 5.1, 5.4_

- [x] 1.2 Add the ETH API interfaces package entry to the mockery configuration
  - Add a new package entry for the ETH ports package, with output directory set to the ETH infrastructure mocks directory
  - List the 15 leaf ETH interfaces: `ETHLifecycle`, `ETHKeyAccessor`, `ETHTransactionSigner`, `ETHRawKeyImporter`, `ETHNodeAPIClient`, `ERC20er`, `EtherTxMonitor`, `ChainConfigProvider`, `BalanceChecker`, `TxCreator`, `GasEstimator`, `TxSigner`, `TxSender`, `TxMonitor`, `AddressValidator`
  - Exclude `Ethereumer` (DI-layer-only monolithic interface), `ETHTransactionSender` (type alias, not a named interface), and all composed interfaces
  - _Requirements: 3.1, 3.4_

- [x] 1.3 Regenerate all mocks and verify the output
  - Run `make mockery` to generate all configured mocks in one pass
  - Confirm that newly created mock files appear in each target mocks directory and each carries the `DO NOT EDIT` header
  - Verify `make go-lint` and `make check-build` pass with no errors, confirming the generated code compiles correctly against all interface definitions
  - _Requirements: 1.1, 1.4, 2.2, 3.2, 3.3, 4.1, 4.2, 4.3, 5.4_

- [x] 2. (P) Replace BTC keygen manual stubs with generated mocks

- [x] 2.1 (P) Replace cold-repository stubs in the generate-descriptor test file
  - Remove the `stubAccountRepo`, `stubAuthRepo`, and `stubSeedRepo` struct definitions from the test file
  - Add the cold-repository mocks package import and replace each stub construction with the corresponding `NewMock*(t)` constructor
  - Set up EXPECT() expectations that match the exact method calls and arguments the production code makes on those repositories
  - _Requirements: 1.1, 1.3, 1.5_

- [x] 2.2 (P) Replace the account-key and file-writer stubs in the export-descriptor test file
  - Remove the `stubAccountKeyRepo` and `stubDescriptorFileWriter` struct definitions from the test file
  - Add imports for the cold-repository and file-storage mocks packages and use `NewMock*(t)` constructors in their place
  - Preserve the `stubDescriptorGenerator` struct intact — it implements a use case interface and is intentionally out of scope for mockery generation (see Non-Goals in design.md)
  - Set up EXPECT() expectations matching the arguments the production export-descriptor logic passes at runtime
  - _Requirements: 1.2, 1.3, 1.5_

- [x] 2.3 Verify BTC keygen tests pass after stub replacement
  - Run `make go-test` scoped to the BTC keygen use case package
  - Confirm every test case passes and no new failures are introduced by the mock replacements
  - _Requirements: 1.6_

- [x] 3. (P) Replace XRP watch manual mocks with generated mocks

- [x] 3.1 (P) Replace the transaction-submitter mock in the send-transaction test file
  - Remove the `MockTransactionSubmitter` struct and all its manually-written method bodies from the test file
  - Import the XRP infrastructure mocks package and update the test dependency struct to hold the generated mock type
  - Convert each `On(…).Return(…)` setup to an equivalent EXPECT() builder call on the generated mock
  - _Requirements: 2.3, 2.4_

- [x] 3.2 (P) Replace the account-info and transaction-preparer mocks in the create-transaction test file
  - Remove `MockAccountInfoProvider` and `MockTransactionPreparer` struct definitions and their method bodies
  - Import the XRP infrastructure mocks package and substitute generated constructors for the manual mock construction
  - Convert each `On(…).Return(…)` setup to an equivalent EXPECT() builder call on the respective generated mock
  - _Requirements: 2.3, 2.4_

- [x] 3.3 Verify XRP watch tests pass after mock replacement
  - Run `make go-test` scoped to the XRP watch use case package
  - Confirm every test case passes and no new failures are introduced by the mock replacements
  - _Requirements: 2.5_

- [x] 4. (P) Create the Claude mockery skill

- [x] 4.1 (P) Write the mockery skill file at `.claude/skills/mockery/SKILL.md`
  - Document when to activate the skill: whenever a developer asks Claude to generate a mock for a new interface or to add test coverage that requires a mock
  - Include the placement convention table that maps each ports package to its corresponding infrastructure mocks directory
  - Provide the step-by-step workflow: locate the interface source package → determine the correct target directory using the convention table → add or update the `.mockery.yaml` entry → run `make mockery` → verify with `make go-lint` and `make check-build`
  - Embed a copy-paste `.mockery.yaml` entry template for adding a new interface
  - Include a before/after example showing a manually-written stub replaced by a generated mock using the EXPECT() builder pattern
  - List all exceptions that must NOT be added to `.mockery.yaml`: `Ethereumer`, type aliases (e.g., `ETHTransactionSender`), use case interfaces, and compile-time conformance stubs in ports test files
  - Explicitly state that Claude must never write mock struct code by hand — always delegate to `make mockery`
  - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5, 6.6_

- [x] 5. Finalize convention rule and run full verification

- [x] 5.1 Confirm the mock placement convention rule aligns with the final configuration
  - Review `.claude/rules/internal/mockery.md` against the completed `.mockery.yaml`
  - Ensure the placement convention table in the rule file covers all packages now listed in `.mockery.yaml`, including the new ETH entry
  - Confirm the exceptions list (DI-only monolithic interface, type aliases, use case interfaces) matches the finalized design decisions
  - Update the rule file if any entry is missing or inaccurate
  - _Requirements: 5.1, 5.2, 5.3_

- [x] 5.2 Run complete build and test verification across all changed areas
  - Run `make go-lint` and `make check-build` to confirm the entire codebase compiles with no errors or warnings
  - Run `make go-test` across all packages that were modified to confirm no regressions
  - Confirm that no manually-written mock structs for ports interfaces remain anywhere in the codebase (search for struct types implementing ports interfaces in test files outside `mocks/` directories)
  - _Requirements: 1.6, 2.5, 3.5, 4.3_
