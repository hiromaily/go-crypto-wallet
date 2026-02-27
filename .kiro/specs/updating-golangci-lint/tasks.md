# Implementation Plan

> **Deferred**: Requirement 2 (native Go plugin compilation for `tools/lint/rules.go`) is explicitly deferred per the design decision documented in `research.md`. The native Go plugin system requires CGO and exact transitive dependency matching, which is fragile for a single-rule file. The existing gocritic/ruleguard approach is preserved unchanged.

---

- [x] 1. Upgrade golangci-lint version dependency

- [x] 1.1 Update the go.mod tool dependency to golangci-lint v2.10.1
  - Run `go get github.com/golangci/golangci-lint/v2@v2.10.1` to update the version pin
  - Run `go mod tidy` to regenerate go.sum and update indirect dependencies
  - Confirm the `require` block in `go.mod` reflects `v2.10.1`
  - Review the go.sum diff to ensure no unexpected transitive dependency changes
  - _Requirements: 1.1, 1.2_

- [x] 1.2 Verify backward compatibility of existing `.golangci.yml` with v2.10.1
  - Run `go tool golangci-lint config verify` to validate the existing config parses correctly
  - Run `make go-lint` against the full codebase and confirm the issue set is unchanged
  - If any linter was renamed or removed between v2.8.0 and v2.10.1, update the config to use the replacement name
  - Verify that the new gosec rules (G117, G602, G701-G706) added in v2.10.0 do not produce false positives on the existing codebase; add targeted exclusions if needed
  - _Requirements: 1.3, 1.4_

---

- [x] 2. Create Pattern A fast-mode lint profile

- [x] 2.1 (P) Create `.golangci-fast.yml` for fast-mode PR linting
  - Add a `run` block with `timeout: 2m`, `go` version matching `go.mod`, `build-tags: [integration]`, and `allow-parallel-runners: true`
  - Set `linters.default: fast` to activate all fast-tagged linters in golangci-lint v2
  - Include the `issues` block with `max-issues-per-linter: 0` and `max-same-issues: 0`
  - Carry over the `importas` alias settings and `gocritic` enabled-checks (ruleguard) from the base config, as these apply to fast linters
  - Carry over relevant `exclusions.rules` entries from `.golangci.yml` that apply to fast linters
  - Run `go tool golangci-lint config verify -c .golangci-fast.yml` to confirm the file is valid
  - _Requirements: 4.1, 4.2, 4.4_

- [x] 2.2 (P) Add `go-lint-fast` Makefile target
  - Add a `go-lint-fast` target in `make/lint.mk` that invokes `go tool golangci-lint run -c .golangci-fast.yml --fix`
  - Add a `go-lint-fast-check` variant (no `--fix`) for read-only validation
  - Add both targets to the `.PHONY` declaration
  - Verify the target runs to completion locally within 2 minutes
  - _Requirements: 4.1, 4.2_

---

- [x] 3. Create Pattern B lightweight lint profile

- [x] 3.1 (P) Create `.golangci-lightweight.yml` for curated lightweight PR linting
  - Set `linters.default: none` and enumerate an explicit `linters.enable` list: `govet`, `staticcheck`, `errcheck`, `unused`, `ineffassign`, `misspell`, `copyloopvar`, `fatcontext`, `nolintlint`, `revive`
  - Add a `run` block with `timeout: 2m`, `go` version matching `go.mod`, and `build-tags: [integration]`
  - Copy relevant settings blocks from the base config: `govet.enable: [shadow]`, `revive.rules`, `staticcheck.checks`
  - Include the same `issues` and `exclusions` sections as the base config where applicable
  - Run `go tool golangci-lint config verify -c .golangci-lightweight.yml` to confirm the file is valid
  - _Requirements: 5.1, 5.2, 5.4_

- [x] 3.2 (P) Add `go-lint-lightweight` Makefile target
  - Add a `go-lint-lightweight` target in `make/lint.mk` that invokes `go tool golangci-lint run -c .golangci-lightweight.yml --fix`
  - Add a `go-lint-lightweight-check` variant (no `--fix`) for read-only validation
  - Add both targets to the `.PHONY` declaration
  - Record the observed execution time when running locally to use as a baseline for the A vs B comparison
  - _Requirements: 5.1, 5.6_

---

- [ ] 4. Create strict lint profile with modernize linter

