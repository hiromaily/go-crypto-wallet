---
name: docs-ssot
description: Set up docs-ssot SSOT documentation structure — migrate existing docs, build, and validate
allowed-tools:
  - Read
  - Write
  - Edit
  - Bash(docs-ssot *)
  - Bash(make docs*)
  - Bash(git diff *)
---
# docs-ssot: Documentation SSOT Setup

This skill migrates existing Markdown documentation into a modular Single Source of Truth (SSOT) structure managed by `docs-ssot`, then builds and validates the output.

## When to use this skill

- Setting up `docs-ssot` in a new or existing repository (follow the main **Workflow**)
- Adding a new generated output file (e.g., `SETUP.md`) to an existing SSOT structure (follow **Adding a new output file**)
- Regenerating documentation after editing source templates (run `docs-ssot build`)

---

## Quick Command Reference

| Command | Purpose |
|---------|---------|
| `docs-ssot migrate <files>` | Decompose existing docs into SSOT section structure |
| `docs-ssot migrate --dry-run <files>` | Preview migration without writing files |
| `docs-ssot validate` | Dry-run: check all includes resolve (no files written) |
| `docs-ssot build` | Generate all output files from templates |
| `docs-ssot check` | Detect near-duplicate sections (SSOT violations) |
| `docs-ssot index` | Show include relationships and orphan detection |
| `docs-ssot include <template>` | Expand and print a template to stdout (debugging) |

---

## Workflow

### Step 1 — Check prerequisites

Verify required tools are installed:

```sh
docs-ssot version   # docs-ssot itself
jq --version        # required by the generated-file protection hook
yq --version        # required by the generated-file protection hook and lefthook check
```

Install `docs-ssot` if missing:

```sh
# Homebrew (macOS/Linux)
brew tap hiromaily/tap && brew install docs-ssot

# or Go install
go install github.com/hiromaily/docs-ssot/cmd/docs-ssot@latest
```

---

### Step 2 — Check if already migrated

Before running migration, check whether this repository already has a docs-ssot structure:

```sh
ls docsgen.yaml 2>/dev/null
ls template/sections/ 2>/dev/null
```

If both exist, the repository is already (partially or fully) migrated. **Do not run `docs-ssot migrate` on generated output files.** Instead, skip to Step 5 to review the structure, or Step 8 to validate the current state.

If `docsgen.yaml` does not exist, or `template/sections/` is empty or missing, proceed with Steps 3–6.

---

### Step 3 — Identify existing documentation files

List Markdown files in the repository root:

```sh
ls *.md
```

Common candidates: `README.md`, `CLAUDE.md`, `AGENTS.md`, `SETUP.md`, `ARCHITECTURE.md`

**Do not migrate files that are already outputs in `docsgen.yaml`** — migrating a build artifact instead of a source file will produce incorrect results.

---

### Step 4 — Preview and run the migration

First, preview without writing files:

```sh
docs-ssot migrate --dry-run README.md CLAUDE.md
```

Review the output to understand how sections will be categorised and which are detected as duplicates. Then run the actual migration:

```sh
docs-ssot migrate README.md CLAUDE.md AGENTS.md
```

This will:
1. Split each file by H2 headings into section files under `template/sections/<category>/`
2. Create template files under `template/pages/` with `@include` directives
3. Create or update `docsgen.yaml` with build targets
4. Verify round-trip: build and compare output against originals

**`docsgen.yaml` reference** (the file created/updated by migration):

```yaml
index:
  output: template/INDEX.md   # optional: generates include-relationship index

targets:
  - input: template/pages/README.tpl.md
    output: README.md
  - input: template/pages/CLAUDE.tpl.md
    output: CLAUDE.md
```

Add or remove targets to control which files are generated.

---

### Step 5 — Review generated structure

```sh
docs-ssot index
```

Inspect:
- `template/sections/` — modular section files (**edit these, not the generated outputs**)
- `template/pages/*.tpl.md` — template files defining document structure
- `docsgen.yaml` — build targets mapping templates to output files

---

### Step 6 — Audit for duplicate content

**Run this before creating any new section file.** The `docs-ssot migrate` command deduplicates across the files it processes, but any additional content you create manually must be checked:

```sh
docs-ssot check
```

Also inspect existing sections:

```sh
ls template/sections/
```

**Rule:** If content for a new section already exists in a section used by another template (e.g., `installation-guide.md` already covers prerequisites and setup), do **not** create a duplicate. Instead, have the new template include the existing section:

```markdown
<!-- template/pages/SETUP.tpl.md -->
<!-- @include: ../sections/development/installation-guide.md -->
<!-- @include: ../sections/development/setup-release.md -->
```

Only create a new section file when the content is genuinely new and not covered anywhere else.

---

### Step 7 — Heading level convention

**All section files under `template/sections/` must start at heading level 2 (`##`).**

```markdown
## My Section Title       ← ✅ correct

# My Section Title        ← ❌ wrong
```

**Why:** Section files are embedded into larger documents where `#` is reserved for the document title. Starting at `##` means most includes need no `level` parameter.

Exception — when a section file is used as a **standalone output** (e.g., `.claude/rules/*.md`), use `level=-1` in the template to shift `##` → `#`:

```markdown
<!-- @include: ../sections/ai/rules/docs.md level=-1 -->
```

---

### Step 8 — Validate and build

First validate (dry-run — checks includes without writing files):

```sh
docs-ssot validate
```

If validation passes, build:

