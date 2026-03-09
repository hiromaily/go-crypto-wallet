# Implementation Plan

## Task Format

Tasks follow the migration phases from design.md:
Phase 1–2 → Tasks 1–3 (toolchain setup, configuration, home page)
Phase 3 → Task 4 (CI/CD)
Phase 5 → Task 5 (SSOT enforcement)

---

- [x] 1. Set up VitePress toolchain and project files

- [x] 1.1 Initialize root package.json and install VitePress with bun
  - Create `package.json` at repo root declaring `vitepress@^1.6.4` as a dev dependency (never use `"latest"`)
  - Add scripts: `docs:dev` (`vitepress dev docs`), `docs:build` (`vitepress build docs`), `docs:preview` (`vitepress preview docs`)
  - Run `bun install` to generate `bun.lock`; commit both `package.json` and `bun.lock`
  - Add `node_modules/`, `docs/.vitepress/dist/`, and `docs/.vitepress/cache/` to `.gitignore`
  - Verify `bun run docs:dev` starts the development server at `http://localhost:5173` without errors
  - _Requirements: 1.1, 1.2, 1.3, 5.4_

- [x] 1.2 Create VitePress TypeScript config with site metadata, search, and theme settings
  - Create `docs/.vitepress/config.ts` using `defineConfig` from vitepress (TypeScript; never use `.js`)
  - Set `title`, `description`, `base: '/go-crypto-wallet/'`, and `cleanUrls: true`
  - Enable local full-text search: `themeConfig.search = { provider: 'local' }`
  - Keep the default VitePress theme with no custom CSS overrides
  - Add `socialLinks` pointing to the GitHub repository
  - Confirm the default theme renders breadcrumb navigation and per-page table of contents automatically (no extra config needed)
  - Run `bun run docs:build` and assert zero errors before proceeding to Task 2
  - _Requirements: 1.1, 4.1, 4.2, 4.4, 6.2, 6.3, 6.4_

---

- [x] 2. Configure navigation bar and sidebar

- [x] 2.1 Define the top-level navigation bar
  - Add a `nav` array to `themeConfig` with the following top-level items:
    - Getting Started (links to installation / commands page)
    - Architecture (links to guidelines/architecture)
    - Chains (dropdown with BTC, BCH, ETH, XRP, Cosmos entries)
    - Database (links to database overview)
    - Guidelines (links to guidelines index)
    - AI Agent & Dev Workflow (links to agent-skills or task-contexts overview)
  - Each nav item must link to a page that exists in `docs/`; verify with `bun run docs:build`
  - _Requirements: 1.4, 2.6, 3.1, 3.2, 3.3, 3.4, 3.5, 3.6_

- [x] 2.2 Define the complete sidebar hierarchy covering all docs/ content
  - Define a path-keyed `sidebar` object in `themeConfig` with groups for each major section:
    - **Overview**: overview, transaction-flow (top-level pages)
    - **Getting Started**: Installation, commands, devcontainer, proto
    - **Architecture**: guidelines/architecture, guidelines/core, guidelines/multi-chain, directory_structure
    - **Chains — BTC**: README, architecture, overview (address-types, technical-reference), taproot, descriptor (all sub-pages), psbt (all sub-pages), musig2 (all sub-pages), keygen, operations, testing; archive sub-pages as a collapsed group
    - **Chains — BCH**: README, interface-separation
    - **Chains — ETH**: README, architecture, transaction-patterns, solidity-development, json-schema, multisig
    - **Chains — XRP**: README, architecture-2026, architecture-xrpl-grpc-server-version, transaction-flow, library-selection, network-devnet, testing-strategy, setup, xrpl-go
    - **Chains — Cosmos**: README
    - **Database**: architecture, db-management, atlas-migration-flow, sqlc-code-generation-flow, schema-changes, quick-reference
    - **Guidelines**: README, architecture, core, multi-chain, coding-conventions, testing, workflow, security, release, code-generation, task-classification, requirements, claude-mem
    - **AI Agent & Dev Workflow**: agent-skills, design/ai-agents-instruction, task-contexts (all sub-pages)
    - **Design Notes**: all pages under design/ (btc-network-mode-switching, revise-db-atlas-sqlc-flow, superpowers-integration, pg2sqlite-alter-table-support, claude-best-practice)
    - **Tools & CI**: tools/golangci-lint, github-actions/investigation-cache
    - **Issues / Archive**: issues/REFACTORING_PLAN, issues/REFACTORING_CHECKLIST (collapsed by default)
  - Japanese-language files (`overview-ja.md`, `keygen/improvements-2025-ja.md`) may be included in their respective sections or omitted; note this decision in the config
  - Every `link` value must resolve to an existing `.md` file; `bun run docs:build` catches broken links automatically
  - _Requirements: 1.5, 2.1, 2.6, 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 3.7, 6.1_

