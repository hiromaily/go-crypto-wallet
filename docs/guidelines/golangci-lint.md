<!--
⚠️ AUTO-GENERATED FILE — DO NOT EDIT
Source: template/pages/docs/guidelines/golangci-lint.tpl.md · Run `make docs` to regenerate.
-->

# golangci-lint

## Performance Profiling

### How to Measure

#### 1. Verbose mode (`-v`)

Shows active linter count and per-analyzer execution times.

```bash
go tool golangci-lint run -v
```

Key output lines:

- `[lintersdb] Active N linters: [...]` — lists all enabled linters
- `analyzers took ... with top 10 stages: ...` — per-analyzer breakdown (parallel; sum exceeds wall time)
- `linters took ... with stages: goanalysis_metalinter: ...` — actual wall time
- `Execution took ...` — total execution time

#### 2. Trace output (`--trace-path`)

Generates a Go execution trace viewable in browser.

```bash
go tool golangci-lint run --trace-path trace.out
go tool trace trace.out
```

#### 3. CPU profiling (`--cpu-profile-path`)

Generates a pprof CPU profile for detailed function-level analysis.

```bash
go tool golangci-lint run --cpu-profile-path cpu.prof
go tool pprof cpu.prof
```

> **Note:** golangci-lint does not provide per-linter timing directly. The `-v` output groups all linters under `goanalysis_metalinter`, but the `analyzers took` line shows individual analyzer stages. To isolate a specific linter, use `--enable-only <linter>`.

### Measurement Results (2026-02-27)

**Environment:** macOS Darwin 23.6.0, Go 1.25, golangci-lint v2 (via `go tool`)

**Conditions:** Cache cleared before each run (`go tool golangci-lint cache clean`)

#### Profile Comparison

| Profile | Config | Linters | Wall Time |
|---------|--------|---------|-----------|
| Base | `.golangci.yml` | 45 | ~31s |

**Difference:** Strict profile adds only `modernize` linter compared to base. The previously observed 65s for base was due to cold cache; with consistent cache clearing, both profiles take similar time.

#### Per-Analyzer Breakdown (Base Profile, Cache Cleared)

Total analyzer time: 11m12s (parallel execution, wall time ~31s)

| # | Analyzer | Time | Notes |
|---|----------|------|-------|
| 1 | buildir | 1m58s | SSA IR construction (prerequisite for SSA-based linters) |
| 2 | wastedassign | 56s | Depends on buildir |
| 3 | goimports | 56s | Import formatting check |
| 4 | gocritic | 50s | Multi-checker with many sub-analyzers |
| 5 | gosec | 43s | Security analysis (depends on buildir) |
| 6 | revive | 15s | Extensible linter framework |
| 7 | gofumpt | 12s | Strict gofmt |
| 8 | unparam | 11s | Unused function parameter detection |
| 9 | unconvert | 10s | Unnecessary type conversion detection |
| 10 | buildssa | 9s | SSA construction (used by multiple linters) |

#### Resource Usage

- Memory: avg 1,452 MB, max 2,398 MB (346 samples)
- Package loading: ~378ms

---
