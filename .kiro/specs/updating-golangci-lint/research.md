# Research & Design Decisions

---
**Purpose**: Capture discovery findings, architectural investigations, and rationale that inform the technical design.

---

## Summary

- **Feature**: `updating-golangci-lint`
- **Discovery Scope**: Extension (existing toolchain configuration)
- **Key Findings**:
  - golangci-lint is managed as a `go tool` dependency in `go.mod` (not via direct binary download), so the upgrade path is `go get github.com/golangci/golangci-lint/v2@v2.10.1` + `go mod tidy`.
  - The v1 `linters.fast: true` property was removed in v2; the equivalent is `linters.default: fast` (configuration) or the `--fast-only` CLI flag. These are **not strictly equivalent** to the old behavior — linter selection may differ.
  - The native Go plugin system (`-buildmode=plugin`) requires `CGO_ENABLED=1`, exact library version matching, and a specific `New(conf any) ([]*analysis.Analyzer, error)` export. The current `tools/lint/rules.go` is a ruleguard DSL file (loaded by gocritic's ruleguard check) and is **not** a Go plugin; a separate `plugin/main.go` wrapper would be required.
  - The `modernize` linter is not tagged "fast" in v2 documentation, making it suitable for the strict profile only.
  - v2.9.0 added Go 1.26 support; v2.10.0 added 8 new gosec rules (G117, G602, G701–G706).

## Research Log

### golangci-lint Versioning and Installation

- **Context**: Needed to understand how to upgrade from v2.8.0 to v2.10.1 in this repo.
- **Sources Consulted**: `go.mod` (line 172, 382), `make/lint.mk`, GitHub Releases
- **Findings**:
  - golangci-lint is registered as `go tool` in `go.mod` at the tool directive (line 382): `github.com/golangci/golangci-lint/v2/cmd/golangci-lint`.
  - Version is pinned via the `require` block: `github.com/golangci/golangci-lint/v2 v2.8.0 // indirect`.
  - Makefile targets (`go-lint`, `go-fmt`, `go-lint-verify-config`) all use `go tool golangci-lint ...`.
  - CI (`lint-test.yml`) uses `golangci/golangci-lint-action@v8` with `version: v2.8.0`.
- **Implications**: Upgrade requires both `go get` to update `go.mod`/`go.sum` AND updating the `version:` field in CI workflow. Both must be kept in sync.

### linters.fast Behavior Change (v1 → v2)

- **Context**: Requirement 4 specifies a "fast mode" config using the v1 `linters.fast: true` property.
- **Sources Consulted**: https://golangci-lint.run/docs/product/migration-guide/#lintersfast
- **Findings**:
  - In v2, `linters.fast: true` is replaced by `linters.default: fast`.
  - The `--fast-only` CLI flag also filters to "fast" linters only.
  - The documentation explicitly warns these are "not strictly equivalent" to v1 behavior.
  - "Fast" linters as of v2 documentation: canonicalheader, copyloopvar, dupword, err113, errorlint, exptostd, fatcontext, ginkgolinter, gocritic, godot, goheader, govet, iface, importas, intrange, mirror, misspell, nakedret, nlreturn, noinlineerr, nolintlint, perfsprint, protogetter, revive, sloglint, staticcheck, tagalign, testifylint, usestdlibvars, usetesting, whitespace, wsl_v5.
- **Implications**: Pattern A config must use `linters.default: fast` (not `linters.fast: true`). Some linters currently in `.golangci.yml` are not "fast" and will be absent in Pattern A.

### Go Plugin System for Custom Rules

- **Context**: Requirement 2 proposes compiling `tools/lint/rules.go` as a native Go plugin.
- **Sources Consulted**: https://golangci-lint.run/docs/plugins/go-plugins/
- **Findings**:
  - The native Go plugin system requires a `plugin/main.go` exporting `func New(conf any) ([]*analysis.Analyzer, error)`.
  - Requirements: `CGO_ENABLED=1`, all overlapping library versions must exactly match golangci-lint's transitive dependencies.
  - Build command: `go build -buildmode=plugin -o plugin.so plugin/main.go`
  - Configuration: added under `linters.settings.custom.{name}.path`.
  - The current `tools/lint/rules.go` uses the `github.com/quasilyte/go-ruleguard/dsl` package (ruleguard DSL). It is loaded by gocritic's `ruleguard` enabled-check, **not** as a plugin.
  - The golangci-lint docs recommend the "Module Plugin System" over the native Go plugin system due to complexity and fragility of version matching.
- **Implications**: Converting to a native Go plugin requires:
  1. Creating a new `tools/lint/plugin/main.go` that wraps ruleguard logic as a `go/analysis.Analyzer`.
  2. Accepting significant CGO and version-pinning complexity.
  3. Alternative: keep the current gocritic/ruleguard approach (simpler, already working) and drop Requirement 2 if the performance gain is marginal.

### modernize Linter

- **Context**: Requirement 3 proposes introducing the `modernize` linter.
- **Sources Consulted**: https://golangci-lint.run/docs/linters/
- **Findings**:
  - `modernize` is described as: "A suite of analyzers that suggest simplifications to Go code, using modern language and library features."
  - It is **not** tagged "fast" in v2 documentation.
  - Suitable for the strict (nightly) profile; evaluate for lightweight if execution time is acceptable.
- **Implications**: Enable `modernize` in `.golangci-strict.yml`. Add exclusions for any patterns idiomatic to this project. Do not include in Pattern A or Pattern B initially.

### New gosec Rules in v2.10.0

- **Context**: Validate whether new gosec rules require exclusions in the existing config.
- **Sources Consulted**: golangci-lint changelog v2.10.0
- **Findings**:
  - Eight new gosec rules added: G117, G602, G701, G702, G703, G704, G705, G706.
  - The existing config already excludes G104, G115, G306, G307.
- **Implications**: After upgrade, run the strict profile to assess whether any new gosec rules fire on the existing codebase. Add targeted exclusions if needed.

### Nightly Workflow Pattern

- **Context**: Requirement 6 requires a scheduled GitHub Actions job.
- **Sources Consulted**: `.github/workflows/nightly-e2e.yml` (existing nightly pattern)
- **Findings**:
  - Existing nightly job uses `cron: "0 18 * * *"` (18:00 UTC daily).
  - Uses `ubuntu-latest` as runner (non-standard for this repo; other jobs use `ubuntu-slim`).
  - Has `workflow_dispatch` for manual triggers.
- **Implications**: New `nightly-lint.yml` should use `ubuntu-slim`, a different cron time (e.g., `0 2 * * *` for 02:00 UTC), `workflow_dispatch`, and run only on `main`.

## Architecture Pattern Evaluation

| Option | Description | Strengths | Risks / Limitations | Notes |
|--------|-------------|-----------|---------------------|-------|
| Three separate config files | `.golangci-fast.yml`, `.golangci-lightweight.yml`, `.golangci-strict.yml` | Clear separation; easy to switch by `-c` flag; independently testable | Config drift between files if not maintained; DRY violation | Preferred: matches user requirement for separate configs per usage |
| Base + override YAMLs | Single base config with profile-specific overrides via YAML anchors | DRY; single source of truth | golangci-lint v2 YAML does not support anchors between files natively | Not viable without custom merge tooling |
| Single config + CLI flags | One config, use `--fast-only` or `--enable`/`--disable` CLI flags | Minimal files | Hard to document; not reproducible from config file alone | Viable for fast mode but poor for CI documentation |

**Selected**: Three separate config files.

## Design Decisions

### Decision: Keep gocritic/ruleguard for Custom Rules

- **Context**: Requirement 2 proposes native Go plugin compilation for `tools/lint/rules.go`.
- **Alternatives Considered**:
  1. Native Go plugin (`-buildmode=plugin`) — requires CGO, exact version matching, new `plugin/main.go` wrapper
  2. Keep current gocritic/ruleguard source approach — zero new complexity, already working
  3. Module plugin system — experimental, requires separate module
- **Selected Approach**: Keep the current gocritic/ruleguard approach. Document the native plugin path as an optional future enhancement.
- **Rationale**: The native Go plugin system's CGO requirement and strict version matching introduce fragility that outweighs the marginal performance improvement for a single-rule ruleguard file. The `rules: "tools/lint/rules.go"` configuration in gocritic already compiles the rules on first run and caches them.
- **Trade-offs**: Simpler maintenance vs. potentially slower first-run on cache miss.
- **Follow-up**: Measure actual ruleguard overhead in CI timing; revisit if it contributes significantly to execution time.

### Decision: Pattern A uses `linters.default: fast`

- **Context**: Requirement 4 requires fast-mode using v2-equivalent of `linters.fast: true`.
- **Alternatives Considered**:
  1. `linters.default: fast` — v2 native config approach
  2. `--fast-only` CLI flag — not expressible in config file alone
  3. Manually list "fast" linters — fragile, not forward-compatible
- **Selected Approach**: `linters.default: fast` in `.golangci-fast.yml`.
- **Rationale**: Config-file-driven approach is reproducible and visible in CI workflow; aligns with v2 migration guide.
- **Trade-offs**: Linter set changes when golangci-lint adds/removes "fast" tags on linters.

### Decision: Strict Profile Extends Main Config

- **Context**: Requirement 6 needs a strict full-linter profile for nightly CI.
- **Selected Approach**: `.golangci-strict.yml` = `.golangci.yml` linter set + `modernize` + any new gosec rules not excluded.
- **Rationale**: The current `.golangci.yml` already represents the "standard" config. The strict profile adds depth rather than replacing it.

## Risks & Mitigations

- **go.sum regeneration on upgrade** — `go mod tidy` after `go get` may update many indirect dependencies; review diff carefully to avoid unexpected breakage.
- **gosec new rules (G117, G602, G701-G706)** may fire on existing code — mitigate by running strict profile locally before merging and adding targeted exclusions.
- **`linters.default: fast` may include linters not in the current config** — review the fast linter list against existing exclusion rules; ensure new fast linters don't produce noise.
- **Plugin system complexity** — native Go plugin approach is deferred; if revisited, test in an isolated branch first due to CGO requirements.
- **Config drift across three files** — mitigate by using shared `importas`, `depguard`, and `gosec` settings blocks and documenting in development guidelines which file to edit.

## References

- [golangci-lint Go Plugin System](https://golangci-lint.run/docs/plugins/go-plugins/) — plugin build requirements and configuration
- [golangci-lint v2 Migration Guide: linters.fast](https://golangci-lint.run/docs/product/migration-guide/#lintersfast) — `linters.default: fast` replacement
- [golangci-lint Linters List](https://golangci-lint.run/docs/linters/) — fast-tagged linters and modernize description
- [golangci-lint Changelog](https://golangci-lint.run/docs/product/changelog/) — v2.9.0 and v2.10.0 changes
- [golangci-lint v2.10.1 Release](https://github.com/golangci/golangci-lint/releases/tag/v2.10.1) — buildssa panic fix
- [golangci-lint modernize linter](https://golangci-lint.run/docs/linters/configuration/#modernize) — modernize configuration