---

- [x] 3. Create the documentation home page
  - Create `docs/index.md` with VitePress home layout frontmatter (`layout: home`)
  - Add a hero section: project name (go-crypto-wallet), a one-sentence tagline describing the multi-chain wallet, and action buttons ("Get Started" → Installation, "View on GitHub" → repo URL)
  - Add a features section with cards linking to the six major documentation areas: Architecture, Chains, Getting Started, Development Guidelines, Database, AI Agent & Dev Workflow
  - Run `bun run docs:build` and `bun run docs:preview` to confirm the home page renders correctly at the site root
  - _Requirements: 4.3_

---

- [x] 4. Set up GitHub Actions CI/CD

- [x] 4.1 Create the docs deploy workflow
  - Create `.github/workflows/docs.yml` with:
    - Trigger: push to `main` with path filter `docs/**`, `package.json`, `bun.lock`; also `workflow_dispatch`
    - Permissions: `contents: read`, `pages: write`, `id-token: write`
    - Concurrency group `pages` with `cancel-in-progress: false` (allow production deploy to finish)
    - Single job (`deploy`) on `ubuntu-slim`, `timeout-minutes: 10`
    - Steps: `actions/checkout@v6` → `oven-sh/setup-bun@v2` → `bun install --frozen-lockfile` → `bun run docs:build` → `actions/configure-pages@v5` → `actions/upload-pages-artifact@v3` (path: `docs/.vitepress/dist`) → `actions/deploy-pages@v4`
  - Note in implementation: GitHub Pages must be enabled in repo Settings → Pages → Source: GitHub Actions (one-time manual step)
  - _Requirements: 5.1, 5.3_

- [x] 4.2 (P) Extend lint-test.yml with a docs build check job
  - Add a `docs` entry to the `dorny/paths-filter` step: paths `docs/**`, `package.json`, `bun.lock`
  - Add a `docs-build` job that runs when `needs.changes.outputs.docs == 'true'`
  - Job config: `ubuntu-slim`, `timeout-minutes: 5`, `permissions: contents: read`
  - Steps: `actions/checkout@v6` → `oven-sh/setup-bun@v2 (bun-version: latest)` → `bun install --frozen-lockfile` → `bun run docs:build`
  - This job can run in parallel with other lint jobs (independent file set, no shared resources)
  - _Requirements: 5.2, 5.3_

---

- [x] 5. Enforce SSOT across .claude/ files and directory READMEs

- [x] 5.1 (P) Audit and update .claude/rules/ files
  - For each file under `.claude/rules/`, identify sections where body content is a verbatim or near-verbatim copy (>50% overlap) of a section in `docs/`
  - Replace duplicate body content with a Markdown link to the canonical `docs/` page and a one-line summary of the linked content
  - Preserve meta-instructions, agent-specific guidance, and content that merely references or links to `docs/` — these are not duplicates
  - After updates, confirm every `docs/` page referenced actually exists (search the VitePress site or check the file tree)
  - Can be executed in parallel with 5.2, 5.3, 5.4 (disjoint file sets)
  - _Requirements: 2.2_

- [x] 5.2 (P) Audit and update .claude/skills/ files
  - For each `SKILL.md` under `.claude/skills/*/`, identify sections duplicating `docs/` content using the same >50% overlap criterion
  - Replace duplicate sections with links to the canonical `docs/` page; preserve skill-specific invocation instructions and examples
  - After updates, confirm every `docs/` page referenced actually exists (search the VitePress site or check the file tree)
  - Can be executed in parallel with 5.1, 5.3, 5.4
  - _Requirements: 2.3_

