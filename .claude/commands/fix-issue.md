# Fix Issue #{issue_number}

Work on a GitHub issue using the label-context-mapping skill.

## Process

1. **Fetch issue**: `gh issue view {issue_number}`
2. **Load label-context-mapping skill**: Determine what to load from labels
   - See [label-context-mapping](../skills/label-context-mapping/SKILL.md)
3. **Load determined Skills and Contexts**:
   - Always: `git-workflow`
   - Based on labels: Skills and contexts from mapping
4. **Follow loaded skill workflows**

## How It Works

```
fix-issue command
    │
    ├─ 1. Fetch issue labels
    │
    ├─ 2. Load label-context-mapping skill
    │      │
    │      ├─ Type label → Context document
    │      ├─ Lang/Scope label → Development skill
    │      ├─ Chain label → Chain context
    │      └─ Test label → Verification commands
    │
    └─ 3. Execute workflow from loaded skills + contexts
```

## Key References

| Reference | Purpose |
|-----------|---------|
| [label-context-mapping](../skills/label-context-mapping/SKILL.md) | Label → Skill/Context mapping |
| [task-classification.md](../../../docs/guidelines/task-classification.md) | SSOT for label definitions |

## Example

```
/fix-issue #123
```

If issue #123 has labels `bug`, `lang:go`, `chain:btc`:

1. Load `label-context-mapping` skill
2. Mapping determines:
   - Context: `docs/task-contexts/bug-fix.md` (from `bug`)
   - Skill: `go-development` (from `lang:go`)
   - Chain context: `docs/task-contexts/chains/btc.md` (from `chain:btc`)
3. Load `git-workflow` skill (always)
4. Load `go-development` skill
5. Follow workflow: branch → implement → verify → commit → PR
