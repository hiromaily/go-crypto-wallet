# Requirements Document

## Introduction

This specification defines requirements for setting up and managing a VitePress-based documentation site for the go-crypto-wallet repository. Documentation is currently scattered across `docs/`, per-directory `README.md` files, `.claude/rules/`, and `.claude/skills/`. The goal is to consolidate this into a unified, navigable documentation site while enforcing SSOT (Single Source of Truth): `.claude/rules`, `.claude/skills`, `.claude/commands` and directory READMEs shall reference canonical content under `docs/` rather than duplicating it.

## Requirements

### Requirement 1: VitePress Site Setup

**Objective:** As a developer, I want a VitePress site configured in this repository, so that all project documentation is accessible as a navigable static website.

#### Acceptance Criteria

1. The VitePress site shall be initialized under a dedicated directory (e.g., `docs-site/` or within `docs/.vitepress/`) with a valid `package.json` and VitePress configuration file, using bun as the package manager.
2. When `bun run docs:dev` is executed, the VitePress development server shall start and serve the documentation site locally without errors.
3. When `bun run docs:build` is executed, the VitePress static site build shall complete without errors and output files to a designated build directory.
4. The VitePress configuration shall define a navigation bar covering all major documentation areas (architecture, guidelines, chain-specific docs, usage).
5. The VitePress configuration shall define a sidebar that organizes all content hierarchically, matching the directory structure under `docs/`.

---

### Requirement 2: Documentation Consolidation and SSOT Enforcement

**Objective:** As a developer, I want all authoritative documentation to reside under `docs/`, so that `.claude/rules`, `.claude/skills`, `.claude/commands`, and directory READMEs can reference it without duplication.

#### Acceptance Criteria

1. The VitePress documentation site shall include all existing content from `docs/` as browsable pages, organized by topic area (guidelines, chains, database, design, task-contexts).
2. Where `.claude/rules/*.md` files contain documentation that duplicates content in `docs/`, the rules files shall be updated to reference the canonical `docs/` source instead of duplicating it.
3. Where `.claude/skills/*/SKILL.md` files contain documentation that duplicates content in `docs/`, the skill files shall be updated to reference the canonical `docs/` source instead of duplicating it.
4. Where `.claude/commands/*.md` files contain documentation that duplicates content in `docs/`, the command files shall be updated to reference the canonical `docs/` source instead of duplicating it.
5. Where directory-level `README.md` files contain documentation that duplicates content in `docs/`, those READMEs shall be updated to reference the canonical `docs/` source.
6. The VitePress site shall include a dedicated section documenting AI-agent usage (Claude Code rules, skills, commands, and workflow conventions) by referencing the content in `.claude/rules/`, `.claude/skills/`, and `.claude/commands/`.

---

### Requirement 3: Content Coverage

**Objective:** As a project contributor or user, I want the documentation site to cover all major topic areas, so that I can find architecture, implementation design, and usage instructions in one place.

#### Acceptance Criteria

1. The VitePress site shall include an **Architecture** section covering Clean Architecture layers, multi-chain design, and dependency rules (sourced from `docs/guidelines/architecture.md`, `docs/guidelines/core.md`).
2. The VitePress site shall include a **Chain-Specific** section with subsections for BTC, BCH, ETH, and XRP, each covering operations, key generation, signing, and transaction patterns.
3. The VitePress site shall include a **Getting Started / Usage** section covering installation (`docs/Installation.md`), CLI commands (`docs/commands.md`), and wallet operation flows.
4. The VitePress site shall include a **Development Guidelines** section covering coding conventions, testing, workflow, security, and task classification.
5. The VitePress site shall include a **Database** section covering schema management, Atlas migrations, and SQLC code generation.
6. The VitePress site shall include an **AI Agent & Dev Workflow** section explaining the spec-driven development process, Claude Code rules, and available skills.
7. If a documentation page is referenced from a rule or skill file, the VitePress site shall render that page without broken links.

---

### Requirement 4: Navigation and Discoverability

**Objective:** As a reader, I want intuitive navigation and search capabilities, so that I can quickly find the information I need.

#### Acceptance Criteria

1. The VitePress site shall enable the built-in full-text search feature so that readers can search across all documentation pages.
2. The VitePress site shall display breadcrumb navigation so that readers always know their location in the documentation hierarchy.
3. When a reader visits the documentation root, the VitePress site shall display an index or overview page that describes the project and links to major sections.
4. The VitePress site shall include a table of contents (page outline) in the sidebar or article column for pages with multiple headings.

---

### Requirement 5: Automation and CI Integration

**Objective:** As a maintainer, I want the documentation build to be automated, so that the site is always up to date without manual intervention.

#### Acceptance Criteria

1. When code is pushed to the main branch, the CI pipeline shall automatically build the VitePress site and publish it to the configured hosting target (e.g., GitHub Pages).
2. When a pull request is opened, the CI pipeline shall run `bun run docs:build` and report success or failure as a required status check.
3. If the VitePress build fails due to broken links or invalid configuration, the CI pipeline shall fail and provide an actionable error message.
4. The VitePress build output shall be excluded from Git tracking via `.gitignore`.

---

### Requirement 6: Maintainability

**Objective:** As a developer, I want documentation management to be straightforward, so that documentation stays current as the codebase evolves.

#### Acceptance Criteria

1. The VitePress sidebar configuration shall be maintainable: when a new `docs/` file is added, it shall require only minimal configuration changes to appear in the sidebar.
2. The VitePress site shall support Markdown frontmatter so that individual pages can override title, description, and layout.
3. The VitePress configuration shall be written in TypeScript or JavaScript and co-located with the documentation source, making it easy to locate and modify.
4. The VitePress site shall use a consistent theme (default VitePress theme or a minimal customization) so that the visual style is coherent without requiring bespoke CSS maintenance.