- [x] 5.3 (P) Audit and update .claude/commands/ files
  - For each file under `.claude/commands/`, identify sections duplicating `docs/` content
  - Replace duplicate sections with links to the canonical `docs/` page; preserve command-specific prompts and invocation logic
  - After updates, confirm every `docs/` page referenced actually exists (search the VitePress site or check the file tree)
  - Can be executed in parallel with 5.1, 5.2, 5.4
  - _Requirements: 2.4_

- [x] 5.4 (P) Audit and update directory-level README.md files
  - Focus on READMEs in `internal/`, `apps/`, `scripts/`, `cmd/`, and other directories that contain prose duplicating `docs/` content
  - Replace duplicated explanatory sections with links to the canonical `docs/` page; preserve per-directory quick-start steps that are directory-specific and not in `docs/`
  - After updates, run `bun run docs:build` to confirm all newly linked `docs/` pages exist (no broken references)
  - Can be executed in parallel with 5.1, 5.2, 5.3
  - _Requirements: 2.5_

---

- [ ] 6. Fix dead links in docs/ and remove ignoreDeadLinks workaround

  **Context**: `ignoreDeadLinks: true` was set in Task 1 as a temporary workaround because 156 relative
  links inside `docs/` point outside the `docs/` directory (e.g. `../../tools/atlas/README.md`,
  `../../../AGENTS.md`, `../../.claude/skills/…`). These cause VitePress build failures. This task
  cleans them up so the flag can be removed.

- [ ] 6.1 Audit and categorize all out-of-docs relative links
  - Run: `grep -rn "](\.\./" docs/ --include="*.md"` (excluding srcExclude paths) to collect all 156 links
  - Categorize each into one of three fix strategies:
    - **A — Remap to docs/ page**: link target has an equivalent page already in `docs/` → update to VitePress-style absolute path (e.g. `/database/architecture`)
    - **B — Convert to GitHub URL**: target is a source file, script, or config with no docs/ equivalent → replace with a `https://github.com/hiromaily/go-crypto-wallet/blob/main/…` URL
    - **C — Convert to plain text**: target is an agent-only file (`.claude/`, `AGENTS.md`, `CLAUDE.md`) that should not be a hyperlink in user-facing docs → remove the markdown link syntax, keep the text as a code-span reference
  - Produce a mapping table (file → line → category → fix) before making any edits
  - _Requirements: 1.3_

- [ ] 6.2 Apply fixes — category A (remap to docs/ pages)
  - For each category-A link, replace the relative `../` path with the VitePress absolute path (leading `/`)
  - Key path mappings expected:
    - `../../tools/atlas/README.md` → `/database/architecture` (or `/database/db-management`)
    - `../development/database.md` → `/database/architecture`
    - `../guidelines/coding-standards.md` → `/guidelines/coding-conventions`
    - `../crypto/btc/operations/e2e-transaction-patterns.md` → `/chains/btc/operations/e2e-transaction-patterns`
    - `../Installation.md` → `/Installation`
  - After each file edit, confirm the target path exists in `docs/`
  - _Requirements: 1.3_

- [ ] 6.3 Apply fixes — category B (convert to GitHub URLs)
  - For each category-B link (source files, scripts, configs), replace with:
    `https://github.com/hiromaily/go-crypto-wallet/blob/main/<repo-relative-path>`
  - Key examples: `../../tools/atlas/README.md` references from non-database docs, `../internal/interface-adapters/cli/README.md`
  - _Requirements: 1.3_

- [ ] 6.4 Apply fixes — category C (convert to plain text)
  - For each category-C link (`.claude/`, `AGENTS.md`, `CLAUDE.md`, `ARCHITECTURE.md`), remove the hyperlink and render as a code-span: `` `.claude/skills/label-context-mapping/` `` instead of a broken link
  - These files are agent-configuration files not served by the docs site
  - _Requirements: 1.3, 2.2_

- [ ] 6.5 Remove ignoreDeadLinks and verify clean build
  - Remove `ignoreDeadLinks: true` (and its comment) from `docs/.vitepress/config.ts`
  - Run `bun run docs:build` — must complete with **zero dead-link errors**
  - If any remaining dead links appear, fix them before proceeding
  - _Requirements: 1.3_