```sh
docs-ssot build
# or, if a Makefile target exists:
make docs
```

Verify output is correct:

```sh
git diff README.md CLAUDE.md
```

**Never edit** generated outputs (`README.md`, `CLAUDE.md`, or any `docsgen.yaml` output) directly — they are overwritten on every build.

---

## Include directive syntax

Templates use include directives to compose sections:

```markdown
<!-- @include: ../sections/project/overview.md -->
<!-- @include: ../sections/development/ -->
<!-- @include: ../sections/**/*.md level=+1 -->
```

| Parameter | Effect |
|-----------|--------|
| `level=+1` | `##` → `###` (deepen by one) |
| `level=-1` | `###` → `##` (shallow by one) |
| `level=0` | no change (same as omitting) |

Paths are resolved relative to the file containing the directive. Includes are expanded recursively.

---

## Adding a new output file

To add a new generated file (e.g., `SETUP.md`) to an existing SSOT structure:

1. **Audit existing sections** (see Step 6) — do not duplicate content that already exists.
2. **Create the template page** at `template/pages/SETUP.tpl.md`, including existing sections where applicable.
3. **Register in `docsgen.yaml`**:
   ```yaml
   - input: template/pages/SETUP.tpl.md
     output: SETUP.md
   ```
4. **Validate and build** (see Step 8):
   ```sh
   docs-ssot validate && docs-ssot build
   git diff SETUP.md
   ```

---

## Optional enhancements

### VitePress integration

If the repository has a VitePress site, detect it by checking for `docs/.vitepress/`:

```sh
ls docs/.vitepress/ 2>/dev/null
```

If present, convert each `docs/` page from standalone content into a thin `@include` wrapper so `template/sections/` becomes the single source of truth for both the generated root-level files and the VitePress site.

**Before** (`docs/guide/installation.md` with standalone content):
```markdown
# Installation
## Prerequisites
...full content here...
```

**After** (`docs/guide/installation.md` as a thin wrapper):
```markdown
<!-- @include: ../../template/sections/development/installation-guide.md -->
```

The canonical content lives in `template/sections/` and is referenced by any template that needs it.

---

### Generated-file protection hook (Claude Code)

Prevent AI agents from directly editing generated files. The hook reads `docsgen.yaml` at runtime so the block list stays in sync automatically.

Create `.claude/hooks/prevent-generated-edit.sh`:

```sh
#!/bin/sh
# Blocks Edit and Write tool use on auto-generated files.
# Generated files are defined as outputs in docsgen.yaml.
# Edit the source files in template/ instead, then run `make docs`.

# Requires: jq, yq
command -v jq >/dev/null 2>&1 || { echo "jq required but not found — hook inactive" >&2; exit 0; }
command -v yq >/dev/null 2>&1 || { echo "yq required but not found — hook inactive" >&2; exit 0; }

FILE=$(echo "$TOOL_INPUT" | jq -r '.file_path // empty')
if [ -z "$FILE" ]; then
  exit 0
fi

REPO_ROOT=$(git rev-parse --show-toplevel 2>/dev/null)
REL_PATH="${FILE#"${REPO_ROOT}/"}"

CONFIG="docsgen.yaml"
if [ ! -f "$CONFIG" ]; then
  exit 0
fi

GENERATED=$(yq -r '.targets[].output' "$CONFIG" 2>/dev/null)
if [ -z "$GENERATED" ]; then
  exit 0
fi

if echo "$GENERATED" | grep -Fqx "$REL_PATH"; then
  echo "BLOCKED: '$REL_PATH' is auto-generated by docs-ssot. Edit the source in template/ instead, then run 'make docs'." >&2
  exit 2
fi

exit 0
```

Make it executable:

```sh
chmod +x .claude/hooks/prevent-generated-edit.sh
```

Register in `.claude/settings.json`. **If the file already exists, merge into the existing `hooks.PreToolUse` array** rather than overwriting:

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Edit",
        "hooks": [
          { "type": "command", "command": ".claude/hooks/prevent-generated-edit.sh" }
        ]
      },
      {
        "matcher": "Write",
        "hooks": [
          { "type": "command", "command": ".claude/hooks/prevent-generated-edit.sh" }
        ]
      }
    ]
  }
}
```

---

### Pre-push docs-check (lefthook)

Add to `lefthook.yml` to fail a push when generated files are stale. Derive the file list from `docsgen.yaml` so it stays in sync automatically:

```yaml
pre-push:
  commands:
    docs-check:
      glob: "{template/**/*,docsgen.yaml}"
      run: |
        docs-ssot build
        docs-ssot index
        if ! yq -r '(.targets[].output, .index.output) | select(. != null)' docsgen.yaml | xargs -I {} git diff --quiet {}; then
          echo "ERROR: Generated files are out of date. Run 'make docs' and commit the changes." >&2
          exit 1
        fi
```

---

### Document the tooling in the README

Add a note so contributors know the documentation is managed by docs-ssot.

Find the section file for the README intro:

```sh
head -5 template/pages/README.tpl.md  # the first @include is the intro section
```

Open that section file and append:

```markdown
> **Documentation** is managed as a Single Source of Truth using [docs-ssot](https://github.com/hiromaily/docs-ssot).
> Files such as `README.md`, `CLAUDE.md`, and `ARCHITECTURE.md` are auto-generated —
> edit the source files under `template/` and run `make docs` to regenerate.
```

Then rebuild:

```sh
docs-ssot build
```
