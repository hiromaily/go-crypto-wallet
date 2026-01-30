---
name: verify-buildkit-cache
description: Verify BuildKit cache effectiveness from CI logs for Go + Docker + GitHub Actions workflows using buildkit-cache-dance. Use when debugging cache issues, checking if Docker layer caching works, or when go build is slow despite cache configuration.
---

# Verifying BuildKit Cache Effectiveness

Systematically verify—**from CI logs only**—that Docker `RUN --mount=type=cache`
(Go `GOMODCACHE` and `GOCACHE`) is:

1. Persisted across jobs
2. Restored by GitHub Actions
3. Injected back into BuildKit
4. Actually speeding up `go build`

## Prerequisites

This checklist assumes:

- Dockerfile uses BuildKit cache mounts:
  - `RUN --mount=type=cache,target=/go/pkg/mod` (GOMODCACHE)
  - `RUN --mount=type=cache,target=/root/.cache/go-build` (GOCACHE)
- GitHub Actions workflow:
  - Uses `actions/cache@v5` to persist a host directory (`CACHE_DIR`, e.g. `cache-mount/`)
  - Uses `buildkit-cache-dance` (upstream or fork) to inject/extract cache mounts
- **Critical ordering**:
  - `actions/cache@v5` runs **before** `buildkit-cache-dance`
  - This guarantees `extract → save` via post-step execution order

## Final Pass/Fail Rule

Caching is **working correctly** if all three are true:

1. `actions/cache` restore shows **non-zero cache size**
2. Docker build logs show **large cache mounts BEFORE build**
3. `go build` runtime is dramatically reduced (e.g. 40s → a few seconds)

---

## CI Log Checklist

### A. `actions/cache@v5` Restore (Most Important)

#### ✅ PASS

Log contains:

```
Cache hit for: ...
Cache Size: ~XXX MB
```

(Any value well above 0 MB is acceptable)

#### ❌ FAIL

Restore shows:

```
Cache Size: ~0 MB (a few hundred bytes)
```

Or repeatedly logs:

```
Cache hit occurred on the primary key ... not saving cache.
```

→ Indicates an **empty cache permanently locked to the key**

#### Fix

- Bump `cache_version` to escape polluted keys
- Optionally include `github.ref_name` in the cache key
- Verify that buildkit-cache-dance actually writes into `CACHE_DIR`

---

### B. buildkit-cache-dance (Inject / Extract)

#### ✅ PASS

- Inject phase runs without errors
- Extract (post step) produces **real data inside `CACHE_DIR/<cache-id>/`**
- On the next run, `actions/cache` restores hundreds of MB

#### ❌ FAIL (Common)

- Extract appears to run, but `CACHE_DIR` remains empty
- Cache data is written somewhere else (e.g. repository root)
- Cache IDs or mount targets do not match the Dockerfile

#### Fix

Use absolute path for cache-dir:

```yaml
cache-dir: ${{ github.workspace }}/${{ env.CACHE_DIR }}
```

Ensure `cache-map` matches Dockerfile cache mounts:

- Same `target`
- Same `id` (if `id=` is used in Dockerfile)

For debugging, enable `is_debug=true` and inspect:

- Generated `Dancefile.inject`
- Generated `Dancefile.extract`
- Actual destination paths of extracted data

---

### C. Docker Build Stage (Runtime Proof)

Instrument the Dockerfile build stage with simple checks.

#### Recommended Observations

- `du -sh /go/pkg/mod /root/.cache/go-build` before and after `go build`
- `time go build ...`
- Presence of `go: downloading ...` logs

#### ✅ PASS (2nd run and later)

BEFORE build:

```
/go/pkg/mod            hundreds of MB
/root/.cache/go-build  hundreds of MB
```

- `go build` completes in **seconds**
- Very few or no `go: downloading ...` lines

#### ❌ FAIL

BEFORE build shows:

```
8.0K /go/pkg/mod
8.0K /root/.cache/go-build
```

- Every run downloads modules
- Build time remains ~40s+

---

## Expected Behavior by Run

### First Run (New cache_version)

- Cache restore: miss or very small
- Docker build BEFORE: small (OK)
- Build: slow (expected)
- Post step: extract + save creates the cache

### Second Run (Same cache_version)

- Cache restore: **hundreds of MB**
- Docker build BEFORE: **hundreds of MB**
- Build time: dramatically faster
- Minimal module downloads

---

## Quick Failure → Fix Mapping

| Symptom | Likely Cause | Fix |
|---------|--------------|-----|
| Cache Size ~0 MB | Empty cache locked to key | Bump `cache_version` |
| BEFORE = 8K | Inject failed | Fix cache-dance / cache-map |
| Primary hit, not saving | Empty cache frozen | New key space |
| Cache written to repo root | Bug in cache-dance | Force output into `CACHE_DIR` |
| Downloads every run | GOMODCACHE not restored | Verify `/go/pkg/mod` inject |

---

## Final Automated Pass Criteria

Declare **"cache is effective"** if:

- `actions/cache` restore reports **≥ 100 MB**
- Docker build BEFORE shows:
  - `/go/pkg/mod` ≥ 100 MB
  - `/root/.cache/go-build` ≥ 50 MB
- `go build` time is significantly reduced (e.g. < 10s)
- `go: downloading` logs are mostly absent

---

## Operational Notes

- Always bump `cache_version` after:
  - Dockerfile cache mount changes
  - cache-map changes
  - buildkit-cache-dance logic changes
- During debugging, use timestamped versions (e.g. `v4-0130-0943`)
- Once stable, freeze to a fixed version (e.g. `stable-v1`)
