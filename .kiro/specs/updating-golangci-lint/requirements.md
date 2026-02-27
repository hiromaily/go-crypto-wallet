# Requirements Document

## Project Description (Input)

Enhance <https://github.com/golangci/golangci-lint> for Go's lint environment.

- Update to the latest version, v2.10.1.

- For extension rules, the plugin system should probably be available. <https://golangci-lint.run/docs/plugins/go-plugins/>
If current rule @tools/lint/rules.go is built beforehand using `go build -buildmode=plugin plugin/example.go`, it should run faster.

- Introducing the new modernize feature: <https://golangci-lint.run/docs/linters/configuration/#modernize>

- This is the most important feature: lint strictness can be adjusted depending on usage (meaning there are separate config files for each).
  - Local development pattern A (runs in CI even when creating PRs): A version using fast mode: <https://golangci-lint.run/docs/product/migration-guide/#lintersfast>. CI configuration for local must be run within 2 minutes.
  - Local development pattern B (runs in CI even when creating PRs): A version that does not use fast mode, but is configured only with other lightweight linters.
  - A version that strictly implements many linters and runs only once a day in CI (requires a GitHub action to run once a day).

- I wanna compare pattern A and pattern B.

## Introduction

This specification defines requirements for upgrading and restructuring the golangci-lint configuration in the go-crypto-wallet project. The upgrade targets golangci-lint v2.10.1 and introduces three distinct lint profiles (fast, lightweight, and strict) to balance developer velocity with code quality across different execution contexts (local development, CI on PRs, and scheduled nightly CI). The existing `tools/lint/rules.go` ruleguard plugin will be integrated using the native golangci-lint plugin system for improved performance.

## Requirements

### Requirement 1: Version Upgrade to v2.10.1

**Objective:** As a developer, I want golangci-lint upgraded to v2.10.1, so that the project benefits from the latest linter improvements, bug fixes, and new linter support.

#### Acceptance Criteria

1. The Lint Tool shall use golangci-lint version v2.10.1 in all CI workflows.
2. When golangci-lint runs locally, the Lint Tool shall use v2.10.1 as specified in the Makefile target or tooling configuration.
3. The Lint Tool shall maintain full backward compatibility with existing linter rules already defined in `.golangci.yml`.
4. If a linter used in the current configuration is removed or renamed in v2.10.1, the Lint Tool shall update the configuration to use the equivalent replacement.

---

### Requirement 2: Plugin System Integration for Custom Rules

**Objective:** As a developer, I want the existing `tools/lint/rules.go` ruleguard rules compiled as a native golangci-lint plugin, so that custom crypto-naming enforcement runs faster than the interpreted ruleguard mode.

#### Acceptance Criteria

1. The Lint Tool shall support loading `tools/lint/rules.go` as a compiled Go plugin using `go build -buildmode=plugin`.
2. When the plugin binary is present, the Lint Tool shall load custom rules from the compiled plugin rather than interpreting the source file directly.
3. The build system shall provide a `make` target (e.g., `make lint-plugin`) that compiles `tools/lint/rules.go` into a plugin binary.
4. Where the plugin system is enabled, the Lint Tool shall apply the same `forbidCryptoNaming` rule logic as the current ruleguard source-based configuration.
5. If the plugin binary is absent or stale, the Lint Tool shall fall back gracefully to the source-based ruleguard approach.
6. The plugin build process shall be documented in the project's development guidelines.

---

### Requirement 3: Modernize Linter Introduction

**Objective:** As a developer, I want the `modernize` linter enabled in the appropriate lint profiles, so that the codebase is automatically flagged for idiomatic Go modernization opportunities.

#### Acceptance Criteria

1. The Lint Tool shall enable the `modernize` linter in the strict lint profile.
2. When `modernize` detects outdated Go patterns, the Lint Tool shall report each issue with file path and line number.
3. The Lint Tool shall configure `modernize` settings consistent with the project's minimum supported Go version as declared in `go.mod`.
4. Where `modernize` produces false positives for project-specific patterns, the Lint Tool shall define exclusion rules in the configuration.
5. The `modernize` linter shall be evaluated for inclusion in the lightweight profile based on its execution time impact.

---

### Requirement 4: Fast Mode Lint Profile (Pattern A)

**Objective:** As a developer, I want a fast-mode lint configuration for use in PR CI checks, so that lint feedback is available within 2 minutes to maintain developer velocity.

#### Acceptance Criteria