- [x] 4.1 Create `.golangci-strict.yml` mirroring the base config plus `modernize`
  - Copy the full linter list, settings, and exclusion rules from `.golangci.yml`
  - Add `modernize` to the `linters.enable` list
  - Configure `modernize` settings aligned with the Go version declared in `go.mod` (e.g., `min-version`)
  - Set `run.timeout: 15m` for full-codebase analysis
  - Do not set `only-new-issues`; the strict profile analyses the entire codebase
  - Run `go tool golangci-lint config verify -c .golangci-strict.yml` to confirm the file is valid
  - _Requirements: 3.1, 3.3, 6.1, 6.2_

- [x] 4.2 Triage new warnings from the strict profile and add required exclusions
  - Run `go tool golangci-lint run -c .golangci-strict.yml ./...` locally on the full codebase
  - For each `modernize` finding, determine whether it is a true actionable improvement or a false positive for this project's patterns
  - Add exclusion entries in `.golangci-strict.yml` for any confirmed false positives, with comments referencing the rule ID and rationale
  - For each new gosec rule (G117, G602, G701-G706) that fires, apply the same triage process and add exclusions where appropriate
  - Confirm the profile either produces zero issues or documents known suppressions before enabling nightly enforcement
  - _Requirements: 3.2, 3.4, 3.5, 6.7_

- [x] 4.3 Add `go-lint-strict` Makefile target
  - Add a `go-lint-strict` target in `make/lint.mk` that invokes `go tool golangci-lint run -c .golangci-strict.yml`
  - Add the target to the `.PHONY` declaration
  - _Requirements: 6.1_

---

- [ ] 5. Update PR CI workflow to use fast-only mode

- [ ] 5.1 Update the `go-lint` job in `lint-test.yml` to use `--fast-only` and upgrade action to v9
  - Upgrade `golangci/golangci-lint-action` from `@v8` to `@v9` (node24 runtime; node20 is being deprecated by GitHub)
  - Add `args: --fast-only` to the action step
  - Keep the single `.golangci.yml` config (no separate config file needed)
  - Keep `only-new-issues: true` and `fetch-depth: 0` for the checkout step
  - Verify the golangci-lint version pin is `v2.10.1`
  - _Requirements: 4.3, 4.5, 7.1, 7.2, 7.6, 7.7, 8.1, 8.2_

- [ ] 5.2 Evaluate runner choice for the PR lint job
  - `ubuntu-slim` has fewer CPU cores, which may degrade golangci-lint parallel analysis performance
  - Run the PR lint job on both `ubuntu-slim` and `ubuntu-latest` (or a larger runner) and compare execution times
  - If `ubuntu-slim` causes timeout or significant slowdown, switch the `go-lint` job to `ubuntu-latest`
  - Document the chosen runner and rationale in the workflow comment
  - _Requirements: 7.6, 8.1_

---

- [ ] 6. Create nightly strict (full) lint workflow

- [ ] 6.1 Create `.github/workflows/nightly-lint.yml` for the daily full lint run
  - Set `on.schedule.cron: "0 2 * * *"` (02:00 UTC daily) and include `on.workflow_dispatch` for manual triggering
  - Add a single job `strict-lint` with `timeout-minutes: 15`
  - Use `ubuntu-latest` (or the runner determined by task 5.2); full lint is CPU-intensive and `ubuntu-slim` cores may be insufficient
  - Include `permissions: contents: read`
  - Use `actions/checkout@v6` with default ref (checks out `main` when triggered by schedule)
  - Use `actions/setup-go@v6` with `go-version` matching the repo's Go version and `cache: true`
  - Invoke `golangci/golangci-lint-action@v9` with `version: v2.10.1` (no `--fast-only`, no `only-new-issues`; runs all linters on the full codebase)
  - The workflow must be triggered via `workflow_dispatch` once manually to confirm it completes successfully before the first automated run
  - _Requirements: 6.3, 6.4, 6.5, 6.6, 7.3, 7.6, 7.7_

---

- [ ] 7. Validate both execution modes end-to-end

- [ ] 7.1 Verify the single config and both execution modes produce expected results
  - Run `go tool golangci-lint config verify` to confirm `.golangci.yml` is valid
  - Run `make go-lint-fast` and record completion time to confirm it stays under 2 minutes
  - Run `make go-lint` and record completion time for full (strict) mode baseline
  - _Requirements: 4.2, 4.4, 6.1, 7.4_

- [ ] 7.2 Open a test PR to verify CI jobs run correctly
  - Confirm `go-lint` job uses `--fast-only` and completes within the timeout
  - Manually trigger `nightly-lint.yml` via `workflow_dispatch` to validate the full lint workflow runs without errors before the first scheduled execution
  - _Requirements: 4.5, 4.6, 6.3, 6.4, 8.1, 8.2_
