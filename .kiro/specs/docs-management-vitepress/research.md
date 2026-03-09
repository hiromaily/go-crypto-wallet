# Research & Design Decisions

---
**Purpose**: Capture discovery findings, architectural investigations, and rationale that inform the technical design.

---

## Summary

- **Feature**: `docs-management-vitepress`
- **Discovery Scope**: New Feature (greenfield documentation site, but extending an existing docs/ directory)
- **Key Findings**:
  - VitePress supports bun natively; `bun run docs:dev/build/preview` work out of the box
  - Built-in local search (MiniSearch) requires zero external services — no Algolia account needed
  - The project's existing `docs/` directory is the natural `srcDir`; placing `.vitepress/` inside it avoids a separate `docs-site/` folder and eliminates duplicate content paths

## Research Log

### VitePress Setup and Bun Compatibility

- **Context**: Required to confirm bun works as the package manager for VitePress
- **Sources Consulted**: https://vitepress.dev/guide/getting-started
- **Findings**:
  - VitePress fully supports bun: `bun add -D vitepress`, `bun run docs:dev`
  - Scripts pattern: `docs:dev`, `docs:build`, `docs:preview` in `package.json`
  - Config file: `.vitepress/config.ts` (TypeScript supported natively)
  - Dev server: `http://localhost:5173`
  - Build output: `docs/.vitepress/dist`
- **Implications**: No npm fallback needed; CI already uses `oven-sh/setup-bun@v2` for other jobs — consistent setup

### VitePress Search Options

- **Context**: Requirement 4.1 mandates full-text search
- **Sources Consulted**: https://vitepress.dev/reference/default-theme-search
- **Findings**:
  - `provider: 'local'` enables MiniSearch-based in-browser full-text search
  - Zero external dependencies; index built at `bun run docs:build` time
  - Algolia is optional and not needed for this project
- **Implications**: Use `search: { provider: 'local' }` in `themeConfig`; no API keys or third-party accounts required

### GitHub Pages Deployment

- **Context**: Requirement 5.1 mandates auto-deploy on push to main
- **Sources Consulted**: https://vitepress.dev/guide/deploy#github-pages
- **Findings**:
  - Required permissions: `contents: read`, `pages: write`, `id-token: write`
  - Recommended actions: `actions/configure-pages@v5`, `actions/upload-pages-artifact@v3`, `actions/deploy-pages@v4`
  - `base` config option in `config.ts` must match the GitHub Pages sub-path (e.g., `/go-crypto-wallet/`)
  - `concurrency` group prevents overlapping deployments
- **Implications**: New workflow `docs.yml` handles deploy; `lint-test.yml` handles PR build check; runner is `ubuntu-slim` per project rules

### Existing Project Structure

- **Context**: Understand what already exists to avoid duplication
- **Sources Consulted**: Codebase analysis via Glob/Bash
- **Findings**:
  - `docs/` contains ~80+ markdown files organized into: `chains/`, `database/`, `design/`, `guidelines/`, `task-contexts/`, `tools/`, plus top-level files
  - `.claude/rules/` has ~15 markdown files, some duplicating guidance from `docs/guidelines/`
  - `.claude/skills/` has ~20 skill directories with `SKILL.md` files
  - `.claude/commands/` has markdown files for slash commands
  - No existing `package.json` at repo root; existing JS packages are under `apps/`
  - Project already uses bun in CI for `apps/xrpl-grpc-server` and `apps/eth-contracts`
- **Implications**: VitePress `package.json` placed at repo root (not in `apps/`); `srcDir` is `docs/`; sidebar must be manually curated to match the existing directory hierarchy

## Architecture Pattern Evaluation

| Option | Description | Strengths | Risks / Limitations | Notes |
|--------|-------------|-----------|---------------------|-------|
| Co-locate in `docs/.vitepress/` (root `package.json`) | VitePress config inside existing `docs/`, scripts at repo root | No new top-level directory; docs and config travel together | Root `package.json` is new — must not conflict with Go tooling | **Selected** |
| Separate `docs-site/` directory | Dedicated folder for VitePress setup | Clean isolation | Requires duplicating or symlinking all content from `docs/` | Rejected — violates SSOT |
| Inside `apps/` | Co-located with other JS apps | Consistent with existing apps structure | `apps/` is for runtime services; docs tooling is different | Rejected — wrong semantic layer |