1. The Lint Tool shall provide a dedicated configuration file (e.g., `.golangci-fast.yml`) that enables `linters.fast: true`.
2. While running in fast mode, the Lint Tool shall complete execution within 2 minutes on the CI runner (ubuntu-slim) for typical PR changesets.
3. When the fast-mode profile is used in CI, the Lint Tool shall apply it only to changed files (using `only-new-issues: true`).
4. The fast-mode configuration shall include the subset of linters available under `linters.fast` mode as defined by golangci-lint v2.10.1.
5. The CI workflow shall invoke Pattern A as the primary go-lint job on pull requests targeting the `main` branch.
6. If Pattern A execution exceeds 2 minutes, the CI job shall be considered failed and require investigation.

---

### Requirement 5: Lightweight Lint Profile (Pattern B)

**Objective:** As a developer, I want a lightweight lint configuration without fast mode but limited to low-overhead linters for PR CI checks, so that a curated set of meaningful linters runs efficiently as an alternative to Pattern A.

#### Acceptance Criteria

1. The Lint Tool shall provide a dedicated configuration file (e.g., `.golangci-lightweight.yml`) that does not use `linters.fast: true`.
2. The lightweight profile shall include only explicitly selected linters known to have low execution overhead.
3. When the lightweight profile is used in CI, the Lint Tool shall apply it only to changed files (using `only-new-issues: true`).
4. The lightweight profile shall include at minimum: `govet`, `staticcheck`, `errcheck`, `unused`, `ineffassign`, and `misspell`.
5. The CI workflow shall provide a separate job or workflow for running Pattern B in parallel with Pattern A on pull requests, enabling direct comparison.
6. The Lint Tool shall report execution time for Pattern B to facilitate comparison with Pattern A.

---

### Requirement 6: Strict Lint Profile (Nightly CI)

**Objective:** As a developer, I want a strict full-linter configuration that runs once daily in CI, so that comprehensive code quality checks are enforced without blocking PR workflows.

#### Acceptance Criteria

1. The Lint Tool shall provide a dedicated configuration file (e.g., `.golangci-strict.yml`) that mirrors the current `.golangci.yml` linter set, extended with `modernize` and any additional linters appropriate for thorough analysis.
2. The strict profile shall run all linters against the entire codebase (not only changed files).
3. A GitHub Actions scheduled workflow shall trigger the strict lint profile once per day (e.g., at 02:00 UTC).
4. When the nightly strict lint job detects issues, the GitHub Actions workflow shall mark the run as failed and surface the report in the repository's Actions tab.
5. The nightly workflow shall run on the `main` branch only.
6. The nightly workflow shall use `ubuntu-slim` as the runner and set a `timeout-minutes` limit appropriate for full-codebase analysis (e.g., 15 minutes).
7. Where the strict profile introduces new linter warnings not present in the current configuration, the Lint Tool shall initially report them as warnings to allow a grace period for remediation.

---

### Requirement 7: CI Workflow Configuration

**Objective:** As a maintainer, I want CI workflows updated to use the new lint profiles correctly, so that the right profile runs in the right context without duplicating effort.

#### Acceptance Criteria

1. The existing `lint-test.yml` workflow shall be updated to run Pattern A as the default go-lint job for PRs.
2. The CI workflow shall add a parallel job to run Pattern B on PRs for comparison purposes.
3. A new `nightly-lint.yml` GitHub Actions workflow file shall be created for the strict profile scheduled run.
4. When golangci-lint-action is updated, the CI workflow shall pin it to a version compatible with golangci-lint v2.10.1.
5. The Lint Tool shall update the `version` field in `golangci/golangci-lint-action` from `v2.8.0` to `v2.10.1` across all workflow files.
6. All new and updated workflow jobs shall use `ubuntu-slim` as the runner.
7. All workflow jobs shall define a `timeout-minutes` value.

---

### Requirement 8: Pattern A vs Pattern B Comparison Tooling

**Objective:** As a developer, I want a mechanism to compare the lint coverage and execution time of Pattern A (fast mode) and Pattern B (lightweight), so that an informed decision can be made about which profile to standardize on for PR CI.

#### Acceptance Criteria

1. The CI workflow shall run both Pattern A and Pattern B as separate parallel jobs on the same PR to produce side-by-side results in the GitHub Actions summary.
2. Each profile job shall report the number of issues found and the total execution time.
3. The Lint Tool's documentation shall include a comparison table of linters included in Pattern A vs Pattern B.
4. After a defined evaluation period, the project shall select one profile as the primary PR lint check and document the rationale in the project guidelines.