## Design Decisions

### Decision: Place VitePress Config in `docs/.vitepress/`, Scripts at Repo Root

- **Context**: VitePress needs a `srcDir` and a config location; content lives in `docs/`
- **Alternatives Considered**:
  1. Separate `docs-site/` folder — clean but creates SSOT violation (content duplication or symlinks)
  2. Inside `apps/` — inconsistent; docs tooling ≠ runtime app
- **Selected Approach**: `docs/.vitepress/config.ts` as config; `package.json` at repo root with `docs:*` scripts; `srcDir` implicitly `docs/`
- **Rationale**: Keeps content and config together; no duplication; consistent with VitePress conventions
- **Trade-offs**: Adds a root `package.json` (new for this Go project); must ensure Go tools ignore it
- **Follow-up**: Verify `go.mod` tooling ignores `package.json`; add `node_modules/` and `docs/.vitepress/dist` to `.gitignore`

### Decision: Use Local Search (MiniSearch) Over Algolia

- **Context**: Requirement 4.1 requires search; two options available
- **Alternatives Considered**:
  1. Algolia DocSearch — cloud-hosted, requires account and crawler config
  2. Local MiniSearch — built-in, zero dependencies
- **Selected Approach**: `search: { provider: 'local' }`
- **Rationale**: No external accounts needed; works offline and in development; adequate for a single-repository docs site
- **Trade-offs**: Index built at build time (not crawled); fine for this project's scale

### Decision: Manual Sidebar Configuration Over Auto-generation

- **Context**: Requirement 1.5 requires hierarchical sidebar matching `docs/` structure
- **Alternatives Considered**:
  1. Auto-generate sidebar from filesystem at build time (custom plugin/script)
  2. Manually curated sidebar in `config.ts`
- **Selected Approach**: Manual sidebar defined in `config.ts`
- **Rationale**: The `docs/` directory has a stable, well-defined hierarchy; manual config gives full control over ordering and labeling; auto-generation adds complexity for minimal gain
- **Trade-offs**: New docs files require a sidebar config update (acceptable per Requirement 6.1)

### Decision: Separate `docs.yml` Workflow for Deploy, Update `lint-test.yml` for PR Check

- **Context**: Requirements 5.1 and 5.2 require both deploy-on-push and PR build check
- **Alternatives Considered**:
  1. Single workflow handling both PR check and deploy (with conditional job)
  2. Two separate workflows: `docs.yml` (deploy) and update to `lint-test.yml` (PR check)
- **Selected Approach**: Option 2 — separate workflows
- **Rationale**: `lint-test.yml` already handles path-filtered PR checks for all file types; adding a `docs` filter follows the existing pattern cleanly; deploy workflow is independent and has different permissions
- **Trade-offs**: Two files to maintain; the pattern is consistent with existing CI structure

## Risks & Mitigations

- **Root `package.json` conflicts with Go tooling** — Mitigation: Go tools ignore `package.json`; add `bun.lock` and `node_modules/` to `.gitignore`; verify `make` targets unaffected
- **`base` URL mismatch on GitHub Pages** — Mitigation: Set `base: '/go-crypto-wallet/'` in config; test with `bun run docs:preview` locally before deploy
- **Sidebar drift as docs grow** — Mitigation: Document the sidebar update process; Requirement 6.1 accepts minimal manual changes
- **SSOT enforcement is labour-intensive** — Mitigation: Scope to clearly duplicated content; do not force-reference every line; use PR reviews to enforce the principle incrementally

## References

- [VitePress Getting Started](https://vitepress.dev/guide/getting-started) — installation, directory structure, bun support
- [VitePress Site Config](https://vitepress.dev/reference/site-config) — `base`, `srcDir`, `outDir`, `title`, `description`
- [VitePress Search](https://vitepress.dev/reference/default-theme-search) — local MiniSearch vs Algolia
- [VitePress GitHub Pages Deploy](https://vitepress.dev/guide/deploy#github-pages) — GitHub Actions workflow, permissions, artifact upload
